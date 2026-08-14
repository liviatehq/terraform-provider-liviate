# Sample: Bucket

Provisions S3-compatible object storage buckets using the `liviate_bucket`
resource. The sample includes **two** buckets:

- **Private** (`bucket.tf`) — default access control; only the owning account
  can read/write.
- **Public** (`public.tf`) — a bucket policy that grants public read access to
  objects.

## Files

| File | Purpose |
|------|---------|
| `main.tf` | Provider + `required_providers` block |
| `variables.tf` | Input variables (credentials are `sensitive`) |
| `bucket.tf` | Private bucket resource |
| `public.tf` | Public bucket resource (with a bucket policy) |
| `outputs.tf` | Bucket URL, access key, and (sensitive) secret key |
| `terraform.tfvars.example` | Copy to `terraform.tfvars` and set real values |

## Access control

Access to a bucket is controlled through:

1. **Ownership** — `account`, `domain_id`, and `project_id` arguments assign
   the bucket to a specific account/domain/project (default: the provider's own
   account).
2. **Bucket policy** — the `policy` argument takes an S3-style bucket policy
   JSON. The public example grants `s3:GetObject` to `Principal = "*"`.

## Credentials

After `terraform apply`, use the outputs to reach the bucket with any
S3-compatible client:

```sh
# from `terraform output`
url        = https://s3.liviate.com/<bucket-name>
access_key = <accesskey>
secret_key = <usersecretkey>   # sensitive
```

The `accesskey` and `usersecretkey` attributes are populated from the
CloudStack API on the `liviate_bucket` resource.

## Run it

```sh
cd sample/bucket
cp terraform.tfvars.example terraform.tfvars   # then fill in your API key/secret
terraform init
terraform plan
terraform apply
```

## Notes

- `objectstorageid` is optional. When omitted, the provider auto-discovers the
  single object storage pool configured in the environment. If you have more
  than one pool, set `objectstorageid` explicitly.
- The `policy` JSON in `public.tf` is a standard S3 public-read policy; adjust
  the `Resource` ARN and actions to match your CloudStack S3 backend if needed.
- Buckets are exposed **path-style** (`https://s3.liviate.com/<bucket>`), not
  virtual-hosted style (`<bucket>.s3.liviate.com`). S3 clients that default to
  virtual-hosted style must be forced to path-style (e.g. `s3_use_path_style =
  true` on the `hashicorp/aws` provider).
