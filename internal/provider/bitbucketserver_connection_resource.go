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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bitbucketServerConnectionResource{}
	_ resource.ResourceWithConfigure   = &bitbucketServerConnectionResource{}
	_ resource.ResourceWithImportState = &bitbucketServerConnectionResource{}
)

// NewBitbucketServerConnectionResource is a helper function to simplify the provider implementation.
func NewBitbucketServerConnectionResource() resource.Resource {
	return &bitbucketServerConnectionResource{}
}

// bitbucketServerConnectionResource is the resource implementation.
type bitbucketServerConnectionResource struct {
	client *client.Client
}

// bitbucketServerConnectionResourceModel maps the resource schema data.
type bitbucketServerConnectionResourceModel struct {
	ID               types.String `tfsdk:"id"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	CreatedAt        types.String `tfsdk:"created_at"`
	Endpoint         types.String `tfsdk:"endpoint"`
	Name             types.String `tfsdk:"name"`
	Password         types.String `tfsdk:"password"`
	Proxy            types.String `tfsdk:"proxy"`
	RateLimitPerHour types.Int64  `tfsdk:"rate_limit_per_hour"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	Username         types.String `tfsdk:"username"`
}

// Metadata returns the resource type name.
func (r *bitbucketServerConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bitbucketserver_connection"
}

// Schema defines the schema for the resource.
func (r *bitbucketServerConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the connection was created in devlake.",
			},
			"endpoint": schema.StringAttribute{
				Description: "The base endpoint URL.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the bitbucket server connection.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Service account password or token, the following permissions are required to collect data from Bitbucket repositories: Repository read.",
				Required:    true,
				Sensitive:   true,
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
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the connection was updated in devlake.",
			},
			"username": schema.StringAttribute{
				Description: "Service account username.",
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *bitbucketServerConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	now := time.Now().Format(time.RFC850)

	// Generate API request body from plan
	var bitbucketServerConnectionCreate = client.BitbucketServerConnection{
		CreatedAt:        now,
		Endpoint:         plan.Endpoint.ValueString(),
		Name:             plan.Name.ValueString(),
		Password:         plan.Password.ValueString(),
		Proxy:            plan.Proxy.ValueString(),
		RateLimitPerHour: int(plan.RateLimitPerHour.ValueInt64()),
		UpdatedAt:        now,
		Username:         plan.Username.ValueString(),
	}

	// Create new bitbucketserverconnection
	bitbucketServerConnection, err := r.client.CreateBitbucketServerConnection(bitbucketServerConnectionCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection",
			"Could not create devlake bitbucket server connection, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(bitbucketServerConnection.ID))
	plan.LastUpdated = types.StringValue(now)
	plan.CreatedAt = types.StringValue(bitbucketServerConnection.CreatedAt)
	plan.Endpoint = types.StringValue(bitbucketServerConnection.Endpoint)
	plan.Name = types.StringValue(bitbucketServerConnection.Name)
	plan.Proxy = types.StringValue(bitbucketServerConnection.Proxy)
	plan.RateLimitPerHour = types.Int64Value(int64(bitbucketServerConnection.RateLimitPerHour))
	plan.UpdatedAt = types.StringValue(bitbucketServerConnection.UpdatedAt)
	plan.Username = types.StringValue(bitbucketServerConnection.Username)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bitbucketServerConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state bitbucketServerConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed bitbucket server connection value from Devlake
	bitbucketServerConnection, err := r.client.ReadBitbucketServerConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake bitbucket server connection",
			err.Error(),
		)
		return
	}

	// Overwrite connection with refreshed state
	state.ID = types.StringValue(strconv.Itoa(bitbucketServerConnection.ID))
	state.CreatedAt = types.StringValue(bitbucketServerConnection.CreatedAt)
	state.Endpoint = types.StringValue(bitbucketServerConnection.Endpoint)
	state.Name = types.StringValue(bitbucketServerConnection.Name)
	state.Proxy = types.StringValue(bitbucketServerConnection.Proxy)
	state.RateLimitPerHour = types.Int64Value(int64(bitbucketServerConnection.RateLimitPerHour))
	state.UpdatedAt = types.StringValue(bitbucketServerConnection.UpdatedAt)
	state.Username = types.StringValue(bitbucketServerConnection.Username)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update fetches the resource and sets the updated Terraform state on success.
func (r *bitbucketServerConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection",
			"Could not update devlake bitbucket server connection, unexpected error: "+err.Error(),
		)
		return
	}
	var bitbucketServerConnectionUpdate = client.BitbucketServerConnection{
		ID:               id,
		CreatedAt:        plan.CreatedAt.ValueString(),
		Endpoint:         plan.Endpoint.ValueString(),
		Name:             plan.Name.ValueString(),
		Password:         plan.Password.ValueString(),
		Proxy:            plan.Proxy.ValueString(),
		RateLimitPerHour: int(plan.RateLimitPerHour.ValueInt64()),
		UpdatedAt:        time.Now().Format(time.RFC850),
		Username:         plan.Username.ValueString(),
	}

	// Update existing connection
	updatedBitbucketServerConnection, err := r.client.UpdateBitbucketServerConnection(plan.ID.ValueString(), bitbucketServerConnectionUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection",
			"Could not update devlake bitbucket server connection, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(updatedBitbucketServerConnection.ID))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.CreatedAt = types.StringValue(updatedBitbucketServerConnection.CreatedAt)
	plan.Endpoint = types.StringValue(updatedBitbucketServerConnection.Endpoint)
	plan.Name = types.StringValue(updatedBitbucketServerConnection.Name)
	plan.Proxy = types.StringValue(updatedBitbucketServerConnection.Proxy)
	plan.RateLimitPerHour = types.Int64Value(int64(updatedBitbucketServerConnection.RateLimitPerHour))
	plan.UpdatedAt = types.StringValue(updatedBitbucketServerConnection.UpdatedAt)
	plan.Username = types.StringValue(updatedBitbucketServerConnection.Username)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bitbucketServerConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state bitbucketServerConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing apikey
	err := r.client.DeleteBitbucketServerConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake bitbucket server connection",
			"Could not delete devlake bitbucket server connection, unexpected error: "+err.Error()+"..",
		)
		return
	}
}

func (r *bitbucketServerConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *bitbucketServerConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
