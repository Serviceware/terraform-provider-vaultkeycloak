---
page_title: "vaultkeycloak_secret_backend Resource - terraform-provider-vaultkeycloak"
subcategory: ""
description: |-
  The vaultkeycloak_secret_backend resource allows you to configure the secret backend engine.
---

# Resource `hashicups_order`

The order resource allows you to configure a secret backend engine.

## Example Usage

```terraform
resource "vaultkeycloak_secret_backend" "my_realm_backend" {
  client_id     = "vault"
  client_secret = "secret123"
  server_url    = "http://127.0.0.1:8080"
  realm         = "my-realm"
  path          = "keycloak-secrets"

}
```

## Argument Reference

- `client_id` - (Required) The client id used to access the client secrets
- `client_secret` - (Required) The client secret for the given client id
- `server_url` - (Required) The server url to the keycloak server. We use [gocloak](https://github.com/Nerzal/gocloak) internally therefore the url should not contain `/auth` at the end
- `realm` - (Required) The realm from which the secrets should be read
- `path` - (Required) The path under which the engine is registered


## Attributes Reference
