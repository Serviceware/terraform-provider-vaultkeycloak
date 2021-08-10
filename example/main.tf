terraform {
  required_providers {
    vaultkeycloak = {
      version = "0.1.0"
      source  = "github.com/Serviceware/vault-keycloak"
    }

    keycloak = {
      source  = "mrparkers/keycloak"
      version = "3.3.0"
    }
  }
}


provider "keycloak" {
  client_id = "admin-cli"
  username  = "admin"
  password  = "admin"
  url       = "http://localhost:8080"
}

locals {
  keycloak_vault_client_id = "vault"
}

resource "keycloak_realm" "demo" {

  realm   = "demo"
  enabled = true

}
module "keycloak_vault_config" {

  source = "../terraform/tfmodule-vault-keycloak-config"
  realm           = keycloak_realm.realm
  vault_client_id = locals.keycloak_vault_client_id

}

##################
provider "vaultkeycloak" {
  vault_address = "http://127.0.0.1:8200"
  vault_token   = "root"
}
resource "vault_keycloak_secret_backend" "default" {
  client_id     = locals.keycloak_vault_client_id
  client_secret = "vault123"
  server_url    = "http://127.0.0.1:8080"
  realm         = "my-realm"
  path          = "keycloak-secrets"
  provider      = vaultkeycloak

  depends_on = [
    module.keycloak_vault_config
  ]

}
