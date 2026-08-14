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
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackNIC_basic(t *testing.T) {
	var nic cloudstack.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNICExists(
						"liviate_instance.foobar", "liviate_nic.foo", &nic),
					testAccCheckCloudStackNICAttributes(&nic),
				),
			},
		},
	})
}

func TestAccCloudStackNIC_update(t *testing.T) {
	var nic cloudstack.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNICExists(
						"liviate_instance.foobar", "liviate_nic.foo", &nic),
					testAccCheckCloudStackNICAttributes(&nic),
				),
			},

			{
				Config: testAccCloudStackNIC_ipaddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNICExists(
						"liviate_instance.foobar", "liviate_nic.foo", &nic),
					testAccCheckCloudStackNICIPAddress(&nic),
					resource.TestCheckResourceAttr(
						"liviate_nic.foo", "ip_address", "10.1.2.123"),
				),
			},
		},
	})
}

func TestAccCloudStackNIC_macaddress(t *testing.T) {
	var nic cloudstack.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNIC_macaddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNICExists(
						"liviate_instance.foobar", "liviate_nic.foo", &nic),
					testAccCheckCloudStackNICMacAddress(&nic),
					resource.TestCheckResourceAttr(
						"liviate_nic.foo", "mac_address", "02:1a:4b:3c:5d:6e"),
				),
			},
		},
	})
}

func testAccCheckCloudStackNICExists(
	v, n string, nic *cloudstack.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsv, ok := s.RootModule().Resources[v]
		if !ok {
			return fmt.Errorf("Not found: %s", v)
		}

		if rsv.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		rsn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rsn.Primary.ID == "" {
			return fmt.Errorf("No NIC ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rsv.Primary.ID)

		if err != nil {
			return err
		}

		for _, n := range vm.Nic {
			if n.Id == rsn.Primary.ID {
				*nic = n
				return nil
			}
		}

		return fmt.Errorf("NIC not found")
	}
}

func testAccCheckCloudStackNICAttributes(
	nic *cloudstack.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkname != "terraform-network-secondary" {
			return fmt.Errorf("Bad network name: %s", nic.Networkname)
		}

		return nil
	}
}

func testAccCheckCloudStackNICIPAddress(
	nic *cloudstack.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkname != "terraform-network-secondary" {
			return fmt.Errorf("Bad network name: %s", nic.Networkname)
		}

		if nic.Ipaddress != "10.1.2.123" {
			return fmt.Errorf("Bad IP address: %s", nic.Ipaddress)
		}

		return nil
	}
}

func testAccCheckCloudStackNICMacAddress(
	nic *cloudstack.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkname != "terraform-network-secondary" {
			return fmt.Errorf("Bad network name: %s", nic.Networkname)
		}

		if nic.Macaddress != "02:1a:4b:3c:5d:6e" {
			return fmt.Errorf("Bad MAC address: %s", nic.Macaddress)
		}

		return nil
	}
}

func testAccCheckCloudStackNICDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	// Deleting the instance automatically deletes any additional NICs
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "liviate_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Virtual Machine %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackNIC_basic = `
resource "liviate_network" "foo" {
  name = "terraform-network-primary"
  display_text = "terraform-network-primary"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_network" "bar" {
  name = "terraform-network-secondary"
  display_text = "terraform-network-secondary"
  cidr = "10.1.2.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Medium Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "liviate_nic" "foo" {
  network_id = liviate_network.bar.id
  virtual_machine_id = liviate_instance.foobar.id
}`

const testAccCloudStackNIC_ipaddress = `
resource "liviate_network" "foo" {
  name = "terraform-network-primary"
  display_text = "terraform-network-primary"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_network" "bar" {
  name = "terraform-network-secondary"
  display_text = "terraform-network-secondary"
  cidr = "10.1.2.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Medium Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "liviate_nic" "foo" {
  network_id = liviate_network.bar.id
  virtual_machine_id = liviate_instance.foobar.id
  ip_address = "10.1.2.123"
}`

const testAccCloudStackNIC_macaddress = `
resource "liviate_network" "foo" {
  name = "terraform-network-primary"
  display_text = "terraform-network-primary"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_network" "bar" {
  name = "terraform-network-secondary"
  display_text = "terraform-network-secondary"
  cidr = "10.1.2.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Medium Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "liviate_nic" "foo" {
  network_id = liviate_network.bar.id
  virtual_machine_id = liviate_instance.foobar.id
  mac_address = "02:1a:4b:3c:5d:6e"
}`
