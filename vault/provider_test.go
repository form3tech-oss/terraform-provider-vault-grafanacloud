package vault

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// How to run the acceptance tests for this provider:
//
// - Obtain an official Vault release from the Vault website at
//   https://vaultproject.io/ and extract the "vault" binary
//   somewhere.
//
// - Run the following to start the Vault server in development mode:
//       vault server -dev
//
// - Take the "Root Token" value printed by Vault as the server started
//   up and set it as the value of the VAULT_TOKEN environment variable
//   in a new shell whose current working directory is the root of the
//   Terraform repository.
//
// - As directed by the Vault server output, set the VAULT_ADDR environment
//   variable. e.g.:
//       export VAULT_ADDR='http://127.0.0.1:8200'
//
// - Run the Terraform acceptance tests as usual:
//       make testacc TEST=./builtin/providers/vault
//
// The tests expect to be run in a fresh, empty Vault and thus do not attempt
// to randomize or otherwise make the generated resource paths unique on
// each run. In case of weird behavior, restart the Vault dev server to
// start over with a fresh Vault. (Remember to reset VAULT_TOKEN.)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var (
	testProvider  *schema.Provider
	testProviders map[string]*schema.Provider
)

func init() {
	testProvider = Provider()
	testProviders = map[string]*schema.Provider{
		"vault": testProvider,
	}
}
