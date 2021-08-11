package vaultkeycloak

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVaultKeycloakBasic(t *testing.T) {

	os.Setenv("TF_ACC", "1") // needs to be set otherwise resource.Test doesn't do anything
	// addr, token, cleanup := createTestVault(t)
	// defer cleanup()

	// os.Setenv("VAULT_ADDR", addr)
	// os.Setenv("VAULT_TOKEN", token)
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
		client_secret = "6d1c7871-25a3-44d3-941e-d759fec65170"
		server_url    = "http://localhost:8080"
		realm         = "demo"
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
