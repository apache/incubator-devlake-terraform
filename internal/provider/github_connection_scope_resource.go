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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &githubConnectionScopeResource{}
	_ resource.ResourceWithConfigure   = &githubConnectionScopeResource{}
	_ resource.ResourceWithImportState = &githubConnectionScopeResource{}
)

// NewGithubConnectionScopeResource is a helper function to simplify the provider implementation.
func NewGithubConnectionScopeResource() resource.Resource {
	return &githubConnectionScopeResource{}
}

// githubConnectionScopeResource is the resource implementation.
type githubConnectionScopeResource struct {
	client *client.Client
}

// githubConnectionScopeResourceModel maps the resource schema data.
type githubConnectionScopeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	ConnectionId  types.String `tfsdk:"connection_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	Description   types.String `tfsdk:"description"`
	FullName      types.String `tfsdk:"full_name"`
	ScopeConfigId types.String `tfsdk:"scope_config_id"`
}

// Metadata returns the resource type name.
func (r *githubConnectionScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_connection_scope"
}

// Schema defines the schema for the resource.
func (r *githubConnectionScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the repository in github.",
				Required:    true,
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
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "A description for the connection scope.",
				Optional:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "The Github org and repository in the format '<ORG>/<REPOSITORY>'.",
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
func (r *githubConnectionScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan githubConnectionScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope",
			"Could not create devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope",
			"Could not create devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	scopeConfigId, err := strconv.Atoi(plan.ScopeConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope",
			"Could not create devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	now := time.Now().Format(time.RFC3339)
	var githubConnectionScopeCreate = client.GithubConnectionScope{
		GithubId:      id,
		CloneUrl:      "https://github.com/" + plan.FullName.ValueString() + ".git",
		ConnectionId:  connectionId,
		CreatedAt:     now,
		Description:   plan.Description.ValueString(),
		FullName:      plan.FullName.ValueString(),
		HTMLUrl:       "https://github.com/" + plan.FullName.ValueString(),
		Name:          strings.Split(plan.FullName.ValueString(), "/")[1],
		ScopeConfigId: scopeConfigId,
		UpdatedAt:     now,
		CreatedDate:   now,
		UpdatedDate:   now,
	}

	// Create new githubconnectionscope
	githubConnectionScope, err := r.client.CreateGithubConnectionScope(plan.ConnectionId.ValueString(), githubConnectionScopeCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope",
			"Could not create devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(githubConnectionScope.GithubId))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	plan.ConnectionId = types.StringValue(strconv.Itoa(githubConnectionScope.ConnectionId))
	plan.CreatedAt = types.StringValue(githubConnectionScope.CreatedAt)
	plan.Description = types.StringValue(githubConnectionScope.Description)
	plan.FullName = types.StringValue(githubConnectionScope.FullName)
	plan.ScopeConfigId = types.StringValue(strconv.Itoa(githubConnectionScope.ScopeConfigId))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *githubConnectionScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state githubConnectionScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed github connection scope value from Devlake
	githubConnectionScope, err := r.client.ReadGithubConnectionScope(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake github connection scope",
			err.Error(),
		)
		return
	}

	// Overwrite connection with refreshed state
	state.ID = types.StringValue(strconv.Itoa(githubConnectionScope.GithubId))
	state.ConnectionId = types.StringValue(strconv.Itoa(githubConnectionScope.ConnectionId))
	state.CreatedAt = types.StringValue(githubConnectionScope.CreatedAt)
	state.Description = types.StringValue(githubConnectionScope.Description)
	state.FullName = types.StringValue(githubConnectionScope.FullName)
	state.ScopeConfigId = types.StringValue(strconv.Itoa(githubConnectionScope.ScopeConfigId))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update fetches the resource and sets the updated Terraform state on success.
func (r *githubConnectionScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan githubConnectionScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scope",
			"Could not update devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scope",
			"Could not update devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	scopeConfigId, err := strconv.Atoi(plan.ScopeConfigId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope",
			"Could not create devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	var githubConnectionScopeUpdate = client.GithubConnectionScope{
		GithubId:      id,
		CloneUrl:      "https://github.com/" + plan.FullName.ValueString() + ".git",
		ConnectionId:  connectionId,
		CreatedAt:     plan.CreatedAt.ValueString(),
		Description:   plan.Description.ValueString(),
		FullName:      plan.FullName.ValueString(),
		HTMLUrl:       "https://github.com/" + plan.FullName.ValueString(),
		Name:          strings.Split(plan.FullName.ValueString(), "/")[1],
		ScopeConfigId: scopeConfigId,
		UpdatedAt:     time.Now().Format(time.RFC3339),
		CreatedDate:   plan.CreatedAt.ValueString(),
		UpdatedDate:   time.Now().Format(time.RFC3339),
	}

	// Update existing connection scope
	updatedGithubConnectionScope, err := r.client.UpdateGithubConnectionScope(plan.ConnectionId.ValueString(), plan.ID.ValueString(), githubConnectionScopeUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scope",
			"Could not update devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(strconv.Itoa(updatedGithubConnectionScope.GithubId))
	plan.ConnectionId = types.StringValue(strconv.Itoa(updatedGithubConnectionScope.ConnectionId))
	plan.CreatedAt = types.StringValue(updatedGithubConnectionScope.CreatedAt)
	plan.Description = types.StringValue(updatedGithubConnectionScope.Description)
	plan.FullName = types.StringValue(updatedGithubConnectionScope.FullName)
	plan.ScopeConfigId = types.StringValue(strconv.Itoa(updatedGithubConnectionScope.ScopeConfigId))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *githubConnectionScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state githubConnectionScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing apikey
	err := r.client.DeleteGithubConnectionScope(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake github connection scope",
			"Could not delete devlake github connection scope, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *githubConnectionScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *githubConnectionScopeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
