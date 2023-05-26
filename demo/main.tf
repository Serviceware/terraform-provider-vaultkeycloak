terraform {
  required_providers {
 vaultkeycloak = {
      source = "Serviceware/vaultkeycloak"
      version = "0.2.0"
    }

    keycloak = {
      source  = "mrparkers/keycloak"
      version = "4.2.0"
    }
  }
}


provider "keycloak" {
  client_id = "admin-cli"
  username  = "admin"
  password  = "admin"
  url       = "http://localhost:8080"
  base_path = "/auth"
}

locals {
  keycloak_vault_client_id = "vault"
  keycloak_vault_client_secret = "secret123"
}

resource "keycloak_realm" "demo" {

  realm   = "demo"
  enabled = true

}
module "keycloak_vault_config" {
  depends_on = [
    keycloak_realm.demo
  ]

  source  = "Serviceware/keycloak-client/vaultkeycloak"
  version = "0.1.1"

  realm           = keycloak_realm.demo.realm
  vault_client_id = local.keycloak_vault_client_id
  vault_client_secret = local.keycloak_vault_client_secret

}

##################
provider "vaultkeycloak" {
  vault_address = "http://127.0.0.1:8200"
  vault_token   = "root"
}
resource "vaultkeycloak_secret_backend" "default" {
  client_id     = local.keycloak_vault_client_id
  client_secret = local.keycloak_vault_client_secret
  server_url    = "http://keycloak:8080/auth"
  realm         = keycloak_realm.demo.realm
  path          = "keycloak-secrets"


  depends_on = [
    module.keycloak_vault_config
  ]

}



#### test client
resource "keycloak_openid_client" "test_client" {
  realm_id            = keycloak_realm.demo.realm
  client_id           = "test-client"

  name                = "test client"
  enabled             = true

  access_type         = "CONFIDENTIAL"

}