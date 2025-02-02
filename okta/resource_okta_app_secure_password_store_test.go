package okta

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccAppSecurePasswordStoreApplication_credsSchemes(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSecurePasswordStore)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", appSecurePasswordStore)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appSecurePasswordStore, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "ADMIN_SETS_CREDENTIALS"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
					resource.TestCheckResourceAttr(resourceName, "label", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "url", "http://test.com"),
					resource.TestCheckResourceAttr(resourceName, "username_field", "user"),
					resource.TestCheckResourceAttr(resourceName, "password_field", "pass"),
					resource.TestCheckResourceAttr(resourceName, "credentials_scheme", "EXTERNAL_PASSWORD_SYNC"),
				),
			},
		},
	})
}

func TestAccAppSecurePasswordStoreApplication_timeouts(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(appSecurePasswordStore)
	resourceName := fmt.Sprintf("%s.test", appSecurePasswordStore)
	config := `
resource "okta_app_secure_password_store" "test" {
  label              = "testAcc_replace_with_uuid"
  username_field     = "user"
  password_field     = "pass"
  url                = "http://test.com"
  credentials_scheme = "ADMIN_SETS_CREDENTIALS"
  timeouts {
    create = "60m"
    read = "2h"
    update = "30m"
  }
}`
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(appSecurePasswordStore, createDoesAppExist(okta.NewSecurePasswordStoreApplication())),
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(config, ri),
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, createDoesAppExist(okta.NewAutoLoginApplication())),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "60m"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.read", "2h"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.update", "30m"),
				),
			},
		},
	})
}
