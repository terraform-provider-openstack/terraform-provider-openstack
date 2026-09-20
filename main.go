package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack"
)

const providerAddr = "registry.terraform.io/terraform-provider-openstack/openstack"

func main() {
	// added debugMode to enable debugging for provider per https://developer.hashicorp.com/terraform/plugin/sdkv2/debugging
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        debugMode,
		ProviderAddr: providerAddr,
		ProviderFunc: openstack.Provider,
	})
}
