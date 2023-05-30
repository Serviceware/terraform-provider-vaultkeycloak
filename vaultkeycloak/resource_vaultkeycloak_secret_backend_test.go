package vaultkeycloak

import (
	"fmt"
	"os"
	"testing"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"

	tcc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
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

func dockerSetup(t *testing.T) func() {

	t.Helper()

	identifier := tcc.StackIdentifier("vaultkeycloak-" + strings.ToLower(uuid.New().String()))
	ctx, cancel := context.WithCancel(context.Background())

	compose, err := tcc.NewDockerComposeWith(tcc.WithStackFiles("../testing/docker-compose.yaml"), identifier)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tcc.RemoveOrphans(true), tcc.RemoveImagesLocal), "compose.Down()")
	})
	t.Cleanup(cancel)

	execError := compose.
		WaitForService("vault", wait.NewHTTPStrategy("/v1/sys/health").WithPort("8200/tcp")).
		WaitForService("keycloak", wait.NewHTTPStrategy("/").WithPort("8080/tcp").WithStartupTimeout(3*time.Minute)).
		Up(ctx, tcc.Wait(true))

	if execError != nil {
		t.Fatalf("compose.Up() failed with %s", execError.Error())
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = "http://127.0.0.1:8200"

	vaultClient, err := api.NewClient(vaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	vaultClient.SetToken("root")

	err = vaultClient.Sys().Mount("keycloak-secrets", &api.MountInput{
		Type: "vault-plugin-secrets-keycloak",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = setupAdminClient("master", "vault", "vault")
	if err != nil {
		t.Fatal(err)
	}
	err = setupAdminClient("master", "vault2", "vault2")
	if err != nil {
		t.Fatal(err)
	}

	return func() {
		defer compose.Down(ctx, tcc.RemoveOrphans(true), tcc.RemoveImagesLocal)
	}

}
func setupAdminClient(realm, client_id, client_secret string) error {
	keycloakServerUrl := fmt.Sprintf("http://%s:%s/auth", "127.0.0.1", "8080")
	keycloakCLient := gocloak.NewClient(keycloakServerUrl)
	ctx := context.Background()
	loginToken, err := keycloakCLient.Login(ctx, "admin-cli", "", realm, "admin", "admin")
	if err != nil {
		return err
	}
	serviceAccountsEnabled := true
	createClientResponse, err := keycloakCLient.CreateClient(ctx, loginToken.AccessToken, realm, gocloak.Client{
		ID:                     &client_id,
		ClientID:               &client_id,
		Secret:                 &client_secret,
		ServiceAccountsEnabled: &serviceAccountsEnabled,
	})
	if err != nil {
		return err
	}
	clientServiceAccount, err := keycloakCLient.GetClientServiceAccount(ctx, loginToken.AccessToken, realm, createClientResponse)
	if err != nil {
		return err
	}
	adminRole, err := keycloakCLient.GetRealmRole(ctx, loginToken.AccessToken, realm, "admin")
	if err != nil {
		return err
	}
	err = keycloakCLient.AddRealmRoleToUser(ctx, loginToken.AccessToken, realm, *clientServiceAccount.ID, []gocloak.Role{*adminRole})
	if err != nil {
		return err
	}
	return nil
}
