package vaultkeycloak

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vault_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_token": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vault_keycloak_secret_backend": resourceKeycloakSecretBackend(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	vault_address := d.Get("vault_address").(string)
	vault_token := d.Get("vault_token").(string)

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
