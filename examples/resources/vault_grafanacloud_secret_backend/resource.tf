resource "vault_grafanacloud_secret_backend" "backend" {
  backend      = "grafanacloud"
  key          = var.your_secret_api_key
  url          = "https://grafana.com/api"
  organisation = "my-org"
  user         = "my-user"
}
