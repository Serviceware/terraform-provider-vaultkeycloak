---
page_title: "vaultkeycloak_secret_backend_per_realm_config Resource - terraform-provider-vaultkeycloak"
subcategory: ""
description: |-
  The vaultkeycloak_secret_backend resource allows you to configure the secret backend engine.
---

# Resource `vaultkeycloak_secret_backend_per_realm_config`

The vaultkeycloak_secret_backend_per_realm_config resource allows you to configure a secret backend engine for a specific realm.

This is useful if different realm have different server urls or different client ids.

## Example Usage

```terraform
resource "vaultkeycloak_secret_backend_per_realm_config" "my_realm_config" {
  client_id     = "vault"
  client_secret = "secret123"
  server_url    = "http://127.0.0.1:8080/auth"
  realm         = "my-realm"
  path          = "keycloak-secrets"

  ignore_connectivity_check = true # optional
}
```

## Argument Reference

- `client_id` - (Required) The client id used to access the client secrets
- `client_secret` - (Required) The client secret for the given client id
- `server_url` - (Required) The server url to the keycloak server. For older keycloaks this means the server url usually ends with `/auth`
- `realm` - (Required) The realm from which the secrets should be read
- `path` - (Required) The path under which the engine is registered
- `ignore_connectivity_check` - (Optional) If set to true, the plugin will not check the connectivity to the keycloak server. This is useful if you want to use the plugin in a vault cluster that is not able to reach the keycloak server. Defaults to false.


## Attributes Reference
