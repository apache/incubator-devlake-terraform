// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-devlake/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &devlakeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &devlakeProvider{
			version: version,
		}
	}
}

// devlakeProvider is the provider implementation.
type devlakeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// devlakeProviderModel maps provider schema data to a Go type.
type devlakeProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *devlakeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devlake"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *devlakeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "URI for Devlake API. May also be provided via DEVLAKE_HOST environment variable.",
			},
			"token": schema.StringAttribute{
				Description: "Token for Devlake API. May also be provided via DEVLAKE_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a devlake API client for data sources and resources.
func (p *devlakeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Devlake client")

	// Retrieve provider data from configuration
	var config devlakeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown devlake api host",
			"The provider cannot create the devlake api client as there is an unknown configuration value for the devlake api host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DEVLAKE_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown devlake api token",
			"The provider cannot create the devlake api client as there is an unknown configuration value for the devlake api token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DEVLAKE_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("DEVLAKE_HOST")
	token := os.Getenv("DEVLAKE_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing devlake api host",
			"The provider cannot create the devlake api client as there is a missing or empty value for the devlake api host. "+
				"Set the host value in the configuration or use the DEVLAKE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing devlake api token",
			"The provider cannot create the devlake api client as there is a missing or empty value for the devlake api token. "+
				"Set the token value in the configuration or use the DEVLAKE_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "devlake_host", host)
	ctx = tflog.SetField(ctx, "devlake_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "devlake_token")

	tflog.Debug(ctx, "Creating Devlake client")

	// Create a new Devlake client using the configuration values
	client, err := client.NewClient(&host, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create devlake api client",
			"An unexpected error occurred when creating the devlake api client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"devlake client error: "+err.Error(),
		)
		return
	}

	// Make the Devlake client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Devlake client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *devlakeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApiKeysDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *devlakeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiKeyResource,
		NewBitbucketServerConnectionResource,
		NewBitbucketServerConnectionScopeConfigResource,
		NewBitbucketServerConnectionScopeResource,
		NewGithubConnectionResource,
		NewGithubConnectionScopeConfigResource,
	}
}
