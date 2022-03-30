# Terraform Provider Vault Grafana Cloud

This terraform provider is intended to be used to provision a [Grafana Cloud Secrets](https://github.com/form3tech-oss/vault-plugin-secrets-grafanacloud) backend and roles in a Vault instance.

## Testing

To test the terraform provider, you will need to perform some set-up steps.

1. Compile the [vault-plugin-secrets-grafanacloud](https://github.com/form3tech-oss/vault-plugin-secrets-grafanacloud) plugin and copy to `./bin/`.
2. Run `docker-compose up -d`
3. Run `make test`
