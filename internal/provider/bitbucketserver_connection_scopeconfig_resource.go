// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"terraform-provider-devlake/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bitbucketServerConnectionScopeConfigResource{}
	_ resource.ResourceWithConfigure   = &bitbucketServerConnectionScopeConfigResource{}
	_ resource.ResourceWithImportState = &bitbucketServerConnectionScopeConfigResource{}
)

// NewBitbucketServerConnectionScopeConfigResource is a helper function to simplify the provider implementation.
func NewBitbucketServerConnectionScopeConfigResource() resource.Resource {
	return &bitbucketServerConnectionScopeConfigResource{}
}

// bitbucketServerConnectionScopeConfigResource is the resource implementation.
type bitbucketServerConnectionScopeConfigResource struct {
	client *client.Client
}

// bitbucketServerConnectionScopeConfigResourceModel maps the resource schema data.
type bitbucketServerConnectionScopeConfigResourceModel struct {
	ID           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	ConnectionId types.String `tfsdk:"connection_id"`
	CreatedAt    types.String `tfsdk:"created_at"`
	Entities     types.List   `tfsdk:"entities"`
	Name         types.String `tfsdk:"name"`
	PrComponent  types.String `tfsdk:"pr_component"`
	PrType       types.String `tfsdk:"pr_type"`
	RefDiff      *refDiff     `tfsdk:"ref_diff"`
}

type refDiff struct {
	TagsLimit   types.Int64  `tfsdk:"tags_limit"`
	TagsPattern types.String `tfsdk:"tags_pattern"`
}

// Metadata returns the resource type name.
func (r *bitbucketServerConnectionScopeConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bitbucketserver_connection_scopeconfig"
}

// Schema defines the schema for the resource.
func (r *bitbucketServerConnectionScopeConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Numeric identifier for the connection scopeconfig. This is a string for easier resource import.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last Terraform update of the scope config.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_id": schema.StringAttribute{
				Description: "The connection id of the connection this scope config belongs to.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the scope config was created in devlake.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"entities": schema.ListAttribute{
				Computed:    true,
				Description: "The entities this scope config uses, e.g. 'CODEREVIEW', 'CROSS' or 'CODE'. See the documentation for the meaning of the individual values.",
				ElementType: types.StringType,
				Optional:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{
						types.StringValue("CODEREVIEW"),
						types.StringValue("CROSS"),
						types.StringValue("CODE"),
					},
				)),
			},
			"pr_component": schema.StringAttribute{
				Computed:    true,
				Description: "Text (PR body) that matches the RegEx will be set as the component of the pull request.",
				Optional:    true,
			},
			"pr_type": schema.StringAttribute{
				Computed:    true,
				Description: "Text (PR title) that matches the RegEx will be set as the type of a pull request.",
				Optional:    true,
			},
			"ref_diff": schema.SingleNestedAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Calculate the commits diff between two consecutive tags that match the following RegEx. Issues closed by PRs which contain these commits will also be calculated. The result will be shown in table.refs_commits_diffs and table.refs_issues_diffs.",
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					map[string]attr.Type{
						"tags_limit":   types.Int64Type,
						"tags_pattern": types.StringType,
					},
					map[string]attr.Value{
						"tags_limit":   types.Int64Value(10),
						"tags_pattern": types.StringValue(`/v\d+\.\d+(\.\d+(-rc)*\d*)*$/`),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"tags_limit": schema.Int64Attribute{
						Computed:    true,
						Default:     int64default.StaticInt64(10),
						Description: "Compare the last number of tags.",
						Optional:    true,
					},
					"tags_pattern": schema.StringAttribute{
						Computed:    true,
						Default:     stringdefault.StaticString(`/v\d+\.\d+(\.\d+(-rc)*\d*)*$/`),
						Description: "Matching tags are included in the calculation.",
						Optional:    true,
					},
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the scope config.",
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *bitbucketServerConnectionScopeConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionScopeConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var entities []string
	if !plan.Entities.IsNull() && !plan.Entities.IsUnknown() {
		diags = plan.Entities.ElementsAs(ctx, &entities, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scopeconfig",
			"Could not create devlake bitbucket server connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
	var bitbucketServerConnectionScopeConfigCreate = client.BitbucketServerConnectionScopeConfig{
		ConnectionId: connectionId,
		RefDiff: &client.RefDiff{
			TagsLimit:   int(plan.RefDiff.TagsLimit.ValueInt64()),
			TagsPattern: plan.RefDiff.TagsPattern.ValueString(),
		},
		Entities:    entities,
		Name:        plan.Name.ValueString(),
		PrComponent: plan.PrComponent.ValueString(),
		PrType:      plan.PrType.ValueString(),
	}

	// Create new bitbucketserverconnectionscopeconfig
	bitbucketServerConnectionScopeConfig, err := r.client.CreateBitbucketServerConnectionScopeConfig(plan.ConnectionId.ValueString(), bitbucketServerConnectionScopeConfigCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake bitbucket server connection scope config",
			"Could not create devlake bitbucket server connection scope config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, bitbucketServerConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Entities = entitiesVal
	plan.ID = types.StringValue(strconv.Itoa(bitbucketServerConnectionScopeConfig.ID))
	plan.CreatedAt = types.StringValue(bitbucketServerConnectionScopeConfig.CreatedAt)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.Name = types.StringValue(bitbucketServerConnectionScopeConfig.Name)
	plan.PrComponent = types.StringValue(bitbucketServerConnectionScopeConfig.PrComponent)
	plan.PrType = types.StringValue(bitbucketServerConnectionScopeConfig.PrType)
	plan.RefDiff.TagsLimit = types.Int64Value(int64(bitbucketServerConnectionScopeConfig.RefDiff.TagsLimit))
	plan.RefDiff.TagsPattern = types.StringValue(bitbucketServerConnectionScopeConfig.RefDiff.TagsPattern)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bitbucketServerConnectionScopeConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state bitbucketServerConnectionScopeConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed bitbucket server connection scope config value from Devlake
	bitbucketServerConnectionScopeConfig, err := r.client.ReadBitbucketServerConnectionScopeConfig(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake bitbucket server connection scopeconfig",
			err.Error(),
		)
		return
	}

	// Overwrite bitbucket server connection scope config with refreshed state
	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, bitbucketServerConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.CreatedAt = types.StringValue(bitbucketServerConnectionScopeConfig.CreatedAt)
	state.Entities = entitiesVal
	state.Name = types.StringValue(bitbucketServerConnectionScopeConfig.Name)
	state.PrComponent = types.StringValue(bitbucketServerConnectionScopeConfig.PrComponent)
	state.PrType = types.StringValue(bitbucketServerConnectionScopeConfig.PrType)
	if apiRefDiff := bitbucketServerConnectionScopeConfig.RefDiff; apiRefDiff != nil {
		state.RefDiff = &refDiff{
			TagsLimit:   types.Int64Value(int64(apiRefDiff.TagsLimit)),
			TagsPattern: types.StringValue(apiRefDiff.TagsPattern),
		}
	} else {
		state.RefDiff = nil
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update fetches the resource and sets the updated Terraform state on success.
func (r *bitbucketServerConnectionScopeConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan bitbucketServerConnectionScopeConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection scopeconfig",
			"Could not update devlake bitbucket server connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection scopeconfig",
			"Could not update devlake bitbucket server connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
	var entities []string
	if !plan.Entities.IsNull() && !plan.Entities.IsUnknown() {
		diags = plan.Entities.ElementsAs(ctx, &entities, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	var bitbucketServerConnectionScopeConfigUpdate = client.BitbucketServerConnectionScopeConfig{
		ConnectionId: connectionId,
		ID:           id,
		Entities:     entities,
		Name:         plan.Name.ValueString(),
		PrComponent:  plan.PrComponent.ValueString(),
		PrType:       plan.PrType.ValueString(),
		RefDiff: &client.RefDiff{
			TagsLimit:   int(plan.RefDiff.TagsLimit.ValueInt64()),
			TagsPattern: plan.RefDiff.TagsPattern.ValueString(),
		},
	}

	// Update existing connection scope config
	updatedBitbucketServerConnectionScopeConfig, err := r.client.UpdateBitbucketServerConnectionScopeConfig(plan.ConnectionId.ValueString(), plan.ID.ValueString(), bitbucketServerConnectionScopeConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake bitbucket server connection scopeconfig",
			"Could not update devlake bitbucket server connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}

	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, updatedBitbucketServerConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.CreatedAt = types.StringValue(updatedBitbucketServerConnectionScopeConfig.CreatedAt)
	plan.Entities = entitiesVal
	plan.Name = types.StringValue(updatedBitbucketServerConnectionScopeConfig.Name)
	plan.PrComponent = types.StringValue(updatedBitbucketServerConnectionScopeConfig.PrComponent)
	plan.PrType = types.StringValue(updatedBitbucketServerConnectionScopeConfig.PrType)
	plan.RefDiff.TagsLimit = types.Int64Value(int64(updatedBitbucketServerConnectionScopeConfig.RefDiff.TagsLimit))
	plan.RefDiff.TagsPattern = types.StringValue(updatedBitbucketServerConnectionScopeConfig.RefDiff.TagsPattern)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bitbucketServerConnectionScopeConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state bitbucketServerConnectionScopeConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing connection scope config
	err := r.client.DeleteBitbucketServerConnectionScopeConfig(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake bitbucket server connection scopeconfig",
			"Could not delete devlake bitbucket server connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *bitbucketServerConnectionScopeConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and connection id and save to attribute
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: connection_id,scopeconfig_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

// Configure adds the provider configured client to the resource.
func (r *bitbucketServerConnectionScopeConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
