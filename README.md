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