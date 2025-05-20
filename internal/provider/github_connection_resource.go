// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"terraform-provider-devlake/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &githubConnectionResource{}
	_ resource.ResourceWithConfigure   = &githubConnectionResource{}
	_ resource.ResourceWithImportState = &githubConnectionResource{}
)

// NewGithubConnectionResource is a helper function to simplify the provider implementation.
func NewGithubConnectionResource() resource.Resource {
	return &githubConnectionResource{}
}

// githubConnectionResource is the resource implementation.
type githubConnectionResource struct {
	client *client.Client
}

// githubConnectionResourceModel maps the resource schema data.
type githubConnectionResourceModel struct {
	ID               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	AppId            types.Int64  `tfsdk:"app_id"`
	AuthMethod       types.String `tfsdk:"auth_method"`
	CreatedAt        types.String `tfsdk:"created_at"`
	EnableGraphql    types.Bool   `tfsdk:"enable_graphql"`
	Endpoint         types.String `tfsdk:"endpoint"`
	InstallationId   types.Int64  `tfsdk:"installation_id"`
	Name             types.String `tfsdk:"name"`
	Proxy            types.String `tfsdk:"proxy"`
	RateLimitPerHour types.Int64  `tfsdk:"rate_limit_per_hour"`
	SecretKey        types.String `tfsdk:"secret_key"`
	Token            types.String `tfsdk:"token"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *githubConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_connection"
}

// Schema defines the schema for the resource.
func (r *githubConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Numeric identifier for the connection. This is a string for easier resource import.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last Terraform update of the connection.",
			},
			"app_id": schema.Int64Attribute{
				Description: "The app id of the github app used for authentication.",
				Required:    true,
			},
			"auth_method": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("AppKey"),
				Description: "The authentication method. Currently only supports 'AppKey'.",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the connection was created in devlake.",
			},
			"enable_graphql": schema.BoolAttribute{
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to use the faster graphql api endpoints. Defaults to 'true'.",
				Optional:    true,
			},
			"endpoint": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("https://api.github.com/"),
				Description: "The base endpoint URL. Defaults to 'https://api.github.com/'.",
				Optional:    true,
			},
			"installation_id": schema.Int64Attribute{
				Description: "The installation id of the github app used for authentication.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the github connection.",
				Required:    true,
			},
			"proxy": schema.StringAttribute{
				Computed:    true,
				Description: "If you are behind a corporate firewall or VPN you may need to utilize a proxy server.",
				Optional:    true,
			},
			"rate_limit_per_hour": schema.Int64Attribute{
				Optional:    true,
				Description: "DevLake uses a dynamic rate limit to collect Bitbucket Server/Data Center data. You can adjust the rate limit if you want to increase or lower the speed.",
				Computed:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "Github app private key used for authentication.",
				Required:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "PAT for github authentication. Currently not supported.",
				Optional:    true,
				Sensitive:   true,
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the connection was updated in devlake.",
			},
		},
	}
}

// Create a new resource.
func (r *githubConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan githubConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	now := time.Now().Format(time.RFC850)

	// Generate API request body from plan
	var githubConnectionCreate = client.GithubConnection{
		AppId:            strconv.Itoa(int(plan.AppId.ValueInt64())),
		AuthMethod:       plan.AuthMethod.ValueString(),
		CreatedAt:        now,
		EnableGraphql:    plan.EnableGraphql.ValueBool(),
		Endpoint:         plan.Endpoint.ValueString(),
		InstallationId:   int(plan.InstallationId.ValueInt64()),
		Name:             plan.Name.ValueString(),
		Proxy:            plan.Proxy.ValueString(),
		RateLimitPerHour: int(plan.RateLimitPerHour.ValueInt64()),
		SecretKey:        plan.SecretKey.ValueString(),
		// not passing in token
		UpdatedAt: now,
	}

	// Create new githubconnection
	githubConnection, err := r.client.CreateGithubConnection(githubConnectionCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection",
			"Could not create devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	appId, err := strconv.Atoi(githubConnection.AppId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection",
			"Could not create devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(strconv.Itoa(githubConnection.ID))
	plan.LastUpdated = types.StringValue(now)
	plan.AppId = types.Int64Value(int64(appId))
	plan.AuthMethod = types.StringValue(githubConnection.AuthMethod)
	plan.CreatedAt = types.StringValue(githubConnection.CreatedAt)
	plan.EnableGraphql = types.BoolValue(githubConnection.EnableGraphql)
	plan.Endpoint = types.StringValue(githubConnection.Endpoint)
	plan.InstallationId = types.Int64Value(int64(githubConnection.InstallationId))
	plan.Name = types.StringValue(githubConnection.Name)
	plan.Proxy = types.StringValue(githubConnection.Proxy)
	plan.RateLimitPerHour = types.Int64Value(int64(githubConnection.RateLimitPerHour))
	plan.Token = types.StringValue(githubConnection.Token)
	plan.UpdatedAt = types.StringValue(githubConnection.UpdatedAt)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *githubConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state githubConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed github connection value from Devlake
	githubConnection, err := r.client.ReadGithubConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake github connection",
			err.Error(),
		)
		return
	}

	// Overwrite connection with refreshed state
	appId, err := strconv.Atoi(githubConnection.AppId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading devlake github connection",
			"Could not read devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}
	state.ID = types.StringValue(strconv.Itoa(githubConnection.ID))
	state.AppId = types.Int64Value(int64(appId))
	state.AuthMethod = types.StringValue(githubConnection.AuthMethod)
	state.CreatedAt = types.StringValue(githubConnection.CreatedAt)
	state.EnableGraphql = types.BoolValue(githubConnection.EnableGraphql)
	state.Endpoint = types.StringValue(githubConnection.Endpoint)
	state.InstallationId = types.Int64Value(int64(githubConnection.InstallationId))
	state.Name = types.StringValue(githubConnection.Name)
	state.Proxy = types.StringValue(githubConnection.Proxy)
	state.RateLimitPerHour = types.Int64Value(int64(githubConnection.RateLimitPerHour))
	state.Token = types.StringValue(githubConnection.Token)
	state.UpdatedAt = types.StringValue(githubConnection.UpdatedAt)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update fetches the resource and sets the updated Terraform state on success.
func (r *githubConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan githubConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection",
			"Could not update devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}
	var githubConnectionUpdate = client.GithubConnection{
		ID:               id,
		AppId:            strconv.Itoa(int(plan.AppId.ValueInt64())),
		AuthMethod:       plan.AuthMethod.ValueString(),
		CreatedAt:        plan.CreatedAt.ValueString(),
		EnableGraphql:    plan.EnableGraphql.ValueBool(),
		Endpoint:         plan.Endpoint.ValueString(),
		InstallationId:   int(plan.InstallationId.ValueInt64()),
		Name:             plan.Name.ValueString(),
		Proxy:            plan.Proxy.ValueString(),
		RateLimitPerHour: int(plan.RateLimitPerHour.ValueInt64()),
		SecretKey:        plan.SecretKey.ValueString(),
		UpdatedAt:        time.Now().Format(time.RFC850),
	}

	// Update existing connection
	updatedGithubConnection, err := r.client.UpdateGithubConnection(plan.ID.ValueString(), githubConnectionUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection",
			"Could not update devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}

	appId, err := strconv.Atoi(updatedGithubConnection.AppId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading devlake github connection",
			"Could not read devlake github connection, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(strconv.Itoa(updatedGithubConnection.ID))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.AppId = types.Int64Value(int64(appId))
	plan.AuthMethod = types.StringValue(updatedGithubConnection.AuthMethod)
	plan.CreatedAt = types.StringValue(updatedGithubConnection.CreatedAt)
	plan.EnableGraphql = types.BoolValue(updatedGithubConnection.EnableGraphql)
	plan.Endpoint = types.StringValue(updatedGithubConnection.Endpoint)
	plan.Name = types.StringValue(updatedGithubConnection.Name)
	plan.Proxy = types.StringValue(updatedGithubConnection.Proxy)
	plan.RateLimitPerHour = types.Int64Value(int64(updatedGithubConnection.RateLimitPerHour))
	plan.Token = types.StringValue(updatedGithubConnection.Token)
	plan.UpdatedAt = types.StringValue(updatedGithubConnection.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *githubConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state githubConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing apikey
	err := r.client.DeleteGithubConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake github connection",
			"Could not delete devlake github connection, unexpected error: "+err.Error()+"..",
		)
		return
	}
}

func (r *githubConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *githubConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
