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
	"encoding/base64"
	"fmt"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"
)

// Wraps getKubernetesClusterConfig -- lets a Terraform config chain straight from
// liviate_kubernetes_cluster into the standard hashicorp/kubernetes (or helm) provider without a
// manual kubectl step, e.g.:
//
//	data "liviate_kubernetes_cluster_config" "this" { cluster_id = liviate_kubernetes_cluster.this.id }
//	provider "kubernetes" {
//	  host                   = data.liviate_kubernetes_cluster_config.this.host
//	  cluster_ca_certificate = base64decode(data.liviate_kubernetes_cluster_config.this.cluster_ca_certificate)
//	  client_certificate     = base64decode(data.liviate_kubernetes_cluster_config.this.client_certificate)
//	  client_key             = base64decode(data.liviate_kubernetes_cluster_config.this.client_key)
//	}
//
// The individual fields are extracted from the raw kubeconfig's default context (CKS always
// returns a single-cluster, single-user, single-context kubeconfig, client-certificate auth) --
// raw_config is also exposed as an escape hatch for anything unusual.
func dataSourceCloudStackKubernetesClusterConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudStackKubernetesClusterConfigRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"raw_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The full kubeconfig YAML as returned by CloudStack.",
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Base64-encoded -- pass through base64decode() into the kubernetes/helm provider.",
			},
			"client_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

// Minimal shape of the fields we need from a CKS kubeconfig -- not a full kubeconfig schema.
type kubeconfigShape struct {
	Clusters []struct {
		Cluster struct {
			Server                   string `yaml:"server"`
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func dataSourceCloudStackKubernetesClusterConfigRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	clusterID := d.Get("cluster_id").(string)

	p := cs.Kubernetes.NewGetKubernetesClusterConfigParams()
	p.SetId(clusterID)
	r, err := cs.Kubernetes.GetKubernetesClusterConfig(p)
	if err != nil {
		return fmt.Errorf("Error getting Kubernetes cluster config: %s", err)
	}

	var cfg kubeconfigShape
	if err := yaml.Unmarshal([]byte(r.Configdata), &cfg); err != nil {
		return fmt.Errorf("Error parsing kubeconfig for cluster %s: %s", clusterID, err)
	}
	if len(cfg.Clusters) == 0 || len(cfg.Users) == 0 {
		return fmt.Errorf("kubeconfig for cluster %s did not contain the expected clusters/users sections", clusterID)
	}

	d.SetId(clusterID)
	d.Set("raw_config", r.Configdata)
	d.Set("host", cfg.Clusters[0].Cluster.Server)
	// Already base64 in the source YAML -- re-encode ensures a consistent, decodable value even
	// if a future CKS version ever returns these fields any other way.
	d.Set("cluster_ca_certificate", reencodeBase64(cfg.Clusters[0].Cluster.CertificateAuthorityData))
	d.Set("client_certificate", reencodeBase64(cfg.Users[0].User.ClientCertificateData))
	d.Set("client_key", reencodeBase64(cfg.Users[0].User.ClientKeyData))

	return nil
}

func reencodeBase64(alreadyEncoded string) string {
	raw, err := base64.StdEncoding.DecodeString(alreadyEncoded)
	if err != nil {
		return alreadyEncoded // pass through as-is rather than fail the whole read
	}
	return base64.StdEncoding.EncodeToString(raw)
}
