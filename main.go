package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/wistia/terraform-provider-wistia/internal/provider"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New()}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.wistia.io/wistia/wistia", opts)
		if err != nil {
			log.Fatalf("error from Terraform debug plugin: %s", err)
		}
		return
	}

	plugin.Serve(opts)
}
