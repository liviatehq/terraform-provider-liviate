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
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackLoadBalancerRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_update(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// Skip this test on CloudStack 4.22.0.0 due to a known simulator bug
			// that causes "530 Internal Server Error" when updating load balancer rules.
			// This bug does not exist in 4.20.1.0, 4.22.1.0+, or 4.23.0.0+.
			version := getCloudStackVersion(t)
			if version == "4.22.0.0" {
				t.Skip("Skipping TestAccCloudStackLoadBalancerRule_update on CloudStack 4.22.0.0 due to known simulator bug (Error 530: Internal Server Error)")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_forceNew(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_forcenew,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "443"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "protocol", "tcp-proxy"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "cidrlist.0", "20.0.0.0/8"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_vpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_vpcUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_vpc_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("liviate_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "private_port", "443"),
					resource.TestCheckResourceAttr(
						"liviate_loadbalancer_rule.foo", "cidrlist.0", "20.0.0.0/8"),
				),
			},
		},
	})
}

func testAccCheckCloudStackLoadBalancerRuleExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer rule ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		_, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Loadbalancer rule %s not found", n)
		}

		return nil
	}
}

func testAccCheckCloudStackLoadBalancerRuleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "liviate_loadbalancer_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.LoadBalancer.GetLoadBalancerRuleByID(id)
			if err == nil {
				return fmt.Errorf("Loadbalancer rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackLoadBalancerRule_basic = `
resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "liviate_ipaddress" "foo" {
  network_id = liviate_network.foo.id
}

resource "liviate_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = liviate_ipaddress.foo.id
  algorithm = "roundrobin"
  public_port = 80
  private_port = 80
  member_ids = [liviate_instance.foobar1.id]
}`

const testAccCloudStackLoadBalancerRule_update = `
resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "liviate_ipaddress" "foo" {
  network_id = liviate_network.foo.id
}

resource "liviate_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = liviate_ipaddress.foo.id
  algorithm = "leastconn"
  public_port = 80
  private_port = 80
  member_ids = [liviate_instance.foobar1.id]
}`

const testAccCloudStackLoadBalancerRule_forcenew = `
resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "liviate_ipaddress" "foo" {
  network_id = liviate_network.foo.id
}

resource "liviate_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = liviate_ipaddress.foo.id
  algorithm = "leastconn"
  public_port = 443
  private_port = 443
  protocol = "tcp-proxy"
  member_ids = [liviate_instance.foobar1.id]
  cidrlist = ["20.0.0.0/8"]
}`

const testAccCloudStackLoadBalancerRule_vpc = `
resource "liviate_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = liviate_vpc.foo.id
  zone = liviate_vpc.foo.zone
}

resource "liviate_ipaddress" "foo" {
  vpc_id = liviate_vpc.foo.id
  zone = liviate_vpc.foo.zone
}

resource "liviate_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = liviate_ipaddress.foo.id
  algorithm = "roundrobin"
  network_id = liviate_network.foo.id
  public_port = 80
  private_port = 80
  member_ids = [liviate_instance.foobar1.id]
}`

const testAccCloudStackLoadBalancerRule_vpc_update = `
resource "liviate_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = liviate_vpc.foo.id
  zone = liviate_vpc.foo.zone
}

resource "liviate_ipaddress" "foo" {
  vpc_id = liviate_vpc.foo.id
  zone = liviate_vpc.foo.zone
}

resource "liviate_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_instance" "foobar2" {
  name = "terraform-server2"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = liviate_ipaddress.foo.id
  algorithm = "leastconn"
  network_id = liviate_network.foo.id
  public_port = 443
  private_port = 443
  member_ids = [liviate_instance.foobar1.id, liviate_instance.foobar2.id]
  cidrlist = ["20.0.0.0/8"]
}`
