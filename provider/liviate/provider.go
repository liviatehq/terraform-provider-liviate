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
	"errors"

	"github.com/go-ini/ini"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("LIVIATE_API_URL", "https://console.liviate.com/client/api"),
				ConflictsWith: []string{"config", "profile"},
			},

			"api_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("LIVIATE_API_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
				Sensitive:     true,
			},

			"secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("LIVIATE_SECRET_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
				Sensitive:     true,
			},

			"config": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},

			"profile": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},

			"http_get_only": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIVIATE_HTTP_GET_ONLY", false),
			},

			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIVIATE_TIMEOUT", 900),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"liviate_autoscale_policy":          dataSourceCloudstackAutoscalePolicy(),
			"liviate_autoscale_vm_group":        dataSourceCloudstackAutoscaleVMGroup(),
			"liviate_autoscale_vm_profile":      dataSourceCloudstackAutoscaleVMProfile(),
			"liviate_condition":                 dataSourceCloudstackCondition(),
			"liviate_counter":                   dataSourceCloudstackCounter(),
			"liviate_template":                  dataSourceCloudstackTemplate(),
			"liviate_ssh_keypair":               dataSourceCloudstackSSHKeyPair(),
			"liviate_instance":                  dataSourceCloudstackInstance(),
			"liviate_network_offering":          dataSourceCloudstackNetworkOffering(),
			"liviate_zone":                      dataSourceCloudStackZone(),
			"liviate_service_offering":          dataSourceCloudstackServiceOffering(),
			"liviate_volume":                    dataSourceCloudstackVolume(),
			"liviate_vpc":                       dataSourceCloudstackVPC(),
			"liviate_ipaddress":                 dataSourceCloudstackIPAddress(),
			"liviate_user":                      dataSourceCloudstackUser(),
			"liviate_vpn_connection":            dataSourceCloudstackVPNConnection(),
			"liviate_pod":                       dataSourceCloudstackPod(),
			"liviate_domain":                    dataSourceCloudstackDomain(),
			"liviate_project":                   dataSourceCloudstackProject(),
			"liviate_physical_network":          dataSourceCloudStackPhysicalNetwork(),
			"liviate_role":                      dataSourceCloudstackRole(),
			"liviate_cluster":                   dataSourceCloudstackCluster(),
			"liviate_limits":                    dataSourceCloudStackLimits(),
			"liviate_quota":                     dataSourceCloudStackQuota(),
			"liviate_quota_enabled":             dataSourceCloudStackQuotaEnabled(),
			"liviate_quota_tariff":              dataSourceCloudStackQuotaTariff(),
			"liviate_user_data":                 dataSourceCloudstackUserData(),
			"liviate_kubernetes_cluster_config": dataSourceCloudStackKubernetesClusterConfig(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"liviate_affinity_group":                 resourceCloudStackAffinityGroup(),
			"liviate_attach_volume":                  resourceCloudStackAttachVolume(),
			"liviate_autoscale_policy":               resourceCloudStackAutoScalePolicy(),
			"liviate_bucket":                         resourceCloudStackBucket(),
			"liviate_autoscale_vm_group":             resourceCloudStackAutoScaleVMGroup(),
			"liviate_autoscale_vm_profile":           resourceCloudStackAutoScaleVMProfile(),
			"liviate_cni_configuration":              resourceCloudStackCniConfiguration(),
			"liviate_condition":                      resourceCloudStackCondition(),
			"liviate_configuration":                  resourceCloudStackConfiguration(),
			"liviate_counter":                        resourceCloudStackCounter(),
			"liviate_cluster":                        resourceCloudStackCluster(),
			"liviate_disk":                           resourceCloudStackDisk(),
			"liviate_egress_firewall":                resourceCloudStackEgressFirewall(),
			"liviate_firewall":                       resourceCloudStackFirewall(),
			"liviate_host":                           resourceCloudStackHost(),
			"liviate_instance":                       resourceCloudStackInstance(),
			"liviate_ipaddress":                      resourceCloudStackIPAddress(),
			"liviate_kubernetes_cluster":             resourceCloudStackKubernetesCluster(),
			"liviate_kubernetes_version":             resourceCloudStackKubernetesVersion(),
			"liviate_loadbalancer_rule":              resourceCloudStackLoadBalancerRule(),
			"liviate_network":                        resourceCloudStackNetwork(),
			"liviate_network_acl":                    resourceCloudStackNetworkACL(),
			"liviate_network_acl_rule":               resourceCloudStackNetworkACLRule(),
			"liviate_nic":                            resourceCloudStackNIC(),
			"liviate_physical_network":               resourceCloudStackPhysicalNetwork(),
			"liviate_pod":                            resourceCloudStackPod(),
			"liviate_port_forward":                   resourceCloudStackPortForward(),
			"liviate_network_service_provider_state": resourceCloudStackNetworkServiceProviderState(),
			"liviate_private_gateway":                resourceCloudStackPrivateGateway(),
			"liviate_secondary_ipaddress":            resourceCloudStackSecondaryIPAddress(),
			"liviate_secondary_storage":              resourceCloudStackSecondaryStorage(),
			"liviate_security_group":                 resourceCloudStackSecurityGroup(),
			"liviate_security_group_rule":            resourceCloudStackSecurityGroupRule(),
			"liviate_ssh_keypair":                    resourceCloudStackSSHKeyPair(),
			"liviate_static_nat":                     resourceCloudStackStaticNAT(),
			"liviate_static_route":                   resourceCloudStackStaticRoute(),
			"liviate_storage_pool":                   resourceCloudStackStoragePool(),
			"liviate_template":                       resourceCloudStackTemplate(),
			"liviate_traffic_type":                   resourceCloudStackTrafficType(),
			"liviate_vpc":                            resourceCloudStackVPC(),
			"liviate_vpn_connection":                 resourceCloudStackVPNConnection(),
			"liviate_vpn_customer_gateway":           resourceCloudStackVPNCustomerGateway(),
			"liviate_vpn_gateway":                    resourceCloudStackVPNGateway(),
			"liviate_network_offering":               resourceCloudStackNetworkOffering(),
			"liviate_disk_offering":                  resourceCloudStackDiskOffering(),
			"liviate_vlan_ip_range":                  resourceCloudstackVlanIpRange(),
			"liviate_volume":                         resourceCloudStackVolume(),
			"liviate_zone":                           resourceCloudStackZone(),
			"liviate_service_offering":               resourceCloudStackServiceOffering(),
			"liviate_account":                        resourceCloudStackAccount(),
			"liviate_project":                        resourceCloudStackProject(),
			"liviate_user":                           resourceCloudStackUser(),
			"liviate_domain":                         resourceCloudStackDomain(),
			"liviate_network_service_provider":       resourceCloudStackNetworkServiceProvider(),
			"liviate_role":                           resourceCloudStackRole(),
			"liviate_role_permission":                resourceCloudStackRolePermission(),
			"liviate_limits":                         resourceCloudStackLimits(),
			"liviate_snapshot_policy":                resourceCloudStackSnapshotPolicy(),
			"liviate_quota_tariff":                   resourceCloudStackQuotaTariff(),
			"liviate_user_data":                      resourceCloudStackUserData(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (any, error) {
	apiURL, apiURLOK := d.GetOk("api_url")
	apiKey, apiKeyOK := d.GetOk("api_key")
	secretKey, secretKeyOK := d.GetOk("secret_key")
	config, configOK := d.GetOk("config")
	profile, profileOK := d.GetOk("profile")

	switch {
	case apiURLOK, apiKeyOK, secretKeyOK:
		if !(apiURLOK && apiKeyOK && secretKeyOK) {
			return nil, errors.New("'api_url', 'api_key' and 'secret_key' should all have values")
		}
	case configOK, profileOK:
		if !(configOK && profileOK) {
			return nil, errors.New("'config' and 'profile' should both have a value")
		}
	default:
		return nil, errors.New(
			"either 'api_url', 'api_key' and 'secret_key' or 'config' and 'profile' should have values")
	}

	if configOK && profileOK {
		cfg, err := ini.Load(config.(string))
		if err != nil {
			return nil, err
		}

		section, err := cfg.GetSection(profile.(string))
		if err != nil {
			return nil, err
		}

		apiURL = section.Key("url").String()
		apiKey = section.Key("apikey").String()
		secretKey = section.Key("secretkey").String()
	}

	cfg := Config{
		APIURL:      apiURL.(string),
		APIKey:      apiKey.(string),
		SecretKey:   secretKey.(string),
		HTTPGETOnly: d.Get("http_get_only").(bool),
		Timeout:     int64(d.Get("timeout").(int)),
	}

	return cfg.NewClient()
}
