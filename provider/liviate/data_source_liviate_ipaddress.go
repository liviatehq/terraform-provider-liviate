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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackIPAddress() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackIPAddressRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"is_portable": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				// Optional input (name or ID) so a caller can scope the lookup to a CloudStack
				// Project -- without it, resources living inside a Project are invisible to this
				// data source no matter what `filter` blocks are given (they only match
				// client-side against whatever the underlying list call already returned, and
				// that call excludes Project-scoped resources unless projectid is explicitly
				// passed; found live 2026-08-24, matches the WithProject() pattern already used
				// by every liviate_* RESOURCE in this provider). Still Computed so it keeps
				// working as an output attribute when left unset.
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_source_nat": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func datasourceCloudStackIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// `project` (name or ID) explicitly scopes to a CloudStack Project -- plain listall=true does
	// NOT surface Project-owned resources on its own, CloudStack requires the actual projectid
	// param (found live 2026-08-24: listall alone still returned zero results for a Project-scoped
	// IP that genuinely existed). Mirrors WithProject()'s name-or-ID resolution in cloudstack.go.
	// Resolved once, outside the retry loop below -- project-name-to-ID resolution isn't the thing
	// racing here.
	var projectID string
	if project, ok := d.GetOk("project"); ok {
		projectID = project.(string)
		if !cloudstack.IsID(projectID) {
			id, _, err := cs.Project.GetProjectID(projectID)
			if err != nil {
				return fmt.Errorf("Failed to resolve project %q: %s", projectID, err)
			}
			projectID = id
		}
	}

	filters := d.Get("filter").(*schema.Set)

	// GLPI Problem #61: a freshly-created network's SourceNat IP (or any other just-associated
	// public IP) is not always visible via listPublicIpAddresses the instant CloudStack's own API
	// call that created it returns -- the association settles asynchronously, server-side. This
	// data source has no notion of "which IP am I waiting for," only a client-side filter, so
	// unlike resource_liviate_network.go's own Create-time poll (which targets one specific
	// network's SourceNat IP), this retries the exact same list+filter this function already does
	// whenever it comes back with zero matches, up to 60s, before giving up -- covering every
	// caller of this data source (a `depends_on` on the freshly-created network resource is not by
	// itself enough: that resource's own Create only polls its own SourceNat association, not
	// whatever OTHER filter a caller's `data.liviate_ipaddress` block might be using).
	var publicIpAddress *cloudstack.PublicIpAddress
	ctx := context.Background()
	waitErr := retry.RetryContext(ctx, 60*time.Second, func() *retry.RetryError {
		p := cs.Address.NewListPublicIpAddressesParams()
		p.SetListall(true)
		if projectID != "" {
			p.SetProjectid(projectID)
		}
		csPublicIPAddresses, err := cs.Address.ListPublicIpAddresses(p)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("Failed to list ip addresses: %s", err))
		}

		var matches []*cloudstack.PublicIpAddress
		for _, ip := range csPublicIPAddresses.PublicIpAddresses {
			match, ferr := applyIPAddressFilters(ip, filters)
			if ferr != nil {
				return retry.NonRetryableError(ferr)
			}
			if match {
				matches = append(matches, ip)
			}
		}

		if len(matches) == 0 {
			log.Printf("[DEBUG] No ip address matching the specified filter yet, retrying...")
			return retry.RetryableError(fmt.Errorf("No ip address is matching with the specified regex"))
		}

		// return the latest ip address from the list of filtered ip addresses according
		// to its creation date
		latest, lerr := latestIPAddress(matches)
		if lerr != nil {
			return retry.NonRetryableError(lerr)
		}
		publicIpAddress = latest
		return nil
	})
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[DEBUG] Selected ip addresses: %s\n", publicIpAddress.Ipaddress)

	return ipAddressDescriptionAttributes(d, publicIpAddress)
}

func ipAddressDescriptionAttributes(d *schema.ResourceData, publicIpAddress *cloudstack.PublicIpAddress) error {
	d.SetId(publicIpAddress.Id)
	d.Set("is_portable", publicIpAddress.Isportable)
	d.Set("network_id", publicIpAddress.Networkid)
	d.Set("vpc_id", publicIpAddress.Vpcid)
	d.Set("zone_name", publicIpAddress.Zonename)
	d.Set("project", publicIpAddress.Project)
	d.Set("ip_address", publicIpAddress.Ipaddress)
	d.Set("is_source_nat", publicIpAddress.Issourcenat)
	d.Set("tags", tagsToMap(publicIpAddress.Tags))

	return nil
}

func latestIPAddress(publicIpAddresses []*cloudstack.PublicIpAddress) (*cloudstack.PublicIpAddress, error) {
	var latest time.Time
	var publicIpAddress *cloudstack.PublicIpAddress

	for _, ip := range publicIpAddresses {
		created, err := time.Parse("2006-01-02T15:04:05-0700", ip.Allocated)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse allocation date of the ip address: %s", err)
		}

		if created.After(latest) {
			latest = created
			publicIpAddress = ip
		}
	}

	return publicIpAddress, nil
}

func applyIPAddressFilters(publicIpAddress *cloudstack.PublicIpAddress, filters *schema.Set) (bool, error) {
	var publicIPAdressJSON map[string]interface{}
	k, _ := json.Marshal(publicIpAddress)
	err := json.Unmarshal(k, &publicIPAdressJSON)

	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		publicIPAdressField := fmt.Sprintf("%v", publicIPAdressJSON[updatedName])
		if !r.MatchString(publicIPAdressField) {
			return false, nil
		}
	}

	return true, nil
}
