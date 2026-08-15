# Publishing the Liviate Terraform Provider

This guide covers publishing a release to the public Terraform Registry so
users can install it with:

```hcl
terraform {
  required_providers {
    liviate = {
      source  = "liviatehq/liviate"
      version = ">= 1.0.0"
    }
  }
}
```

## Prerequisites

- A [GPG key](https://gnupg.org/) for signing releases
- Write access to `github.com/liviatehq/terraform-provider-liviate`
- [GoReleaser](https://goreleaser.com/install/) installed locally

## Step 1 — Generate a GPG signing key

```sh
gpg --full-generate-key
```

Pick RSA 4096, no expiration. Export the fingerprint:

```sh
gpg --list-secret-keys --keyid-format LONG
# Look for the line:   sec   rsa4096/XXXXXXXXXXXX
# The XXXXXXXXXXXX part is your fingerprint
```

Export the public key (you'll upload this to the Terraform Registry):

```sh
gpg --export --armor XXXXXXXXXXXX > liviate-gpg.pub
```

## Step 2 — Set environment variables

```sh
export GPG_FINGERPRINT=XXXXXXXXXXXX
export GITHUB_TOKEN=ghp_...   # a GitHub token with repo scope
```

## Step 3 — Tag and release

```sh
git tag v1.0.0
goreleaser release --clean
```

GoReleaser will:
- Build the provider for all platforms (linux/darwin/windows/freebsd/openbsd × amd64/386/arm64/arm)
- Create SHA256SUMS and sign them with your GPG key
- Create a GitHub Release on `github.com/liviatehq/terraform-provider-liviate`
- Upload the zip archives and signed checksums

## Step 4 — Publish on the Terraform Registry

1. Go to https://registry.terraform.io
2. Sign in with your GitHub account (must have admin access to `liviatehq/terraform-provider-liviate`)
3. Click **Publish** → **Provider**
4. Select `liviatehq/terraform-provider-liviate`
5. Paste the armored GPG public key (`cat liviate-gpg.pub`)
6. Agree to the terms and submit
7. The registry validates the release and makes it live at `registry.terraform.io/providers/liviatehq/liviate`

## Step 5 — Verify

In a fresh directory, create a `main.tf` with the provider source and run `terraform init`. It should download from the registry.

## Syncing with upstream (apache/cloudstack-terraform-provider)

To fetch new commits from the Apache CloudStack provider without rebasing
(which isn't possible directly due to the repo restructure), use:

```sh
git fetch upstream
```

Then cherry-pick or merge individual upstream commits into `master` as needed.
Because the directory layout differs, you may need to adjust file paths
(e.g. `main.go` in upstream maps to `provider/main.go` here).

## Notes

- The provider address is `registry.terraform.io/liviatehq/liviate`; this is
  baked into `provider/main.go`.
- `.goreleaser.yml` references `owner: liviatehq`.
- Sample configurations in `sample/` use `source = "liviatehq/liviate"`.
