// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-devlake/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &apiKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &apiKeysDataSource{}
)

// NewApiKeysDataSource is a helper function to simplify the provider implementation.
func NewApiKeysDataSource() datasource.DataSource {
	return &apiKeysDataSource{}
}

// apiKeysDataSource is the data source implementation.
type apiKeysDataSource struct {
	client *client.Client
}

// apiKeysDataSourceModel maps the data source schema data.
type apiKeysDataSourceModel struct {
	ApiKeys []apiKeysModel `tfsdk:"apikeys"`
}

// apiKeysModel maps apiKeys schema data.
type apiKeysModel struct {
	ID           types.Int64  `tfsdk:"id"`
	AllowedPath  types.String `tfsdk:"allowed_path"`
	ApiKey       types.String `tfsdk:"api_key"`
	CreatedAt    types.String `tfsdk:"created_at"`
	Creator      types.String `tfsdk:"creator"`
	CreatorEmail types.String `tfsdk:"creator_email"`
	ExpiredAt    types.String `tfsdk:"expired_at"`
	Extra        types.String `tfsdk:"extra"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	Updater      types.String `tfsdk:"updater"`
	UpdaterEmail types.String `tfsdk:"updater_email"`
}

// Metadata returns the data source type name.
func (d *apiKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikeys"
}

// Schema defines the schema for the data source.
func (d *apiKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apikeys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Numeric identifier for the apikey.",
						},
						"allowed_path": schema.StringAttribute{
							Computed:    true,
							Description: "The API URL or endpoint that the API key is permitted to access. It defines the specific resources that the key can interact with.",
						},
						"api_key": schema.StringAttribute{
							Computed:    true,
							Description: "The returned apikey, will always be empty for security reasons. If you want to use the apikey, you need to create a new resource.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "When the apikey was created.",
						},
						"creator": schema.StringAttribute{
							Computed:    true,
							Description: "Who created the apikey, there is no user management yet though.",
						},
						"creator_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email of the person who created the apikey, there is no user management yet though.",
						},
						"expired_at": schema.StringAttribute{
							Computed:    true,
							Description: "When the apikey expires.",
						},
						"extra": schema.StringAttribute{
							Computed:    true,
							Description: "Currently not used.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the apikey.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The apikey type. Currently only 'devlake' is a valid value.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "When the apikey was last updated. Serves no purpose as the UPDATE endpoint doesn't do anything.",
						},
						"updater": schema.StringAttribute{
							Computed:    true,
							Description: "Who updated the apikey, there is no user management yet though.",
						},
						"updater_email": schema.StringAttribute{
							Computed:    true,
							Description: "Email of the person who updated the apikey, there is no user management yet though.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *apiKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state apiKeysDataSourceModel

	apiKeys, err := d.client.GetApiKeys()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Devlake ApiKeys",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, apiKey := range apiKeys {
		apiKeyState := apiKeysModel{
			ID:           types.Int64Value(int64(apiKey.ID)),
			AllowedPath:  types.StringValue(apiKey.AllowedPath),
			ApiKey:       types.StringValue(apiKey.ApiKey),
			CreatedAt:    types.StringValue(apiKey.CreatedAt),
			Creator:      types.StringValue(apiKey.Creator),
			CreatorEmail: types.StringValue(apiKey.CreatorEmail),
			ExpiredAt:    types.StringValue(apiKey.ExpiredAt),
			Extra:        types.StringValue(apiKey.Extra),
			Name:         types.StringValue(apiKey.Name),
			Type:         types.StringValue(apiKey.Type),
			UpdatedAt:    types.StringValue(apiKey.UpdatedAt),
			Updater:      types.StringValue(apiKey.Updater),
			UpdaterEmail: types.StringValue(apiKey.UpdaterEmail),
		}

		state.ApiKeys = append(state.ApiKeys, apiKeyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *apiKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
