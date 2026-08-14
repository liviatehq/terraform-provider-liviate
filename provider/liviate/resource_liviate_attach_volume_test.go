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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudstackAttachVolume_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudstackAttachVolume_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("liviate_attach_volume.foo", "device_id", "1"),
				),
			},
		},
	})
}

const testAccCloudstackAttachVolume_basic = `
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
	disk_offering = "Small"
	zone = liviate_instance.foobar.zone
}

resource "liviate_attach_volume" "foo" {
	volume_id          = liviate_disk.foo.id
	virtual_machine_id = liviate_instance.foobar.id
}
`
