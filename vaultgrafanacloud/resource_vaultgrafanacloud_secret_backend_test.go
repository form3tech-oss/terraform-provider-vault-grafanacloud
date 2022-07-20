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

func TestGrafanaCloudSecretBackend(t *testing.T) {
	backend := acctest.RandomWithPrefix("tf-test-grafanacloud")
	key := uuid.New().String()
	url := "http://localhost"
	organisation := "test_org"
	user := "user"

	resource.Test(t, resource.TestCase{
		Providers:                 testProviders,
		PreCheck:                  func() { testutil.TestAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccGrafanaCloudSecretBackendCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGrafanaCloudSecretBackend_initialConfig(backend, key, url, organisation, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "backend", backend),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "key", key),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "url", url),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "organisation", organisation),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "user", user),
				),
			},
			{
				Config: testGrafanaCloudSecretBackend_updateConfig(backend, key, url, organisation, user),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "backend", backend),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "key", key),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "url", url),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "organisation", organisation),
					resource.TestCheckResourceAttr("vaultgrafanacloud_secret_backend.test", "user", user),
				),
			},
		},
	})
}

func testAccGrafanaCloudSecretBackendCheckDestroy(s *terraform.State) error {
	client := testProvider.Meta().(*api.Client)

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vaultgrafanacloud_secret_backend" {
			continue
		}
		for backend, mount := range mounts {
			backend = strings.Trim(backend, "/")
			rsBackend := strings.Trim(rs.Primary.Attributes["backend"], "/")
			if mount.Type == "grafanacloud" && backend == rsBackend {
				return fmt.Errorf("Mount %q still exists", rsBackend)
			}
		}
	}
	return nil
}

func testGrafanaCloudSecretBackend_initialConfig(backend, key, url, organisation, user string) string {
	return fmt.Sprintf(`
resource "vaultgrafanacloud_secret_backend" "test" {
	backend = "%s"
	key = "%s"
	url = "%s"
	organisation = "%s"
	user = "%s"
}`, backend, key, url, organisation, user)
}

func testGrafanaCloudSecretBackend_updateConfig(backend, key, url, organisation, user string) string {
	return fmt.Sprintf(`
resource "vaultgrafanacloud_secret_backend" "test" {
	backend = "%s"
	key = "%s"
	url = "%s"
	organisation = "%s"
	user = "%s"
}`, backend, key, url, organisation, user)
}
