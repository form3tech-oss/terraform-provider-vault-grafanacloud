resource "vaultgrafanacloud_secret_role" "test" {
  backend         = "grafanacloud"
  name            = "my-role"
  gc_role         = "Viewer"
  ttl_seconds     = "3600"
  max_ttl_seconds = "3600"
}
