# Terraform Provider for Vault Keycloak engine

This terraform provider allows you to configure the Vault Keycloak plugin

https://github.com/Serviceware/vault-plugin-secrets-keycloak

## Demo

Nav into demo
```
cd demo
```


Start keycloak+vault

```
docker-compose -f ../testing/docker-compose.yaml up -d --build
```

Enable keycloak
```
export VAULT_ADDR="http://127.0.0.1:8200"
vault secrets enable -path=keycloak-secrets vault-plugin-secrets-keycloak
```

Apply config
```
terraform init
terraform apply
```

Read a client secret
```
vault read keycloak-secrets/client-secret/test-client
```

## Using this in an existing project

At timse it is useful to upgrade an existing project using the vaultkeycloak provider
to a locally developed version, for example to test fixes or new features.

This can be done using [Developer Overrides](https://developer.hashicorp.com/terraform/cli/v1.1.x/config/config-file#development-overrides-for-provider-developers).

If you haven't done so, build the provider so the built `terraform-provider-keycloak`
binary is in the checkout of the repository.

```
make build
```

Create or open a `~/.terraformrc` file on linux-based systems, and add the following
configuration:

```
provider_installation {
  dev_overrides {
    "Serviceware/vaultkeycloak" = "/path/to/your/checkouf/of/terraform-provider-vaultkeycloak"
  }
}
```

With this in place, the next terraform runs should greet you with a large warning banner:

```
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - serviceware/vaultkeycloak in /path/to/your/checkouf/of/terraform-provider-vaultkeycloak
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
```

Now you can iterate between `make build` and `terraform plan` and `terraform apply` as necessary
to validate your local changes.
