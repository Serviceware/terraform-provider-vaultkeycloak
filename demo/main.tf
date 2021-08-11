terraform {
  required_providers {
    vaultkeycloak = {
      version = "0.1.0"
      source  = "Serviceware/vaultkeycloak"
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
  depends_on = [
    keycloak_realm.demo
  ]
  source = "../terraform/tfmodule-vaultkeycloak-config"
  realm           = keycloak_realm.demo.realm
  vault_client_id = local.keycloak_vault_client_id

}

##################
provider "vaultkeycloak" {
  vault_address = "http://127.0.0.1:8200"
  vault_token   = "root"
}
resource "vaultkeycloak_secret_backend" "default" {
  client_id     = local.keycloak_vault_client_id
  client_secret = "vault123"
  server_url    = "http://127.0.0.1:8080"
  realm         = keycloak_realm.demo.realm
  path          = "keycloak-secrets"
  provider      = vaultkeycloak

  depends_on = [
    module.keycloak_vault_config
  ]

}
