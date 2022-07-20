package vaultgrafanacloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/config"
)

const (
	EnvVaultAddr  = "VAULT_ADDR"
	EnvVaultToken = "VAULT_TOKEN"

	// DefaultMaxHTTPRetries is used for configuring the api.Client's MaxRetries.
	DefaultMaxHTTPRetries = 2
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVaultAddr, nil),
				Description: "URL of the root of the target Vault server.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(EnvVaultToken, ""),
				Description: "Token to use to authenticate to Vault.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vaultgrafanacloud_secret_backend": GrafanaCloudSecretBackendResource(),
			"vaultgrafanacloud_secret_role":    GrafanaCloudSecretRoleResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerToken(d *schema.ResourceData) (string, error) {
	if token := d.Get("token").(string); token != "" {
		return token, nil
	}

	// Use ~/.vault-token, or the configured token helper.
	tokenHelper, err := config.DefaultTokenHelper()
	if err != nil {
		return "", fmt.Errorf("error getting token helper: %s", err)
	}
	token, err := tokenHelper.Get()
	if err != nil {
		return "", fmt.Errorf("error getting token: %s", err)
	}
	return strings.TrimSpace(token), nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientConfig := api.DefaultConfig()
	addr := d.Get("address").(string)
	if addr != "" {
		clientConfig.Address = addr
	}

	// enable ReadYourWrites to support read-after-write on Vault Enterprise
	clientConfig.ReadYourWrites = true

	// set default MaxRetries
	clientConfig.MaxRetries = DefaultMaxHTTPRetries

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Vault API: %s", err)
	}

	// setting this is critical for proper namespace handling
	client.SetCloneHeaders(true)

	// setting this is critical for proper client cloning
	client.SetCloneToken(true)

	// Try an get the token from the config or token helper
	token, err := providerToken(d)
	if err != nil {
		return nil, err
	}

	if token != "" {
		client.SetToken(token)
	}
	if client.Token() == "" {
		return nil, errors.New("no vault token found")
	}

	return client, nil
}
