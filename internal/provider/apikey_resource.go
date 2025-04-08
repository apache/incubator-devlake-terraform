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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &apiKeyResource{}
	_ resource.ResourceWithConfigure   = &apiKeyResource{}
	_ resource.ResourceWithImportState = &apiKeyResource{}
)

// NewApiKeyResource is a helper function to simplify the provider implementation.
func NewApiKeyResource() resource.Resource {
	return &apiKeyResource{}
}

// apiKeyResource is the resource implementation.
type apiKeyResource struct {
	client *client.Client
}

// apiKeyResourceModel maps the resource schema data.
type apiKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	AllowedPath types.String `tfsdk:"allowed_path"`
	ApiKey      types.String `tfsdk:"api_key"`
	ExpiredAt   types.String `tfsdk:"expired_at"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
}

// Metadata returns the resource type name.
func (r *apiKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikey"
}

// Schema defines the schema for the resource.
func (r *apiKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Numeric identifier for the apikey. This is a string for easier resource import.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last Terraform update of the apikey.",
			},
			"allowed_path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_key": schema.StringAttribute{
				Computed:    true,
				Description: "The API URL or endpoint that the API key is permitted to access. It defines the specific resources that the key can interact with.",
				Sensitive:   true,
			},
			"expired_at": schema.StringAttribute{
				Description: "When the apikey expires.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the apikey.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Default:     stringdefault.StaticString("devlake"),
				Description: "The apikey type. Currently only 'devlake' is a valid value.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create a new resource.
func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var apiKeyCreate = client.ApiKeyCreate{
		AllowedPath: plan.AllowedPath.ValueString(),
		ExpiredAt:   plan.ExpiredAt.ValueString(),
		Name:        plan.Name.ValueString(),
		Type:        plan.Type.ValueString(),
	}

	// Create new apikey
	apiKey, err := r.client.CreateApiKey(apiKeyCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating apiKey",
			"Could not create apiKey, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(apiKey.ID))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.AllowedPath = types.StringValue(apiKey.AllowedPath)
	plan.ApiKey = types.StringValue(apiKey.ApiKey)
	plan.ExpiredAt = types.StringValue(apiKey.ExpiredAt)
	plan.Name = types.StringValue(apiKey.Name)
	plan.Type = types.StringValue(apiKey.Type)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed apikey value from Devlake
	apiKeys, err := r.client.GetApiKeys()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Devlake ApiKeys",
			err.Error(),
		)
		return
	}

	// Overwrite apikey with refreshed state
	for _, apiKey := range apiKeys {
		if types.StringValue(strconv.Itoa(apiKey.ID)) == state.ID {
			state = apiKeyResourceModel{
				ID:          types.StringValue(strconv.Itoa(apiKey.ID)),
				AllowedPath: types.StringValue(apiKey.AllowedPath),
				ExpiredAt:   types.StringValue(apiKey.ExpiredAt),
				Name:        types.StringValue(apiKey.Name),
				Type:        types.StringValue(apiKey.Type),
			}
			break
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// all modifyable fields are set to requires replace so no update is needed
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing apikey
	err := r.client.DeleteApiKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting apikey",
			"Could not delete apikey, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *apiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *apiKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
