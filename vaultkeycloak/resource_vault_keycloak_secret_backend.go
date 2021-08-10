package vaultkeycloak

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func resourceKeycloakSecretBackend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakSecretCreate,
		ReadContext:   resourceKeycloakSecretRead,
		UpdateContext: resourceKeycloakSecretUpdate,
		DeleteContext: resourceKeycloakSecretDelete,
		Schema: map[string]*schema.Schema{

			"server_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"realm": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "keycloak",
			},
		},
	}
}

func resourceKeycloakSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	client := m.(*api.Client)
	var diags diag.Diagnostics
	server_url := d.Get("server_url").(string)
	realm := d.Get("realm").(string)
	client_id := d.Get("client_id").(string)
	client_secret := d.Get("client_secret").(string)
	path := d.Get("path").(string)

	c := client.Logical()

	data := map[string]interface{}{
		"server_url":    server_url,
		"realm":         realm,
		"client_id":     client_id,
		"client_secret": client_secret,
	}
	_, err := c.Write(fmt.Sprintf("%s/config/connection", path), data)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(path)
	resourceKeycloakSecretRead(ctx, d, m)
	return diags
}

func resourceKeycloakSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*api.Client)
	c := client.Logical()
	path := d.Id()

	secret, err := c.Read(fmt.Sprintf("%s/config/connection", path))

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("realm", secret.Data["realm"])
	d.Set("server_url", secret.Data["server_url"])
	d.Set("client_id", secret.Data["client_id"])
	d.Set("client_secret", secret.Data["client_secret"])

	return diags
}

func resourceKeycloakSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceKeycloakSecretRead(ctx, d, m)
}

func resourceKeycloakSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*api.Client)
	c := client.Logical()
	path := d.Get("path").(string)

	_, err := c.Delete(fmt.Sprintf("%s/config/connection", path))

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}
