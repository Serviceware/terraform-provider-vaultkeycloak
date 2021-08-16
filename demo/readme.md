
Start keycloak+vault

```
docker-compose up -d --build
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