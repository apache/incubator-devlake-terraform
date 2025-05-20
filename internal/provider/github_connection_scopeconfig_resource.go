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
	_ resource.Resource                = &githubConnectionScopeConfigResource{}
	_ resource.ResourceWithConfigure   = &githubConnectionScopeConfigResource{}
	_ resource.ResourceWithImportState = &githubConnectionScopeConfigResource{}
)

// NewGithubConnectionScopeConfigResource is a helper function to simplify the provider implementation.
func NewGithubConnectionScopeConfigResource() resource.Resource {
	return &githubConnectionScopeConfigResource{}
}

// githubConnectionScopeConfigResource is the resource implementation.
type githubConnectionScopeConfigResource struct {
	client *client.Client
}

// githubConnectionScopeConfigResourceModel maps the resource schema data.
type githubConnectionScopeConfigResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	ConnectionId         types.String `tfsdk:"connection_id"`
	CreatedAt            types.String `tfsdk:"created_at"`
	DeploymentPattern    types.String `tfsdk:"deployment_pattern"`
	Entities             types.List   `tfsdk:"entities"`
	EnvNamePattern       types.String `tfsdk:"env_name_pattern"`
	IssueComponent       types.String `tfsdk:"issue_component"`
	IssuePriority        types.String `tfsdk:"issue_priority"`
	IssueSeverity        types.String `tfsdk:"issue_severity"`
	IssueTypeBug         types.String `tfsdk:"issue_type_bug"`
	IssueTypeIncident    types.String `tfsdk:"issue_type_incident"`
	IssueTypeRequirement types.String `tfsdk:"issue_type_requirement"`
	Name                 types.String `tfsdk:"name"`
	PrBodyClosePattern   types.String `tfsdk:"pr_body_close_pattern"`
	PrComponent          types.String `tfsdk:"pr_component"`
	PrType               types.String `tfsdk:"pr_type"`
	ProductionPattern    types.String `tfsdk:"production_pattern"`
	RefDiff              *refDiff     `tfsdk:"ref_diff"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *githubConnectionScopeConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_github_connection_scopeconfig"
}

// Schema defines the schema for the resource.
func (r *githubConnectionScopeConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"deployment_pattern": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Convert a GitHub workflow run as a DevLake Deployment when: The name of the GitHub workflow run or one of its jobs matches this pattern.",
				Optional:    true,
			},
			"entities": schema.ListAttribute{
				Computed:    true,
				Description: "The entities this scope config uses, e.g. 'CODEREVIEW', 'CROSS' or 'CODE'. See the documentation for the meaning of the individual values.",
				ElementType: types.StringType,
				Optional:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType,
					[]attr.Value{
						types.StringValue("CODE"),
						types.StringValue("CODEREVIEW"),
						types.StringValue("CROSS"),
						types.StringValue("CICD"),
					},
				)),
			},
			"env_name_pattern": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "If its environment name matches this pattern, this deployment is a 'Production Deployment'.",
				Optional:    true,
			},
			"issue_component": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("component(.*)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"issue_priority": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("(highest|high|medium|low|p0|p1|p2|p3)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"issue_severity": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("severity(.*)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"issue_type_bug": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("(bug|broken)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"issue_type_incident": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("(incident|failure)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"issue_type_requirement": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("(feat|feature|proposal|requirement)"),
				Description: "This looks like an error in the API, the webinterface doesn't provide a field for this.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the scope config.",
				Required:    true,
			},
			"pr_body_close_pattern": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString(`(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\s]*.*(((and )?(#|https:\/\/github.com\/%s\/issues\/)\d+[ ]*)+)`),
				Description: "Connect entities across domains to measure metrics such as Bug Count per 1k Lines of Code. Connect PRs and Issues with the following pattern.",
				Optional:    true,
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
			"production_pattern": schema.StringAttribute{
				Computed:    true,
				Description: "Convert a GitHub workflow run as a DevLake Deployment when: If the name or its branch’s name also matches this pattern, this deployment is a 'Production Deployment'. Use only with 'deployment_pattern'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the connection was updated in devlake.",
			},
		},
	}
}

// Create a new resource.
func (r *githubConnectionScopeConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan githubConnectionScopeConfigResourceModel
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
			"Error creating devlake github connection scopeconfig",
			"Could not create devlake github connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
	now := time.Now().Format(time.RFC850)
	var githubConnectionScopeConfigCreate = client.GithubConnectionScopeConfig{
		ConnectionId:         connectionId,
		CreatedAt:            now,
		DeploymentPattern:    plan.DeploymentPattern.ValueString(),
		Entities:             entities,
		EnvNamePattern:       plan.EnvNamePattern.ValueString(),
		IssueComponent:       plan.IssueComponent.ValueString(),
		IssuePriority:        plan.IssuePriority.ValueString(),
		IssueSeverity:        plan.IssueSeverity.ValueString(),
		IssueTypeBug:         plan.IssueTypeBug.ValueString(),
		IssueTypeIncident:    plan.IssueTypeIncident.ValueString(),
		IssueTypeRequirement: plan.IssueTypeRequirement.ValueString(),
		Name:                 plan.Name.ValueString(),
		PrBodyClosePattern:   plan.PrBodyClosePattern.ValueString(),
		PrComponent:          plan.PrComponent.ValueString(),
		PrType:               plan.PrType.ValueString(),
		ProductionPattern:    plan.ProductionPattern.ValueString(),
		RefDiff: &client.RefDiff{
			TagsLimit:   int(plan.RefDiff.TagsLimit.ValueInt64()),
			TagsPattern: plan.RefDiff.TagsPattern.ValueString(),
		},
		UpdatedAt: now,
	}

	// Create new githubconnectionscopeconfig
	githubConnectionScopeConfig, err := r.client.CreateGithubConnectionScopeConfig(plan.ConnectionId.ValueString(), githubConnectionScopeConfigCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating devlake github connection scope config",
			"Could not create devlake github connection scope config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, githubConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Entities = entitiesVal
	plan.ID = types.StringValue(strconv.Itoa(githubConnectionScopeConfig.ID))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.CreatedAt = types.StringValue(githubConnectionScopeConfig.CreatedAt)
	plan.DeploymentPattern = types.StringValue(githubConnectionScopeConfig.DeploymentPattern)
	plan.EnvNamePattern = types.StringValue(githubConnectionScopeConfig.EnvNamePattern)
	plan.IssueComponent = types.StringValue(githubConnectionScopeConfig.IssueComponent)
	plan.IssuePriority = types.StringValue(githubConnectionScopeConfig.IssuePriority)
	plan.IssueSeverity = types.StringValue(githubConnectionScopeConfig.IssueSeverity)
	plan.IssueTypeBug = types.StringValue(githubConnectionScopeConfig.IssueTypeBug)
	plan.IssueTypeIncident = types.StringValue(githubConnectionScopeConfig.IssueTypeIncident)
	plan.IssueTypeRequirement = types.StringValue(githubConnectionScopeConfig.IssueTypeRequirement)
	plan.Name = types.StringValue(githubConnectionScopeConfig.Name)
	plan.PrBodyClosePattern = types.StringValue(githubConnectionScopeConfig.PrBodyClosePattern)
	plan.PrComponent = types.StringValue(githubConnectionScopeConfig.PrComponent)
	plan.PrType = types.StringValue(githubConnectionScopeConfig.PrType)
	plan.ProductionPattern = types.StringValue(githubConnectionScopeConfig.ProductionPattern)
	plan.RefDiff.TagsLimit = types.Int64Value(int64(githubConnectionScopeConfig.RefDiff.TagsLimit))
	plan.RefDiff.TagsPattern = types.StringValue(githubConnectionScopeConfig.RefDiff.TagsPattern)
	plan.UpdatedAt = types.StringValue(githubConnectionScopeConfig.UpdatedAt)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *githubConnectionScopeConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state githubConnectionScopeConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed github connection scope config value from Devlake
	githubConnectionScopeConfig, err := r.client.ReadGithubConnectionScopeConfig(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devlake github connection scopeconfig",
			err.Error(),
		)
		return
	}

	// Overwrite github connection scope config with refreshed state
	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, githubConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.CreatedAt = types.StringValue(githubConnectionScopeConfig.CreatedAt)
	state.DeploymentPattern = types.StringValue(githubConnectionScopeConfig.DeploymentPattern)
	state.Entities = entitiesVal
	state.EnvNamePattern = types.StringValue(githubConnectionScopeConfig.EnvNamePattern)
	state.ID = types.StringValue(strconv.Itoa(githubConnectionScopeConfig.ID))
	state.IssueComponent = types.StringValue(githubConnectionScopeConfig.IssueComponent)
	state.IssuePriority = types.StringValue(githubConnectionScopeConfig.IssuePriority)
	state.IssueSeverity = types.StringValue(githubConnectionScopeConfig.IssueSeverity)
	state.IssueTypeBug = types.StringValue(githubConnectionScopeConfig.IssueTypeBug)
	state.IssueTypeIncident = types.StringValue(githubConnectionScopeConfig.IssueTypeIncident)
	state.IssueTypeRequirement = types.StringValue(githubConnectionScopeConfig.IssueTypeRequirement)
	state.Name = types.StringValue(githubConnectionScopeConfig.Name)
	state.PrBodyClosePattern = types.StringValue(githubConnectionScopeConfig.PrBodyClosePattern)
	state.PrComponent = types.StringValue(githubConnectionScopeConfig.PrComponent)
	state.PrType = types.StringValue(githubConnectionScopeConfig.PrType)
	state.ProductionPattern = types.StringValue(githubConnectionScopeConfig.ProductionPattern)
	state.UpdatedAt = types.StringValue(githubConnectionScopeConfig.UpdatedAt)
	if apiRefDiff := githubConnectionScopeConfig.RefDiff; apiRefDiff != nil {
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
func (r *githubConnectionScopeConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan githubConnectionScopeConfigResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	connectionId, err := strconv.Atoi(plan.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scopeconfig",
			"Could not update devlake github connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scopeconfig",
			"Could not update devlake github connection scopeconfig, unexpected error: "+err.Error(),
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
	var githubConnectionScopeConfigUpdate = client.GithubConnectionScopeConfig{
		ConnectionId:         connectionId,
		ID:                   id,
		DeploymentPattern:    plan.DeploymentPattern.ValueString(),
		Entities:             entities,
		EnvNamePattern:       plan.EnvNamePattern.ValueString(),
		IssueComponent:       plan.IssueComponent.ValueString(),
		IssuePriority:        plan.IssuePriority.ValueString(),
		IssueSeverity:        plan.IssueSeverity.ValueString(),
		IssueTypeBug:         plan.IssueTypeBug.ValueString(),
		IssueTypeIncident:    plan.IssueTypeIncident.ValueString(),
		IssueTypeRequirement: plan.IssueTypeRequirement.ValueString(),
		Name:                 plan.Name.ValueString(),
		PrBodyClosePattern:   plan.PrBodyClosePattern.ValueString(),
		PrComponent:          plan.PrComponent.ValueString(),
		PrType:               plan.PrType.ValueString(),
		ProductionPattern:    plan.ProductionPattern.ValueString(),
		RefDiff: &client.RefDiff{
			TagsLimit:   int(plan.RefDiff.TagsLimit.ValueInt64()),
			TagsPattern: plan.RefDiff.TagsPattern.ValueString(),
		},
		UpdatedAt: time.Now().Format(time.RFC850),
	}

	// Update existing connection scope config
	updatedGithubConnectionScopeConfig, err := r.client.UpdateGithubConnectionScopeConfig(plan.ConnectionId.ValueString(), plan.ID.ValueString(), githubConnectionScopeConfigUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating devlake github connection scopeconfig",
			"Could not update devlake github connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}

	entitiesVal, diags := types.ListValueFrom(ctx, types.StringType, updatedGithubConnectionScopeConfig.Entities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.CreatedAt = types.StringValue(updatedGithubConnectionScopeConfig.CreatedAt)
	plan.DeploymentPattern = types.StringValue(updatedGithubConnectionScopeConfig.DeploymentPattern)
	plan.Entities = entitiesVal
	plan.EnvNamePattern = types.StringValue(updatedGithubConnectionScopeConfig.EnvNamePattern)
	plan.IssueComponent = types.StringValue(updatedGithubConnectionScopeConfig.IssueComponent)
	plan.IssuePriority = types.StringValue(updatedGithubConnectionScopeConfig.IssuePriority)
	plan.IssueSeverity = types.StringValue(updatedGithubConnectionScopeConfig.IssueSeverity)
	plan.IssueTypeBug = types.StringValue(updatedGithubConnectionScopeConfig.IssueTypeBug)
	plan.IssueTypeIncident = types.StringValue(updatedGithubConnectionScopeConfig.IssueTypeIncident)
	plan.IssueTypeRequirement = types.StringValue(updatedGithubConnectionScopeConfig.IssueTypeRequirement)
	plan.Name = types.StringValue(updatedGithubConnectionScopeConfig.Name)
	plan.PrBodyClosePattern = types.StringValue(updatedGithubConnectionScopeConfig.PrBodyClosePattern)
	plan.PrComponent = types.StringValue(updatedGithubConnectionScopeConfig.PrComponent)
	plan.PrType = types.StringValue(updatedGithubConnectionScopeConfig.PrType)
	plan.ProductionPattern = types.StringValue(updatedGithubConnectionScopeConfig.ProductionPattern)
	plan.RefDiff.TagsLimit = types.Int64Value(int64(updatedGithubConnectionScopeConfig.RefDiff.TagsLimit))
	plan.RefDiff.TagsPattern = types.StringValue(updatedGithubConnectionScopeConfig.RefDiff.TagsPattern)
	plan.UpdatedAt = types.StringValue(updatedGithubConnectionScopeConfig.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *githubConnectionScopeConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state githubConnectionScopeConfigResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing connection scope config
	err := r.client.DeleteGithubConnectionScopeConfig(state.ConnectionId.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting devlake github connection scopeconfig",
			"Could not delete devlake github connection scopeconfig, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *githubConnectionScopeConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *githubConnectionScopeConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
