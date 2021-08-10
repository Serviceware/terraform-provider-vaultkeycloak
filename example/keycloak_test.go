package example

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var testProvider *schema.Provider
var testProviders map[string]terraform.ResourceProvider

func init() {
	testProvider = Provider()
	testProviders = map[string]terraform.ResourceProvider{
		"vault": testProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VAULT_ADDR"); v == "" {
		t.Fatal("VAULT_ADDR must be set for acceptance tests")
	}
	if v := os.Getenv("VAULT_TOKEN"); v == "" {
		t.Fatal("VAULT_TOKEN must be set for acceptance tests")
	}
}

// example.Widget represents a concrete Go type that represents an API resource
func TestAccExampleWidget_basic(t *testing.T) {
	var widgetBefore, widgetAfter string

	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExampleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceExists("example_widget.foo", &widgetBefore),
				),
			},
			{
				Config: testAccExampleResource_removedPolicy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceExists("example_widget.foo", &widgetAfter),
				),
			},
		},
	})
}
