package vaultkeycloak

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVaultKeycloakBasicPerRealm(t *testing.T) {

	cleanup := dockerSetup(t)
	defer cleanup()

	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	os.Setenv("VAULT_TOKEN", "root")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVaultKeycloakSecretBackendConfigBasicPerRealm("master", "vault", "vault"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVaultKeycloakSecretBackendExists("vaultkeycloak_secret_backend_per_realm_config.test_config_for_master"),
					resource.TestCheckResourceAttr("vaultkeycloak_secret_backend_per_realm_config.test_config_for_master", "client_id", "vault"),
				),
			},
		},
	})
}
func TestAccVaultKeycloakBasicPerRealmIgnoreConnect(t *testing.T) {

	cleanup := dockerSetup(t)
	defer cleanup()

	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	os.Setenv("VAULT_TOKEN", "root")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVaultKeycloakSecretBackendConfigBasicPerRealmIgnoreConnect("master", "vault", "thisisnottheclientsecret"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVaultKeycloakSecretBackendExists("vaultkeycloak_secret_backend_per_realm_config.test_config_for_master"),
					resource.TestCheckResourceAttr("vaultkeycloak_secret_backend_per_realm_config.test_config_for_master", "client_id", "vault"),
				),
			},
		},
	})
}

func testAccCheckVaultKeycloakSecretBackendConfigBasicPerRealm(realm, client_id, client_secret string) string {
	return fmt.Sprintf(`
	resource "vaultkeycloak_secret_backend_per_realm_config" "test_config_for_master" {
		client_id     = "%s"
		client_secret = "%s"
		server_url    = "http://keycloak:8080/auth"
		realm         = "%s"
		path          = "keycloak-secrets"
	  }
	`, client_id, client_secret, realm)
}
func testAccCheckVaultKeycloakSecretBackendConfigBasicPerRealmIgnoreConnect(realm, client_id, client_secret string) string {
	return fmt.Sprintf(`
	resource "vaultkeycloak_secret_backend_per_realm_config" "test_config_for_master" {
		client_id     = "%s"
		client_secret = "%s"
		server_url    = "http://keycloak:8080/auth"
		realm         = "%s"
		path          = "keycloak-secrets"
		ignore_connectivity_check = true
	  }
	`, client_id, client_secret, realm)
}
