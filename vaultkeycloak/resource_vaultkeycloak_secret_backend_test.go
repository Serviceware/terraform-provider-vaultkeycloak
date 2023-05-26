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
				Config: testAccCheckVaultKeycloakSecretBackendConfigBasic("master", "vault", "vault"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVaultKeycloakSecretBackendExists("vaultkeycloak_secret_backend.test_backend"),
					resource.TestCheckResourceAttr("vaultkeycloak_secret_backend.test_backend", "client_id", "vault"),
				),
			},
			{
				Config: testAccCheckVaultKeycloakSecretBackendConfigBasic("master", "vault2", "vault2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vaultkeycloak_secret_backend.test_backend", "client_id", "vault2"),
				),
			},
		},
	})
}

func testAccCheckVaultKeycloakSecretBackendConfigBasic(realm, client_id, client_secret string) string {
	return fmt.Sprintf(`
	resource "vaultkeycloak_secret_backend" "test_backend" {
		client_id     = "%s"
		client_secret = "%s"
		server_url    = "http://keycloak:8080/auth"
		realm         = "%s"
		path          = "keycloak-secrets"
	  }
	`, client_id, client_secret, realm)
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
