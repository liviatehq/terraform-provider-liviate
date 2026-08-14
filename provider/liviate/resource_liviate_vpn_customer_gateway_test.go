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

func TestAccCloudStackVPNCustomerGateway_basic(t *testing.T) {
	var vpnCustomerGateway cloudstack.VpnCustomerGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"liviate_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "name", "terraform-foo"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "name", "terraform-bar"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "esp_policy", "aes256-sha1"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "ike_policy", "aes256-sha1;modp1536"),
				),
			},
		},
	})
}

func TestAccCloudStackVPNCustomerGateway_update(t *testing.T) {
	var vpnCustomerGateway cloudstack.VpnCustomerGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"liviate_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "name", "terraform-foo"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "name", "terraform-bar"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "esp_policy", "aes256-sha1"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "ike_policy", "aes256-sha1;modp1536"),
				),
			},

			{
				Config: testAccCloudStackVPNCustomerGateway_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"liviate_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayUpdatedAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "name", "terraform-foo-bar"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "name", "terraform-bar-foo"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.bar", "esp_policy", "3des-md5"),
					resource.TestCheckResourceAttr(
						"liviate_vpn_customer_gateway.foo", "ike_policy", "3des-md5;modp1536"),
				),
			},
		},
	})
}

func TestAccCloudStackVPNCustomerGateway_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
			},

			{
				ResourceName:      "liviate_vpn_customer_gateway.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackVPNCustomerGatewayExists(
	n string, vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN CustomerGateway ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		v, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN CustomerGateway not found")
		}

		*vpnCustomerGateway = *v

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayAttributes(
	vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpnCustomerGateway.Esppolicy != "aes256-sha1" {
			return fmt.Errorf("Bad ESP policy: %s", vpnCustomerGateway.Esppolicy)
		}

		if vpnCustomerGateway.Ikepolicy != "aes256-sha1;modp1536" {
			return fmt.Errorf("Bad IKE policy: %s", vpnCustomerGateway.Ikepolicy)
		}

		if vpnCustomerGateway.Ipsecpsk != "terraform" {
			return fmt.Errorf("Bad IPSEC pre-shared key: %s", vpnCustomerGateway.Ipsecpsk)
		}

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayUpdatedAttributes(
	vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpnCustomerGateway.Esppolicy != "3des-md5" {
			return fmt.Errorf("Bad ESP policy: %s", vpnCustomerGateway.Esppolicy)
		}

		if vpnCustomerGateway.Ikepolicy != "3des-md5;modp1536" {
			return fmt.Errorf("Bad IKE policy: %s", vpnCustomerGateway.Ikepolicy)
		}

		if vpnCustomerGateway.Ipsecpsk != "terraform" {
			return fmt.Errorf("Bad IPSEC pre-shared key: %s", vpnCustomerGateway.Ipsecpsk)
		}

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "liviate_vpn_customer_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Customer Gateway ID is set")
		}

		_, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Customer Gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackVPNCustomerGateway_basic = `
resource "liviate_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.1.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_vpc" "bar" {
  name = "terraform-vpc"
  cidr = "10.2.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_vpn_gateway" "foo" {
	vpc_id = liviate_vpc.foo.id
}

resource "liviate_vpn_gateway" "bar" {
	vpc_id = liviate_vpc.bar.id
}

resource "liviate_vpn_customer_gateway" "foo" {
	name = "terraform-foo"
	cidr = liviate_vpc.foo.cidr
	esp_policy = "aes256-sha1"
	gateway = liviate_vpn_gateway.foo.public_ip
	ike_policy = "aes256-sha1;modp1536"
	ipsec_psk = "terraform"
}

resource "liviate_vpn_customer_gateway" "bar" {
  name = "terraform-bar"
  cidr = liviate_vpc.bar.cidr
  esp_policy = "aes256-sha1"
  gateway = liviate_vpn_gateway.bar.public_ip
  ike_policy = "aes256-sha1;modp1536"
	ipsec_psk = "terraform"
}`

const testAccCloudStackVPNCustomerGateway_update = `
resource "liviate_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.1.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_vpc" "bar" {
  name = "terraform-vpc"
  cidr = "10.2.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "liviate_vpn_gateway" "foo" {
  vpc_id = liviate_vpc.foo.id
}

resource "liviate_vpn_gateway" "bar" {
  vpc_id = liviate_vpc.bar.id
}

resource "liviate_vpn_customer_gateway" "foo" {
  name = "terraform-foo-bar"
  cidr = liviate_vpc.foo.cidr
  esp_policy = "3des-md5"
  gateway = liviate_vpn_gateway.foo.public_ip
  ike_policy = "3des-md5;modp1536"
  ipsec_psk = "terraform"
}

resource "liviate_vpn_customer_gateway" "bar" {
  name = "terraform-bar-foo"
  cidr = liviate_vpc.bar.cidr
  esp_policy = "3des-md5"
  gateway = liviate_vpn_gateway.bar.public_ip
  ike_policy = "3des-md5;modp1536"
  ipsec_psk = "terraform"
}`
