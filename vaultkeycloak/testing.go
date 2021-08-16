package vaultkeycloak

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v8"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/api"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func dockerSetup(t *testing.T) func() {

	t.Helper()

	composeFilePaths := []string{"../testing/docker-compose.yaml"}
	identifier := strings.ToLower(uuid.New().String())

	compose := tc.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.
		WithCommand([]string{"up", "-d", "--build"}).
		WithExposedService("vault_1", 8200, wait.NewHTTPStrategy("/v1/sys/health").WithPort("8200/tcp")).
		WithExposedService("keycloak_1", 8080, wait.NewHTTPStrategy("/").WithPort("8080/tcp").WithStartupTimeout(3*time.Minute)).
		Invoke()

	if execError.Error != nil {
		defer compose.Down()
		t.Fatal(execError.Error)
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = "http://127.0.0.1:8200"

	vaultClient, err := api.NewClient(vaultConfig)
	if err != nil {
		defer compose.Down()
		t.Fatal(err)
	}
	vaultClient.SetToken("root")

	err = vaultClient.Sys().Mount("keycloak-secrets", &api.MountInput{
		Type: "vault-plugin-secrets-keycloak",
	})
	if err != nil {
		defer compose.Down()
		t.Fatal(err)
	}
	err = setupAdminClient("master", "vault", "vault")
	if err != nil {
		defer compose.Down()
		t.Fatal(err)
	}
	err = setupAdminClient("master", "vault2", "vault2")
	if err != nil {
		defer compose.Down()
		t.Fatal(err)
	}

	return func() {
		defer compose.Down()
	}

}
func setupAdminClient(realm, client_id, client_secret string) error {
	keycloakServerUrl := fmt.Sprintf("http://%s:%s", "127.0.0.1", "8080")
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
