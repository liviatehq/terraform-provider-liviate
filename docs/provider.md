# Provider Configuration

Configure the Liviate provider with a block like:

```hcl
provider "liviate" {
  api_key    = var.liviate_api_key
  secret_key = var.liviate_secret_key
}
```

The API URL defaults to `https://console.liviate.com/client/api` inside the
provider, so only the key and secret are usually needed.

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `api_url` | string | no | Liviate API URL. Defaults to `https://console.liviate.com/client/api`. Can also be sourced from `LIVIATE_API_URL`. |
| `api_key` | string | no* | Liviate API key. Can also be sourced from `LIVIATE_API_KEY`. |
| `secret_key` | string | no* | Liviate API secret. Can also be sourced from `LIVIATE_SECRET_KEY`. |
| `config` | string | no | Path to a CloudMonkey-style INI config file containing `url`, `apikey`, and `secretkey`. Conflicts with `api_url`/`api_key`/`secret_key`. |
| `profile` | string | no | Profile section in the `config` file. Must be set together with `config`. |
| `http_get_only` | bool | no | Only issue HTTP GET calls (for providers that reject POST). Default `false`. Env: `LIVIATE_HTTP_GET_ONLY`. |
| `timeout` | int | no | Timeout in seconds for async jobs. Default `900`. Env: `LIVIATE_TIMEOUT`. |

\* `api_key`/`secret_key` are required unless you use `config` + `profile`.

## Environment variables

| Variable | Maps to |
|----------|---------|
| `LIVIATE_API_URL` | `api_url` |
| `LIVIATE_API_KEY` | `api_key` |
| `LIVIATE_SECRET_KEY` | `secret_key` |
| `LIVIATE_HTTP_GET_ONLY` | `http_get_only` |
| `LIVIATE_TIMEOUT` | `timeout` |

## Authentication

Use either:

- `api_url` + `api_key` + `secret_key` (all three), or
- `config` + `profile` pointing at an INI file.

Mixing the two sets of options is not allowed.
