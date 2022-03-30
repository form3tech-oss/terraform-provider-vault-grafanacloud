# Replace default with your vault address
variable "vault_addr" {
  default = "http://localhost:8200"
}

# Replace deefault with your vault admin token
variable "vault_token" {
  default = "token"
}

# Replace default with your Grafana Cloud Admin API Key
variable "grafana_cloud_api_key" {
  default = "api-key"
}

variable "grafana_cloud_api_url" {
  default = "https://grafana.com/api"
}

# Replace default with your Grafana Cloud Organisation Slug
variable "grafana_cloud_org" {
  default = "my-org"
}

# Replace default with your Grafana Cloud User
variable "grafana_cloud_user" {
  default = "my-user"
}
