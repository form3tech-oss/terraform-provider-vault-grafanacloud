package vaultgrafanacloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/form3tech-oss/terraform-provider-vault-grafanacloud/testutil"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/vault/api"
)

func TestGrafanaCloudSecretRole(t *testing.T) {
	backend := acctest.RandomWithPrefix("tf-test-grafanacloud")
	key := uuid.New().String()
	url := "http://localhost"
	organisation := "test_org"
	user := "user"
	name := uuid.New().String()
	gcRole := "Viewer"
	updatedGCRole := "Admin"
	ttl := "1"
	updatedTTL := "2"
	maxTTL := "2"
	updatedMaxTTL := "3"

	resource.Test(t, resource.TestCase{
		Providers:    testProviders,
		PreCheck:     func() { testutil.TestAccPreCheck(t) },
		CheckDestroy: testAccGrafanaCloudSecretRoleCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGrafanaCloudSecretRole_initialConfig(backend, key, url, organisation, user, name, gcRole, ttl, maxTTL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "backend", backend),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "name", name),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "gc_role", gcRole),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "ttl_seconds", ttl),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "max_ttl_seconds", maxTTL),
				),
			},
			{
				Config: testGrafanaCloudSecretRole_updateConfig(backend, key, url, organisation, name, user, updatedGCRole, updatedTTL, updatedMaxTTL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "name", name),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "gc_role", updatedGCRole),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "ttl_seconds", updatedTTL),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_role.test", "max_ttl_seconds", updatedMaxTTL),
				),
			},
		},
	})
}

func testAccGrafanaCloudSecretRoleCheckDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vaultgrafanacloud_secret_role" {
			continue
		}
		for backend, mount := range mounts {
			backend = strings.Trim(backend, "/")
			rs := strings.Trim(rs.Primary.Attributes["backend"], "/")
			if mount.Type == "grafanacloud" && backend == rs {
				return fmt.Errorf("Mount %q still exists", rs)
			}
		}
	}
	return nil
}

func testGrafanaCloudSecretRole_initialConfig(backend, key, url, organisation, user, name, gcRole, ttl, maxTTL string) string {
	return fmt.Sprintf(`
resource "vaultgrafanacloud_secret_backend" "test" {
	backend = "%s"
	key = "%s"
	url = "%s"
	organisation = "%s"
	user = "%s"
}

resource "vaultgrafanacloud_secret_role" "test" {
	backend = vaultgrafanacloud_secret_backend.test.backend
	name = "%s"
	gc_role = "%s"
	ttl_seconds = %v
	max_ttl_seconds = %v
}
`, backend, key, url, organisation, user, name, gcRole, ttl, maxTTL)
}

func testGrafanaCloudSecretRole_updateConfig(backend, key, url, organisation, user, name, gcRole, ttl, maxTTL string) string {
	return fmt.Sprintf(`
resource "vaultgrafanacloud_secret_backend" "test" {
	backend = "%s"
	key = "%s"
	url = "%s"
	organisation = "%s"
	user = "%s"
}

resource "vaultgrafanacloud_secret_role" "test" {
	backend = vaultgrafanacloud_secret_backend.test.backend
	name = "%s"
	gc_role = "%s"
	ttl_seconds = %v
	max_ttl_seconds = %v
}`, backend, key, url, organisation, name, user, gcRole, ttl, maxTTL)
}
