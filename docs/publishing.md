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

## Step 1 — Generate a GPG signing key

```sh
gpg --full-generate-key
```

Pick RSA 4096, no expiration. Get the fingerprint:

```sh
gpg --list-secret-keys --keyid-format LONG
# Look for the line:   sec   rsa4096/XXXXXXXXXXXX
```

Export the **private** key (for CI signing):

```sh
gpg --armor --export-secret-keys XXXXXXXXXXXX > liviate-release-key.private.asc
```

Export the **public** key (for the Terraform Registry):

```sh
gpg --armor --export XXXXXXXXXXXX > liviate-gpg.pub
```

## Step 2 — Add the private key to GitHub secrets

1. Go to **Settings → Secrets and variables → Actions** on the
   `github.com/liviatehq/terraform-provider-liviate` repository.
2. Click **New repository secret**.
3. Name: `GPG_PRIVATE_KEY`
4. Value: paste the entire contents of `liviate-release-key.private.asc`
   (including the `-----BEGIN PGP PRIVATE KEY BLOCK-----` header).
5. Click **Add secret**.

The fingerprint (`XXXXXXXXXXXX`) is already hardcoded in the
`.github/workflows/release.yml` workflow file.

## Step 3 — Tag and push (triggers the CI release)

```sh
git tag v1.0.0
git push github v1.0.0
```

This triggers the `release` GitHub Actions workflow, which:
- Imports the GPG private key from the `GPG_PRIVATE_KEY` secret
- Runs GoReleaser to build all platform binaries, create a GitHub Release,
  and sign the SHA256SUMS with the GPG key

## Step 4 — Publish on the Terraform Registry

1. Go to https://registry.terraform.io
2. Sign in with your GitHub account (must have admin access to `liviatehq/terraform-provider-liviate`)
3. Click **Publish** → **Provider**
4. Select `liviatehq/terraform-provider-liviate`
5. Paste the armored GPG public key (`cat liviate-gpg.pub`)
6. Agree to the terms and submit
7. The registry validates the release and makes it live at `registry.terraform.io/providers/liviatehq/liviate`

## Step 5 — Verify

In a fresh directory, create a `main.tf` with the provider source and run
`terraform init`. It should download from the registry.

## Syncing with upstream (apache/cloudstack-terraform-provider)

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
