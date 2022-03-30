resource "vault_grafanacloud_secret_role" "test" {
  backend = "grafanacloud"
  name    = "my-role"
  gc_role = "Viewer"
  ttl     = "1h"
  max_ttl = "1h"
}
