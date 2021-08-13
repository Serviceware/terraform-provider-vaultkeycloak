package vaultkeycloak

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVaultKeycloakBasic(t *testing.T) {

	cleanup := dockerSetup(t)
	defer cleanup()

	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	os.Setenv("VAULT_TOKEN", "root")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVaultKeycloakSecretBackendConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVaultKeycloakSecretBackendExists("vaultkeycloak_secret_backend.test_backend"),
				),
			},
		},
	})
}

func testAccCheckVaultKeycloakSecretBackendConfigBasic() string {
	return `
	resource "vaultkeycloak_secret_backend" "test_backend" {
		client_id     = "vault"
		client_secret = "vault"
		server_url    = "http://keycloak:8080"
		realm         = "master"
		path          = "keycloak-secrets"
	  }
	`
}

func testAccCheckVaultKeycloakSecretBackendExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No backend set")
		}

		return nil
	}
}
