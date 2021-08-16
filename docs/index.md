---
page_title: "Provider: vaultkeycloak"
subcategory: ""
description: |-
  Terraform provider for to configure vault keycloak.
---

# Vault Keycloak Provider

A terraform provider to configure the connection for the vault keycloak engine

https://github.com/Serviceware/vault-plugin-secrets-keycloak


## Example Usage

You can pass in `vault_address` and `vault_token`, otherwise the internal vault api client will look for the environment variables `VAULT_ADDR` and `VAULT_TOKEN`

```terraform
provider "vaultkeycloak" {
  vault_address = "http://127.0.0.1:8200"
  vault_token   = "root"
}
```

## Schema

### Optional

- **vault_address** (String, Optional) Address of vault, otherwise `VAULT_ADDR` environment variable is used
- **vault_token** (String, Optional) vault token for auth, otherwise `VAULT_TOKEN` environment variable is used