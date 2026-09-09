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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackServiceOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackServiceOfferingRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.ServiceOffering.NewListServiceOfferingsParams()
	csServiceOfferings, err := cs.ServiceOffering.ListServiceOfferings(p)

	if err != nil {
		return fmt.Errorf("Failed to list service offerings: %s", err)
	}

	filters := d.Get("filter")
	var serviceOfferings []*cloudstack.ServiceOffering

	for _, s := range csServiceOfferings.ServiceOfferings {
		match, err := applyServiceOfferingFilters(s, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			serviceOfferings = append(serviceOfferings, s)
		}
	}

	if len(serviceOfferings) == 0 {
		return fmt.Errorf("No service offering is matching with the specified regex")
	}
	//return the latest service offering from the list of filtered service according
	//to its creation date
	serviceOffering, err := latestServiceOffering(serviceOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected service offerings: %s\n", serviceOffering.Displaytext)

	return serviceOfferingDescriptionAttributes(d, serviceOffering)
}

func serviceOfferingDescriptionAttributes(d *schema.ResourceData, serviceOffering *cloudstack.ServiceOffering) error {
	d.SetId(serviceOffering.Id)
	d.Set("name", serviceOffering.Name)
	d.Set("display_text", serviceOffering.Displaytext)

	return nil
}

func latestServiceOffering(serviceOfferings []*cloudstack.ServiceOffering) (*cloudstack.ServiceOffering, error) {
	var latest time.Time
	var serviceOffering *cloudstack.ServiceOffering

	for _, s := range serviceOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", s.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of an service offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			serviceOffering = s
		}
	}

	return serviceOffering, nil
}

func applyServiceOfferingFilters(serviceOffering *cloudstack.ServiceOffering, filters *schema.Set) (bool, error) {
	var serviceOfferingJSON map[string]interface{}
	k, _ := json.Marshal(serviceOffering)
	err := json.Unmarshal(k, &serviceOfferingJSON)
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
		serviceOfferingField := stringifyFilterValue(serviceOfferingJSON[updatedName])
		if !r.MatchString(serviceOfferingField) {
			return false, nil
		}

	}
	return true, nil
}

// stringifyFilterValue renders a decoded JSON field (string, float64, bool, or
// nil for a field that doesn't exist on this service offering, e.g. gpu-only
// fields on a non-GPU offering) as a string for filter regex matching. Without
// this, filtering on a numeric/boolean field (cpunumber, memory, gpudisplay,
// dynamicscalingenabled, ...) panics on the hard type assertion to string --
// found live 2026-09-09 trying to filter by cpunumber/memory to resolve a
// service offering by its stable specs instead of its rotating date-stamped
// name (see client-demos-infra GLPI Problem #53).
func stringifyFilterValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}
