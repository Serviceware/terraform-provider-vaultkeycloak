package vaultkeycloak

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"vaultkeycloak": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func createTestVault(t *testing.T) (string, string, func()) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	return addr, rootToken, func() {
		ln.Close()
		core.Shutdown()

	}
}
func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("VAULT_ADDR"); err == "" {
		t.Fatal("VAULT_ADDR must be set for acceptance tests")
	}
	if err := os.Getenv("VAULT_TOKEN"); err == "" {
		t.Fatal("VAULT_TOKEN must be set for acceptance tests")
	}
}
