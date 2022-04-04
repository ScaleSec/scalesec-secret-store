package main

// Import packages for the plugin

import (
	"fmt"
	"os"

	hclogger "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	scaleSecPlugin "scalesec.com/scalesec-secret-store"
)

var (
	pluginBuildDate    string
	pluginBuildVersion string
	pluginBuildInfo    string
)

func main() {

	// DO NOT sent output to stdout in the main function.  Vault will see it and think there is an
	// error loading the plugin.  Instead .. send data to the vault logger

	// Display log information about the plugin in the Vault log file.
	// The Var Values are set by the building of the plugin in the link options.
	// This allows you see and verify that vault is picking up the plugin version you expect

	hclogger.Default().Info("*********************************")
	hclogger.Default().Info("ScaleSec Secret Store Vault Plugin")
	hclogger.Default().Info(fmt.Sprintf("Plugin Built Date: %s ", pluginBuildDate))
	hclogger.Default().Info(fmt.Sprintf("Plugin Built Version: %s ", pluginBuildVersion))
	hclogger.Default().Info(fmt.Sprintf("Plugin Built Information: %s ", pluginBuildInfo))
	hclogger.Default().Info("*********************************")

	// The main function is a very standard function and follow HashiCorp example recomendation for
	// a custom plugin.

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: scaleSecPlugin.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclogger.New(&hclogger.LoggerOptions{})

		logger.Error("scalesecSecret Store plugin shutting down", "error", err)
		os.Exit(1)
	}
}
