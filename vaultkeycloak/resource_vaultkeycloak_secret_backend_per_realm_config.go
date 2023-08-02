package vaultkeycloak

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func resourceKeycloakSecretBackendPerRealmConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakSecretPerRealmConfigCreate,
		ReadContext:   resourceKeycloakSecretPerRealmConfigRead,
		UpdateContext: resourceKeycloakSecretPerRealmConfigUpdate,
		DeleteContext: resourceKeycloakSecretPerRealmConfigDelete,
		Schema: map[string]*schema.Schema{

			"server_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "keycloak",
				ForceNew: true,
			},
		},
	}
}

func resourceKeycloakSecretPerRealmConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type

	var diags diag.Diagnostics
	server_url := d.Get("server_url").(string)
	realm := d.Get("realm").(string)
	client_id := d.Get("client_id").(string)
	client_secret := d.Get("client_secret").(string)
	path := d.Get("path").(string)

	client := m.(*api.Client)
	c := client.Logical()

	data := map[string]interface{}{
		"server_url":    server_url,
		"client_id":     client_id,
		"client_secret": client_secret,
	}
	configPath := calcConfigPath(path, realm)
	_, err := c.Write(configPath, data)

	if err != nil {
		return diag.FromErr(err)
	}
	id := calcId(path, realm)
	d.SetId(id)
	resourceKeycloakSecretRead(ctx, d, m)
	return diags
}

func calcId(path string, realm string) string {
	id := fmt.Sprintf("%s/realms/%s", path, realm)
	return id
}
func pathAndRealmFromId(id string) (string, string) {
	var path string
	var realm string
	fmt.Sscanf(id, "%s/realms/%s", &path, &realm)
	return path, realm
}

func calcConfigPath(path string, realm string) string {
	configPath := fmt.Sprintf("%s/config/realms/%s/connection", path, realm)
	return configPath
}

func resourceKeycloakSecretPerRealmConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*api.Client)
	c := client.Logical()
	path, realm := pathAndRealmFromId(d.Id())

	secret, err := c.Read(calcConfigPath(path, realm))

	if err != nil {
		return diag.FromErr(err)
	}

	if secret == nil {
		d.SetId("")
		return diags
	}

	d.Set("server_url", secret.Data["server_url"])
	d.Set("client_id", secret.Data["client_id"])
	d.Set("client_secret", secret.Data["client_secret"])

	return diags
}

func resourceKeycloakSecretPerRealmConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*api.Client)

	path, realm := pathAndRealmFromId(d.Id())

	if d.HasChangesExcept("path", "realm") {

		server_url := d.Get("server_url").(string)
		client_id := d.Get("client_id").(string)
		client_secret := d.Get("client_secret").(string)

		c := client.Logical()

		data := map[string]interface{}{
			"server_url":    server_url,
			"client_id":     client_id,
			"client_secret": client_secret,
		}
		_, err := c.Write(calcConfigPath(path, realm), data)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceKeycloakSecretPerRealmConfigRead(ctx, d, m)
}

func resourceKeycloakSecretPerRealmConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*api.Client)
	c := client.Logical()
	path := d.Get("path").(string)
	realm := d.Get("realm").(string)

	_, err := c.Delete(calcConfigPath(path, realm))

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}
