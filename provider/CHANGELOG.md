## 1.0.0 (Unreleased)

### Added
- `liviate_bucket` resource for S3-compatible object storage buckets
  (object storage pool auto-discovery; `accesskey`/`usersecretkey` outputs)
- In-place root disk resizing for `liviate_instance`
- Sample configurations for VM, Kubernetes cluster, and object storage buckets
- Project documentation (getting started, provider config, resource reference, samples)

### Changed
- Rebranded to Liviate: namespace `liviatehq`, resources renamed
  `cloudstack_*` → `liviate_*`
- Default API endpoint: `https://console.liviate.com/client/api`
