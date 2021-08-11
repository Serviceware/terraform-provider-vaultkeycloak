resource "vaultkeycloak_secret_backend" "demo" {
  client_id     = "vault"
  client_secret = "secret123"
  server_url    = "http://127.0.0.1:8080"
  realm         = "my-realm"
  path          = "keycloak-secrets"

}
