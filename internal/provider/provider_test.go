// SPECIAL: Environment variables
package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/tfbrew/terraform-provider-aap/internal/configprefix"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	configprefix.Prefix: providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AAP_HOST"); v == "" {
		t.Fatal("AAP_HOST must be set for acceptance tests")
	}
	if os.Getenv("AAP_OAUTH_TOKEN") == "" && os.Getenv("AAP_TOKEN") == "" && (os.Getenv("AAP_USERNAME") == "" || os.Getenv("AAP_PASSWORD") == "") {
		t.Fatal("AAP_OAUTH_TOKEN, AAP_TOKEN, or AAP_USERNAME/AAP_PASSWORD must be set for acceptance tests")
	}
}
