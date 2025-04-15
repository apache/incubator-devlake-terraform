// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
	_ resource.Resource                = &bitbucketServerConnectionScopeResource{}
	_ resource.ResourceWithConfigure   = &bitbucketServerConnectionScopeResource{}
	_ resource.ResourceWithImportState = &bitbucketServerConnectionScopeResource{}
)

// NewBitbucketServerConnectionScopeResource is a helper function to simplify the provider implementation.
func NewBitbucketServerConnectionScopeResource() resource.Resource {
	return &bitbucketServerConnectionScopeResource{}
}

// bitbucketServerConnectionScopeResource is the resource implementation.
type bitbucketServerConnectionScopeResource struct {
	client *client.Client
}

// bitbucketServerConnectionScopeResourceModel maps the resource schema data.
type bitbucketServerConnectionScopeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	CloneUrl      types.String `tfsdk:"clone_url"`
	ConnectionId  types.String `tfsdk:"connection_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	Description   types.String `tfsdk:"description"`
	HTMLUrl       types.String `tfsdk:"html_url"`
	Name          types.String `tfsdk:"name"`
	ScopeConfigId types.String `tfsdk:"scope_config_id"`
}

// Metadata returns the resource type name.
func (r *bitbucketServerConnectionScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bitbucketserver_connection_scope"
}

// Schema defines the schema for the resource.
func (r *bitbucketServerConnectionScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The Bitbucket project and repository in the format '<PROJECT>/repos/<REPOSITORY>'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last Terraform update of the connection scope.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clone_url": schema.StringAttribute{
				Description: "The Bitbucket https clone url.",
				Required:    true,
			},
			"connection_id": schema.StringAttribute{
				Description: "The Connection this scope is part of.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the scope was created in devlake.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for the connection scope.",
				Required:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "The Bitbucket HTML browse url.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "A name for the connection scope.",
				Required:    true,
			},
			"scope_config_id": schema.StringAttribute{
				Description: "The config used for the scope. Needs to be created first.",
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *bitbucketServerConnectionScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scope",
			"Could not create devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	scopeConfigId, err := strconv.Atoi(plan.ScopeConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scope",
			"Could not create devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	now := time.Now().Format(time.RFC3339)
	var bitbucketServerConnectionScopeCreate = client.BitbucketServerConnectionScope{
		BitbucketId:   plan.ID.ValueString(),
		CloneUrl:      plan.CloneUrl.ValueString(),
		ConnectionId:  connectionId,
		CreatedAt:     now,
		Description:   plan.Description.ValueString(),
		HTMLUrl:       plan.HTMLUrl.ValueString(),
		Name:          plan.Name.ValueString(),
		ScopeConfigId: scopeConfigId,
		UpdatedAt:     now,
	}

	// Create new bitbucketserverconnectionscope
	bitbucketServerConnectionScope, err := r.client.CreateBitbucketServerConnectionScope(plan.ConnectionId.ValueString(), bitbucketServerConnectionScopeCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scope",
			"Could not create devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(bitbucketServerConnectionScope.BitbucketId)
	plan.LastUpdated = types.StringValue(now)
	plan.CloneUrl = types.StringValue(bitbucketServerConnectionScope.CloneUrl)
	plan.ConnectionId = types.StringValue(strconv.Itoa(bitbucketServerConnectionScope.ConnectionId))
	plan.CreatedAt = types.StringValue(bitbucketServerConnectionScope.CreatedAt)
	plan.Description = types.StringValue(bitbucketServerConnectionScope.Description)
	plan.HTMLUrl = types.StringValue(bitbucketServerConnectionScope.HTMLUrl)
	plan.Name = types.StringValue(bitbucketServerConnectionScope.Name)
	plan.ScopeConfigId = types.StringValue(strconv.Itoa(bitbucketServerConnectionScope.ScopeConfigId))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bitbucketServerConnectionScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state bitbucketServerConnectionScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed bitbucket server connection scope value from Devlake
	bitbucketServerConnectionScope, err := r.client.ReadBitbucketServerConnectionScope(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake bitbucket server connection scope",
			err.Error(),
		)
		return
	}

	// Overwrite connection with refreshed state
	state.ID = types.StringValue(bitbucketServerConnectionScope.BitbucketId)
	state.CloneUrl = types.StringValue(bitbucketServerConnectionScope.CloneUrl)
	state.ConnectionId = types.StringValue(strconv.Itoa(bitbucketServerConnectionScope.ConnectionId))
	state.CreatedAt = types.StringValue(bitbucketServerConnectionScope.CreatedAt)
	state.Description = types.StringValue(bitbucketServerConnectionScope.Description)
	state.HTMLUrl = types.StringValue(bitbucketServerConnectionScope.HTMLUrl)
	state.Name = types.StringValue(bitbucketServerConnectionScope.Name)
	state.ScopeConfigId = types.StringValue(strconv.Itoa(bitbucketServerConnectionScope.ScopeConfigId))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update fetches the resource and sets the updated Terraform state on success.
func (r *bitbucketServerConnectionScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection scope",
			"Could not update devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	scopeConfigId, err := strconv.Atoi(plan.ScopeConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scope",
			"Could not create devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	var bitbucketServerConnectionScopeUpdate = client.BitbucketServerConnectionScope{
		BitbucketId:   plan.ID.ValueString(),
		CloneUrl:      plan.CloneUrl.ValueString(),
		ConnectionId:  connectionId,
		CreatedAt:     plan.CreatedAt.ValueString(),
		Description:   plan.Description.ValueString(),
		HTMLUrl:       plan.HTMLUrl.ValueString(),
		Name:          plan.Name.ValueString(),
		ScopeConfigId: scopeConfigId,
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	// Update existing connection scope
	updatedBitbucketServerConnectionScope, err := r.client.UpdateBitbucketServerConnectionScope(plan.ConnectionId.ValueString(), plan.ID.ValueString(), bitbucketServerConnectionScopeUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection scope",
			"Could not update devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(updatedBitbucketServerConnectionScope.BitbucketId)
	plan.CloneUrl = types.StringValue(updatedBitbucketServerConnectionScope.CloneUrl)
	plan.ConnectionId = types.StringValue(strconv.Itoa(updatedBitbucketServerConnectionScope.ConnectionId))
	plan.CreatedAt = types.StringValue(updatedBitbucketServerConnectionScope.CreatedAt)
	plan.Description = types.StringValue(updatedBitbucketServerConnectionScope.Description)
	plan.HTMLUrl = types.StringValue(updatedBitbucketServerConnectionScope.HTMLUrl)
	plan.Name = types.StringValue(updatedBitbucketServerConnectionScope.Name)
	plan.ScopeConfigId = types.StringValue(strconv.Itoa(updatedBitbucketServerConnectionScope.ScopeConfigId))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bitbucketServerConnectionScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state bitbucketServerConnectionScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing apikey
	err := r.client.DeleteBitbucketServerConnectionScope(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake bitbucket server connection scope",
			"Could not delete devlake bitbucket server connection scope, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *bitbucketServerConnectionScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and connection id and save to attribute
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: connection_id,scope_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

// Configure adds the provider configured client to the resource.
func (r *bitbucketServerConnectionScopeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
