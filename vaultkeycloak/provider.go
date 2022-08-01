package vaultkeycloak

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/config"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vault_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"add_address_to_env": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "If true, adds the value of the `address` argument to the Terraform process environment.",
			},
			"vault_token": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vaultkeycloak_secret_backend": resourceKeycloakSecretBackend(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

// Taken from:
// https://github.com/hashicorp/terraform-provider-vault/blob/e02ab326651869963bdeb5769e9eb95d3ce8b81b/internal/provider/meta.go#L377
func GetToken(d *schema.ResourceData) (string, error) {
	if token := d.Get("vault_token").(string); token != "" {
		return token, nil
	}

	if addAddr := d.Get("add_address_to_env").(string); addAddr == "true" {
		if addr := d.Get("vault_address").(string); addr != "" {
			addrEnvVar := "VAULT_ADDR"
			if current, exists := os.LookupEnv(addrEnvVar); exists {
				defer func() {
					os.Setenv(addrEnvVar, current)
				}()
			} else {
				defer func() {
					os.Unsetenv(addrEnvVar)
				}()
			}
			if err := os.Setenv(addrEnvVar, addr); err != nil {
				return "", err
			}
		}
	}

	// Use ~/.vault-token, or the configured token helper.
	tokenHelper, err := config.DefaultTokenHelper()
	if err != nil {
		return "", fmt.Errorf("error getting token helper: %s", err)
	}
	token, err := tokenHelper.Get()
	if err != nil {
		return "", fmt.Errorf("error getting token: %s", err)
	}
	return strings.TrimSpace(token), nil
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	vault_address := d.Get("vault_address").(string)
	vault_token := d.GetToken()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	config := api.DefaultConfig()
	if vault_address != "" {
		config.Address = vault_address
	}
	client, err := api.NewClient(config)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Vault client",
			Detail:   "Unable connect Vault client",
		})
		return nil, diags
	}
	if vault_token != "" {
		client.SetToken(vault_token)
	}

	return client, diags
}
