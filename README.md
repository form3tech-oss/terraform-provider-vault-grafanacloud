# terraform-provider-vault-grafanacloud

[![Build Status](https://travis-ci.com/form3tech-oss/terraform-provider-vault-grafanacloud.svg?branch=master)](https://travis-ci.com/form3tech-oss/terraform-provider-vault-grafanacloud)

A Terraform provider for provisioning a [Grafana Cloud Secrets](https://github.com/form3tech-oss/vault-plugin-secrets-grafanacloud) backend and roles in a Vault instance.

## Installation

Download the relevant binary from [releases](https://github.com/form3tech-oss/terraform-provider-vault-grafanacloud/releases) and copy it to `$HOME/.terraform.d/plugins/`.

## Configuration

The following provider block variables are available for configuration:

| Name | Environment Variable | Description |
| ---- | -------------------- | ----------- |
| `address` | `VAULT_ADDR` | URL of the root of the target Vault server. |
| `token` | `VAULT_TOKEN` | Token to use to authenticate to Vault. |

Alternatively, these values can be read from the environment variables in the table.

## Resources

### `vaultgrafanacloud_secret_backend`

The `vaultgrafanacloud_secret_backend` resource mounts the [vault-plugin-secrets-grafanacloud](https://github.com/form3tech-oss/vault-plugin-secrets-grafanacloud) plugin to Vault.

#### Attributes

| Name | Required | Description | Default Value | 
| ---- | -------- | ----------- | ------------- |
| `backend` | `false` | The mount path for a backend, for example, the path given in "$ vault secrets enable -path=grafana-cloud grafana-cloud-plugin". | `grafana-cloud` |
| `key` | `true` | Grafana Cloud API key with Admin role to create user keys | N/A |
| `url` | `true` | The URL for the Grafana Cloud API | N/A |
| `organisation` | `true` | The Organisation slug for the Grafana Cloud API" | N/A |
| `user` | `true` | The User that is needed to interact with prometheus, if set this is returned alongside every issued credential | N/A |

### `vaultgrafanacloud_secret_role`

The `vaultgrafanacloud_secret_role` resource creates a Vault role on the Grafana Cloud secret backend.

#### Attributes

| Name | Required | Description | Default Value | 
| ---- | -------- | ----------- | ------------- |
| `backend` | `false` | The mount path of the Grafana Cloud backend | `grafana-cloud` |
| `name` | `true` | Grafana Cloud API key with Admin role to create user keys | N/A |
| `gc_role` | `true` | The URL for the Grafana Cloud API | N/A |
| `ttl_seconds` | `false` | The Organisation slug for the Grafana Cloud API" | `300` |
| `max_ttl_seconds` | `false` | The User that is needed to interact with prometheus, if set this is returned alongside every issued credential | `300` |

#### Example

```hcl
resource "vaultgrafanacloud_secret_backend" "backend" {
  backend      = "grafanacloud"
  key          = var.your_secret_api_key
  url          = "https://grafana.com/api"
  organisation = "my-org"
  user         = "my-user"
}

resource "vaultgrafanacloud_secret_role" "test" {
  backend         = "grafanacloud"
  name            = "my-role"
  gc_role         = "Viewer"
  ttl_seconds     = "3600"
  max_ttl_seconds = "3600"
}
```

## Testing

To test the terraform provider, you will need to perform some set-up steps.

1. Compile the [vault-plugin-secrets-grafanacloud](https://github.com/form3tech-oss/vault-plugin-secrets-grafanacloud) plugin and copy to `./bin/`.
2. Run `docker-compose up -d`
3. Run `TF_ACC=1 VAULT_ADDR=http://localhost:8200 VAULT_TOKEN=root make test`