
provider "vault-grafanacloud" {
  address = var.vault_addr
  token   = var.vault_token
}

resource "vaultgrafanacloud_secret_backend" "backend" {
  backend      = "vault-grafanacloud"
  key          = var.grafana_cloud_api_key
  url          = var.grafana_cloud_api_url
  organisation = var.grafana_cloud_org
  user         = var.grafana_cloud_user
}

resource "vaultgrafanacloud_secret_role" "test" {
  backend = "vault-grafanacloud"
  name    = "viewer-role"
  gc_role = "Viewer"
  ttl     = "3600"
  max_ttl = "3600"
}
