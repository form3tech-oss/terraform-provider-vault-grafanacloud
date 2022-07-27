
provider "vault-grafanacloud" {
  # Can also be provided via VAULT_ADDR
  address = var.your_vault_addr

  # Can also be provided via VAULT_TOKEN
  token = var.your_secret_token
}