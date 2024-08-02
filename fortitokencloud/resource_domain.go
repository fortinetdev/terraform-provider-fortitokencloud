package fortitokencloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	ftc_client "terraform-provider-fortitokencloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &domainResource{}
	_ resource.ResourceWithConfigure   = &domainResource{}
	_ resource.ResourceWithImportState = &domainResource{}
)

// NewUserSourceResource is a helper function to simplify the provider implementation.
func NewDomainResource() resource.Resource {
	return &domainResource{}
}

// UserSourceResource is the resource implementation.
type domainResource struct {
	client *ftc_client.Client
}

func formatDomainObj(plan domainResourceModel, create bool) *map[string]interface{} {
	obj := make(map[string]interface{})
	obj["name"] = plan.Name.ValueString()
	if create {
		obj["realm_id"] = plan.RealmID.ValueString()
	}
	if plan.UserSourceID.ValueString() == "" {
		obj["user_source_id"] = nil
	} else {
		obj["user_source_id"] = plan.UserSourceID.ValueString()
	}
	return &obj
}

// Metadata returns the resource type name.
func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Configure adds the provider configured client to the resource.
func (r *domainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ftc_client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ftc.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"realm_id": schema.StringAttribute{
				Required: true,
			},
			"user_source_id": schema.StringAttribute{
				Computed: true,
				Required: false,
			},
		},
	}
}

// userSourceResourceModel maps the resource schema data.
type domainResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	RealmID      types.String `tfsdk:"realm_id"`
	UserSourceID types.String `tfsdk:"user_source_id"`
}

// Create a new resource.
func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan domainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := formatDomainObj(plan, true)

	// Create new domain
	domain, err := r.client.CreateDomain(obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating domain",
			"Could not create domain, unexpected error: "+err.Error(),
		)
		return
	}

	plan = domainResourceModel{
		ID:           types.StringValue(domain.ID),
		Name:         types.StringValue(domain.Name),
		RealmID:      types.StringValue(domain.RealmID),
		UserSourceID: types.StringValue(domain.UserSourceID),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state domainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if resource has been deleted upstream, reset state
	if state.ID.ValueString() == "" {
		state = domainResourceModel{}
	} else {
		// Get refreshed user source
		usersource, err := r.client.GetDomain(state.ID.ValueString())
		if err != nil {
			if !strings.Contains(err.Error(), "status: 404") {
				resp.Diagnostics.AddError(
					"Error Reading domain",
					"Could not read domain ID "+state.ID.ValueString()+": "+err.Error(),
				)
				return
			}
			state = domainResourceModel{}
		} else {

			// Overwrite items with refreshed state
			state = domainResourceModel{
				ID:           types.StringValue(usersource.ID),
				Name:         types.StringValue(usersource.Name),
				RealmID:      types.StringValue(usersource.RealmID),
				UserSourceID: types.StringValue(usersource.UserSourceID),
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan domainResourceModel
	var domain *ftc_client.Domain
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := formatDomainObj(plan, false)

	if plan.ID.ValueString() == "" {
		obj = formatDomainObj(plan, true)
		domain, err = r.client.CreateDomain(obj)
	} else {
		domain, err = r.client.UpdateDomain(plan.ID.ValueString(), obj)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating domain",
			"Could not update domain, unexpected error: "+err.Error(),
		)
		return
	}

	plan = domainResourceModel{
		ID:           types.StringValue(domain.ID),
		Name:         types.StringValue(domain.Name),
		RealmID:      types.StringValue(domain.RealmID),
		UserSourceID: types.StringValue(domain.UserSourceID),
	}

	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state domainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing domain
	if state.ID.ValueString() != "" {
		err := r.client.DeleteDomain(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting domain",
				"Could not delete domain, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
