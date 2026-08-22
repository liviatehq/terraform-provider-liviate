# Kubernetes cluster with working CSI storage

End-to-end example of what's needed to get a CKS cluster from `terraform apply` to "a plain PVC
just works": the cluster itself (`enable_csi = true`), the platform role-permission grant its CSI
driver needs but doesn't ship with, and a default `StorageClass` via the standard `kubernetes`
provider, chained straight off the cluster's own kubeconfig.

## Why the role-permission step exists

CloudStack's built-in `"Project Kubernetes Service Role"` (used internally by every CKS cluster's
CSI/CCM sidecar account) predates CSI support and has no volume-management API permissions at all.
Without granting them, every `PersistentVolumeClaim` fails to provision with:

```
API [listVolumes] does not exist or is not available for the account for user id [...]
```

This is a genuinely platform-wide gap -- fixing it once (this role isn't per-cluster or
per-project) fixes it for every cluster on the install, current and future.

## Why the disk offering has to be `iscustomized = true`

The CSI driver's provisioner honors whatever size a PVC actually requests (`storage: 5Gi`, etc.).
A disk offering with a fixed size can't do that -- only one with `iscustomized = true` lets the
driver create a volume at the exact requested size instead of snapping to (or rejecting) a
mismatched fixed size.

## Apply order

`terraform apply` handles the ordering on its own via the implicit dependency graph (data sources
read after their inputs exist) plus the explicit `depends_on` on the `StorageClass` -- there's
nothing manual to sequence. First run will take a few minutes: cluster bootstrap, then the role
grant, then the `kubernetes` provider can actually reach the API server to create the
`StorageClass`.
