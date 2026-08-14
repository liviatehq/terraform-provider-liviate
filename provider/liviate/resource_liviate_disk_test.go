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

func TestAccCloudStackDisk_basic(t *testing.T) {
	var disk cloudstack.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDisk_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskExists(
						"liviate_disk.foo", &disk),
					testAccCheckCloudStackDiskAttributes(&disk),
					testAccCheckResourceTags(&disk),
				),
			},
		},
	})
}

func TestAccCloudStackDisk_update(t *testing.T) {
	var disk cloudstack.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDisk_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskExists(
						"liviate_disk.foo", &disk),
					testAccCheckCloudStackDiskAttributes(&disk),
				),
			},

			{
				Config: testAccCloudStackDisk_resize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskExists(
						"liviate_disk.foo", &disk),
					testAccCheckCloudStackDiskResized(&disk),
					resource.TestCheckResourceAttr(
						"liviate_disk.foo", "disk_offering", "Medium"),
				),
			},
		},
	})
}

func TestAccCloudStackDisk_deviceID(t *testing.T) {
	var disk cloudstack.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDisk_deviceID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskExists(
						"liviate_disk.foo", &disk),
					testAccCheckCloudStackDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"liviate_disk.foo", "device_id", "4"),
				),
			},
		},
	})
}

func TestAccCloudStackDisk_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDisk_basic,
			},

			{
				ResourceName:            "liviate_disk.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shrink_ok", "reattach_on_change"},
			},
		},
	})
}

func testAccCheckCloudStackDiskExists(
	n string, disk *cloudstack.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		volume, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if volume.Id != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *volume

		return nil
	}
}

func testAccCheckCloudStackDiskAttributes(
	disk *cloudstack.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Name != "terraform-disk" {
			return fmt.Errorf("Bad name: %s", disk.Name)
		}

		if disk.Diskofferingname != "Small" {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCloudStackDiskResized(
	disk *cloudstack.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Diskofferingname != "Medium" {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCloudStackDiskDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "liviate_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		_, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Disk %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackDisk_basic = `
resource "liviate_disk" "foo" {
  name = "terraform-disk"
  attach = false
  disk_offering = "Small"
  zone = "Sandbox-simulator"
  tags = {
    terraform-tag = "true"
  }
}`

const testAccCloudStackDisk_update = `
resource "liviate_disk" "foo" {
  name = "terraform-disk"
  disk_offering = "Small"
  zone = "Sandbox-simulator"
}`

const testAccCloudStackDisk_resize = `
resource "liviate_disk" "foo" {
  name = "terraform-disk"
  disk_offering = "Medium"
  zone = "Sandbox-simulator"
}`

const testAccCloudStackDisk_deviceID = `
resource "liviate_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "liviate_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = liviate_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = liviate_network.foo.zone
  expunge = true
}

resource "liviate_disk" "foo" {
  name = "terraform-disk"
  attach = true
  device_id = 4
  disk_offering = "Small"
  virtual_machine_id = liviate_instance.foobar.id
  zone = liviate_instance.foobar.zone
}`
