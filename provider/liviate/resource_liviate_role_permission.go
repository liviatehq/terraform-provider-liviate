//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package liviate

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// A CloudStack role's rules are evaluated in LIST ORDER, first match wins -- a role with a
// trailing "*" deny catch-all (the normal pattern for custom roles) makes any rule added AFTER
// it dead on arrival, but createRolePermission always APPENDS. This resource fixes that: after
// creating a rule, if the role has any wildcard ("*") rule, it reorders the whole list so every
// specific (non-wildcard) rule stays ahead of every wildcard rule -- exactly the ordering a
// specific allow/deny needs to actually take effect against a catch-all. This is a deliberate
// simplification (not a general-purpose "order" attribute): it solves the one case that actually
// comes up in practice -- add a specific override on top of an existing catch-all -- without
// asking every user of this resource to reason about full list ordering themselves.
func resourceCloudStackRolePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackRolePermissionCreate,
		Read:   resourceCloudStackRolePermissionRead,
		Update: resourceCloudStackRolePermissionUpdate,
		Delete: resourceCloudStackRolePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the role this permission belongs to.",
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The API name (e.g. listVolumes) or a wildcard rule (e.g. list*, *).",
			},
			"permission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "allow or deny.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackRolePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	roleID := d.Get("role_id").(string)
	rule := d.Get("rule").(string)
	permission := d.Get("permission").(string)

	p := cs.Role.NewCreateRolePermissionParams(permission, roleID, rule)
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	}

	log.Printf("[DEBUG] Creating RolePermission %s on role %s", rule, roleID)
	r, err := cs.Role.CreateRolePermission(p)
	if err != nil {
		return fmt.Errorf("Error creating RolePermission: %s", err)
	}
	d.SetId(r.Id)

	if err := reorderSpecificRulesBeforeWildcards(cs, roleID); err != nil {
		return err
	}

	return resourceCloudStackRolePermissionRead(d, meta)
}

// Moves every non-"*" rule ahead of every "*" rule on the role, preserving relative order within
// each group. Idempotent -- safe to call after every create.
//
// Retries on CloudStack's "rule permissions list has changed while you were making updates"
// (CSExceptionErrorCode 4250) -- an optimistic-concurrency check that fires when two of these
// reorders race on the same role. Terraform runs a for_each's Create calls concurrently by
// default (parallelism=10), so N liviate_role_permission resources on the same role_id WILL
// collide here in normal use, not just as an edge case -- found live applying 9 of these at once.
func reorderSpecificRulesBeforeWildcards(cs *cloudstack.CloudStackClient, roleID string) error {
	const maxAttempts = 8
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lp := cs.Role.NewListRolePermissionsParams()
		lp.SetRoleid(roleID)
		lr, err := cs.Role.ListRolePermissions(lp)
		if err != nil {
			return fmt.Errorf("Error listing RolePermissions for reorder: %s", err)
		}

		var specific, wildcard []string
		for _, rp := range lr.RolePermissions {
			if rp.Rule == "*" {
				wildcard = append(wildcard, rp.Id)
			} else {
				specific = append(specific, rp.Id)
			}
		}
		if len(wildcard) == 0 {
			return nil // nothing to reorder around
		}

		up := cs.Role.NewUpdateRolePermissionParams(roleID)
		up.SetRuleorder(append(specific, wildcard...))
		if _, err := cs.Role.UpdateRolePermission(up); err != nil {
			if !strings.Contains(err.Error(), "has changed while you were making updates") {
				return fmt.Errorf("Error reordering RolePermissions on role %s: %s", roleID, err)
			}
			lastErr = err
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond) // linear backoff, jitter-free is fine at this scale
			continue
		}
		return nil
	}
	return fmt.Errorf("Error reordering RolePermissions on role %s after %d attempts (concurrent reorders kept colliding): %s", roleID, maxAttempts, lastErr)
}

func resourceCloudStackRolePermissionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	roleID := d.Get("role_id").(string)

	p := cs.Role.NewListRolePermissionsParams()
	p.SetRoleid(roleID)
	r, err := cs.Role.ListRolePermissions(p)
	if err != nil {
		return fmt.Errorf("Error listing RolePermissions: %s", err)
	}

	for _, rp := range r.RolePermissions {
		if rp.Id == d.Id() {
			d.Set("role_id", rp.Roleid)
			d.Set("rule", rp.Rule)
			d.Set("permission", rp.Permission)
			d.Set("description", rp.Description)
			return nil
		}
	}

	log.Printf("[DEBUG] RolePermission %s no longer exists", d.Id())
	d.SetId("")
	return nil
}

func resourceCloudStackRolePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("permission") {
		p := cs.Role.NewUpdateRolePermissionParams(d.Get("role_id").(string))
		p.SetRuleid(d.Id())
		p.SetPermission(d.Get("permission").(string))
		if _, err := cs.Role.UpdateRolePermission(p); err != nil {
			return fmt.Errorf("Error updating RolePermission: %s", err)
		}
	}

	return resourceCloudStackRolePermissionRead(d, meta)
}

func resourceCloudStackRolePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Role.NewDeleteRolePermissionParams(d.Id())
	log.Printf("[DEBUG] Deleting RolePermission %s", d.Id())
	if _, err := cs.Role.DeleteRolePermission(p); err != nil {
		return fmt.Errorf("Error deleting RolePermission: %s", err)
	}

	return nil
}
