# Resources & Data Sources

The provider exposes **62 resources** and **27 data sources**. Resources manage
objects in Liviate; data sources are read-only lookups used to reference
existing objects (for example, finding a template or zone ID).

Full argument references are generated in the
[`provider/website/docs/`](../provider/website/docs/) directory
(`r/` for resources, `d/` for data sources).

## Resources

### Compute

| Resource | Description |
|----------|-------------|
| `liviate_instance` | Deploys and manages a virtual machine |
| `liviate_attach_volume` | Attaches a volume to a virtual machine |
| `liviate_ssh_keypair` | Manages an SSH keypair |
| `liviate_user_data` | Manages user data (cloud-init) content |

### Storage & Object Storage

| Resource | Description |
|----------|-------------|
| `liviate_volume` | Manages a block storage volume |
| `liviate_disk` | Manages a custom disk |
| `liviate_disk_offering` | Manages a disk offering |
| `liviate_storage_pool` | Manages a primary storage pool |
| `liviate_secondary_storage` | Manages a secondary storage provider |
| `liviate_snapshot_policy` | Manages a volume snapshot policy |
| `liviate_bucket` | Manages an S3-compatible object storage bucket |

### Networking

| Resource | Description |
|----------|-------------|
| `liviate_network` | Manages a network |
| `liviate_network_acl` | Manages a network ACL |
| `liviate_network_acl_rule` | Manages network ACL rules |
| `liviate_network_offering` | Manages a network offering |
| `liviate_network_service_provider` | Manages a network service provider |
| `liviate_network_service_provider_state` | Manages the enabled state of a network service provider |
| `liviate_nic` | Manages a NIC on a virtual machine |
| `liviate_ipaddress` | Manages a public IP address |
| `liviate_secondary_ipaddress` | Manages a secondary public IP address |
| `liviate_vlan_ip_range` | Manages a VLAN IP range |
| `liviate_loadbalancer_rule` | Manages a load balancer rule |
| `liviate_port_forward` | Manages a port forwarding rule |
| `liviate_static_nat` | Manages a static NAT rule |
| `liviate_firewall` | Manages a firewall rule |
| `liviate_egress_firewall` | Manages an egress firewall rule |
| `liviate_private_gateway` | Manages a private gateway |
| `liviate_static_route` | Manages a static route |
| `liviate_vpc` | Manages a VPC |

### VPN

| Resource | Description |
|----------|-------------|
| `liviate_vpn_gateway` | Manages a VPN gateway |
| `liviate_vpn_customer_gateway` | Manages a VPN customer gateway |
| `liviate_vpn_connection` | Manages a VPN connection |

### Security

| Resource | Description |
|----------|-------------|
| `liviate_security_group` | Manages a security group |
| `liviate_security_group_rule` | Manages security group rules |

### Kubernetes

| Resource | Description |
|----------|-------------|
| `liviate_kubernetes_cluster` | Manages a Kubernetes cluster |
| `liviate_kubernetes_version` | Registers a Kubernetes version |
| `liviate_cni_configuration` | Manages a CNI configuration |

### Autoscaling

| Resource | Description |
|----------|-------------|
| `liviate_autoscale_policy` | Manages an autoscale policy |
| `liviate_autoscale_vm_group` | Manages an autoscale VM group |
| `liviate_autoscale_vm_profile` | Manages an autoscale VM profile |
| `liviate_condition` | Manages a condition used by autoscaling |
| `liviate_counter` | Manages a counter used by autoscaling |

### Offerings & Templates

| Resource | Description |
|----------|-------------|
| `liviate_service_offering` | Manages a service offering |
| `liviate_service_offering_constrained` | Manages a constrained service offering (Plugin Framework) |
| `liviate_service_offering_fixed` | Manages a fixed service offering (Plugin Framework) |
| `liviate_service_offering_unconstrained` | Manages an unconstrained service offering (Plugin Framework) |
| `liviate_template` | Manages a VM template |
| `liviate_affinity_group` | Manages an affinity group |

### IAM & Accounts

| Resource | Description |
|----------|-------------|
| `liviate_account` | Manages an account |
| `liviate_user` | Manages a user |
| `liviate_domain` | Manages a domain |
| `liviate_role` | Manages a role |
| `liviate_project` | Manages a project |

### Infrastructure

| Resource | Description |
|----------|-------------|
| `liviate_zone` | Manages a zone |
| `liviate_cluster` | Manages a cluster |
| `liviate_pod` | Manages a pod |
| `liviate_host` | Manages a host |
| `liviate_physical_network` | Manages a physical network |
| `liviate_traffic_type` | Manages a traffic type |
| `liviate_configuration` | Manages a global configuration value |

### Quotas & Limits

| Resource | Description |
|----------|-------------|
| `liviate_limits` | Manages resource limits |
| `liviate_quota_tariff` | Manages a quota tariff |

## Data Sources

| Data source | Description |
|-------------|-------------|
| `liviate_autoscale_policy` | Look up an autoscale policy |
| `liviate_autoscale_vm_group` | Look up an autoscale VM group |
| `liviate_autoscale_vm_profile` | Look up an autoscale VM profile |
| `liviate_condition` | Look up a condition |
| `liviate_counter` | Look up a counter |
| `liviate_template` | Look up a template |
| `liviate_ssh_keypair` | Look up an SSH keypair |
| `liviate_instance` | Look up a virtual machine |
| `liviate_network_offering` | Look up a network offering |
| `liviate_zone` | Look up a zone |
| `liviate_service_offering` | Look up a service offering |
| `liviate_volume` | Look up a volume |
| `liviate_vpc` | Look up a VPC |
| `liviate_ipaddress` | Look up a public IP address |
| `liviate_user` | Look up a user |
| `liviate_vpn_connection` | Look up a VPN connection |
| `liviate_pod` | Look up a pod |
| `liviate_domain` | Look up a domain |
| `liviate_project` | Look up a project |
| `liviate_physical_network` | Look up a physical network |
| `liviate_role` | Look up a role |
| `liviate_cluster` | Look up a cluster |
| `liviate_limits` | Look up resource limits |
| `liviate_quota` | Look up quota usage |
| `liviate_quota_enabled` | Check whether quota is enabled |
| `liviate_quota_tariff` | Look up a quota tariff |
| `liviate_user_data` | Look up user data |

## Notes

- `liviate_bucket`, `liviate_attach_volume`, `liviate_host`, and
  `liviate_configuration` do not yet have a full argument page under
  `provider/website/docs/`. For those, the resource schema in
  `provider/liviate/` is the source of truth.
