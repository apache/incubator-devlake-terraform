// Copyright (c) HashiCorp, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the devlake client is properly configured.
	// It is also possible to use the DEVLAKE_ environment variables instead,
	// update the environment variables the Makefile if you want to use that.
	providerConfig = `
provider "devlake" {
  host  = "http://localhost:4000/api"
  token = "whatever"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"devlake": providerserver.NewProtocol6WithError(New("test")()),
	}
)
