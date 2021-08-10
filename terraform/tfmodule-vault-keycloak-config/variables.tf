variable "realm" {
  type = string
  default = "master"
  description = "The realm in which the vault client should be created"
}
variable "vault_client_id" {
  type = string
  default = "vault"

  description = "The name of the keycloak client that will be used by vault to access the client secrets"
}
