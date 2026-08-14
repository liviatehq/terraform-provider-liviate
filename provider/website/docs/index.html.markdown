---
layout: "liviate"
page_title: "Provider: Liviate"
sidebar_current: "docs-liviate-index"
description: |-
  The Liviate provider is used to interact with the many resources supported by Liviate. The provider needs to be configured with a URL pointing to a running Liviate API and the proper credentials before it can be used.
---

# Liviate Provider

The Liviate provider is used to interact with the many resources
supported by Liviate. The provider needs to be configured with a
URL pointing to a running Liviate API and the proper credentials
before it can be used.

In order to provide the required configuration options you can either
supply values for the `api_url`, `api_key` and `secret_key` fields, or
for the `config` and `profile` fields. A combination of both is not
allowed and will not work.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Liviate Provider
provider "liviate" {
  api_url    = "${var.liviate_api_url}"
  api_key    = "${var.liviate_api_key}"
  secret_key = "${var.liviate_secret_key}"
}

# Create a web server
resource "liviate_instance" "web" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `api_url` - (Optional) This is the Liviate API URL. It can also be sourced
  from the `LIVIATE_API_URL` environment variable.

* `api_key` - (Optional) This is the Liviate API key. It can also be sourced
  from the `LIVIATE_API_KEY` environment variable.

* `secret_key` - (Optional) This is the Liviate secret key. It can also be
  sourced from the `LIVIATE_SECRET_KEY` environment variable.

* `config` - (Optional) The path to a `CloudMonkey` config file. If set the API
  URL, key and secret will be retrieved from this file.

* `profile` - (Optional) Used together with the `config` option. Specifies which
  `CloudMonkey` profile in the config file to use.

* `http_get_only` - (Optional) Some cloud providers only allow HTTP GET calls to
  their Liviate API. If using such a provider, you need to set this to `true`
  in order for the provider to only make GET calls and no POST calls. It can also
  be sourced from the `LIVIATE_HTTP_GET_ONLY` environment variable.

* `timeout` - (Optional) A value in seconds. This is the time allowed for Liviate
  to complete each asynchronous job triggered. If unset, this can be sourced from the
  `LIVIATE_TIMEOUT` environment variable. Otherwise, this will default to 300
  seconds.
