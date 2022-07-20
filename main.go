package main

import (
	"flag"

	"github.com/form3tech-oss/terraform-provider-vault-grafanacloud/vaultgrafanacloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vaultgrafanacloud.Provider,
		Debug:        debugMode,
	})
}
