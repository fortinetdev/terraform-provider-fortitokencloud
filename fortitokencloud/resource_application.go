package fortitokencloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
	ftc_client "terraform-provider-fortitokencloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

// NewApplicationResource is a helper function to simplify the provider implementation.
func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

// applicationResource is the resource implementation.
type applicationResource struct {
	client *ftc_client.Client
}

func formatAppObj(plan applicationResourceModel, create bool) (*map[string]interface{}, map[string][]string) {
	obj := make(map[string]interface{})
	user_source_list := make(map[string][]string)
	obj["name"] = plan.Name.ValueString()
	obj["ttl"] = plan.TTL.ValueInt64()
	if create {
		obj["realm_id"] = plan.RealmID.ValueString()
	}
	if plan.BrandingID.ValueString() == "" {
		obj["branding_id"] = nil
	} else {
		obj["branding_id"] = plan.BrandingID.ValueString()
	}
	if plan.AttrMapping.ValueString() == "" {
		obj["attr_mapping"] = nil
	} else {
		var attr_mapping map[string]interface{}
		err := json.Unmarshal([]byte(plan.AttrMapping.ValueString()), &attr_mapping)
		if err != nil {
			return nil, user_source_list
		}
		obj["attr_mapping"] = attr_mapping
	}
	saml_obj := map[string]interface{}{}
	if plan.SigningCertID.ValueString() == "" {
		saml_obj["signing_cert_id"] = nil
	} else {
		saml_obj["signing_cert_id"] = plan.SigningCertID.ValueString()
	}
	if plan.SPEntityID.ValueString() == "" {
		saml_obj["entity_id"] = nil
	} else {
		saml_obj["entity_id"] = plan.SPEntityID.ValueString()
	}
	if plan.SPAcsURL.ValueString() == "" {
		saml_obj["acs_url"] = nil
	} else {
		saml_obj["acs_url"] = plan.SPAcsURL.ValueString()
	}
	if plan.SPSloURL.ValueString() == "" {
		saml_obj["slo_url"] = nil
	} else {
		saml_obj["slo_url"] = plan.SPSloURL.ValueString()
	}
	obj["saml_params"] = saml_obj

	var user_source_ids []string
	for _, user_source_id := range plan.UserSources {
		user_source_ids = append(user_source_ids, user_source_id.ValueString())
	}
	user_source_list["user_source_ids"] = user_source_ids
	return &obj, user_source_list
}

// Metadata returns the resource type name.
func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Configure adds the provider configured client to the resource.
func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ftc_client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ftc_client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"entity_id": schema.StringAttribute{
				Computed: true,
			},
			"sso_url": schema.StringAttribute{
				Computed: true,
			},
			"slo_url": schema.StringAttribute{
				Computed: true,
			},
			"realm_id": schema.StringAttribute{
				Required: true,
			},
			"prefix": schema.StringAttribute{
				Computed: true,
			},
			"branding_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"ttl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"signing_cert_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"sp_entity_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"sp_acs_url": schema.StringAttribute{
				Optional: true,
			},
			"sp_slo_url": schema.StringAttribute{
				Optional: true,
			},
			"user_source_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Computed:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"attr_mapping": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// applicationResourceModel maps the resource schema data.
type applicationResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	EntityID      types.String   `tfsdk:"entity_id"`
	SsoUrl        types.String   `tfsdk:"sso_url"`
	SloUrl        types.String   `tfsdk:"slo_url"`
	RealmID       types.String   `tfsdk:"realm_id"`
	Prefix        types.String   `tfsdk:"prefix"`
	BrandingID    types.String   `tfsdk:"branding_id"`
	TTL           types.Int64    `tfsdk:"ttl"`
	SigningCertID types.String   `tfsdk:"signing_cert_id"`
	SPEntityID    types.String   `tfsdk:"sp_entity_id"`
	SPAcsURL      types.String   `tfsdk:"sp_acs_url"`
	SPSloURL      types.String   `tfsdk:"sp_slo_url"`
	UserSources   []types.String `tfsdk:"user_source_ids"`
	AttrMapping   types.String   `tfsdk:"attr_mapping"`
}

// Create a new resource.
func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan applicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, user_source_ids := formatAppObj(plan, true)

	// Create new application
	application, err := r.client.CreateApplication(obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	attr_mapping, _ := json.Marshal(application.AttrMapping)
	plan = applicationResourceModel{
		ID:            types.StringValue(application.ID),
		Name:          types.StringValue(application.Name),
		EntityID:      types.StringValue(application.EntityID),
		SloUrl:        types.StringValue(application.SloUrl),
		SsoUrl:        types.StringValue(application.SsoUrl),
		RealmID:       types.StringValue(application.RealmID),
		Prefix:        types.StringValue(application.Prefix),
		BrandingID:    types.StringValue(application.BrandingID),
		TTL:           types.Int64Value(int64(application.TTL)),
		SigningCertID: types.StringValue(application.SigningCertID),
		SPEntityID:    types.StringValue(application.SpEntityID),
		SPAcsURL:      types.StringValue(application.SpAcsUrl),
		SPSloURL:      types.StringValue(application.SpSloUrl),
		AttrMapping:   types.StringValue(string(attr_mapping)),
		UserSources:   make([]types.String, 0),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	if len(user_source_ids["user_source_ids"]) > 0 {
		_, err = r.client.UpdateApplicationUserSource(application.ID, user_source_ids)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating app user source mapping",
				"Could not create app user source mapping, unexpected error: "+err.Error(),
			)
			return
		}
		for _, user_source_id := range user_source_ids["user_source_ids"] {
			plan.UserSources = append(plan.UserSources, types.StringValue(user_source_id))
		}
		diags = resp.State.Set(ctx, plan)
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state applicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if resource has been deleted upstream, reset state
	if state.ID.ValueString() == "" {
		state = applicationResourceModel{}
	} else {
		application, err := r.client.GetApplication(state.ID.ValueString())
		if err != nil {
			if !strings.Contains(err.Error(), "status: 404") {
				resp.Diagnostics.AddError(
					"Error Reading application",
					"Could not read application ID "+state.ID.ValueString()+": "+err.Error(),
				)
				return
			}
			state = applicationResourceModel{}
		} else {
			new_user_sources := make([]types.String, 0)
			// Overwrite items with refreshed state
			attr_mapping, _ := json.Marshal(application.AttrMapping)
			for _, usersource := range application.UserSources {
				new_user_sources = append(new_user_sources, types.StringValue(usersource.ID))
			}
			state = applicationResourceModel{
				ID:            types.StringValue(application.ID),
				Name:          types.StringValue(application.Name),
				EntityID:      types.StringValue(application.EntityID),
				SloUrl:        types.StringValue(application.SloUrl),
				SsoUrl:        types.StringValue(application.SsoUrl),
				RealmID:       types.StringValue(application.RealmID),
				Prefix:        types.StringValue(application.Prefix),
				BrandingID:    types.StringValue(application.BrandingID),
				TTL:           types.Int64Value(int64(application.TTL)),
				SigningCertID: types.StringValue(application.SigningCertID),
				SPEntityID:    types.StringValue(application.SpEntityID),
				SPAcsURL:      types.StringValue(application.SpAcsUrl),
				SPSloURL:      types.StringValue(application.SpSloUrl),
				AttrMapping:   types.StringValue(string(attr_mapping)),
				UserSources:   new_user_sources,
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

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan applicationResourceModel
	var state applicationResourceModel
	req.State.Get(ctx, &state)
	var application *ftc_client.Application
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	old_user_sources := state.UserSources
	new_user_sources := plan.UserSources

	obj, user_source_list := formatAppObj(plan, false)
	// if resource was deleted upstream, recreate it.
	if plan.ID.ValueString() == "" {
		obj, user_source_list = formatAppObj(plan, true)
		application, err = r.client.CreateApplication(obj)
	} else {
		// Update existing application
		application, err = r.client.UpdateApplication(plan.ID.ValueString(), obj)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating application",
			"Could not update application, unexpected error: "+err.Error(),
		)
		return
	}

	attr_mapping, _ := json.Marshal(application.AttrMapping)
	plan = applicationResourceModel{
		ID:            types.StringValue(application.ID),
		Name:          types.StringValue(application.Name),
		EntityID:      types.StringValue(application.EntityID),
		SloUrl:        types.StringValue(application.SloUrl),
		SsoUrl:        types.StringValue(application.SsoUrl),
		RealmID:       types.StringValue(application.RealmID),
		Prefix:        types.StringValue(application.Prefix),
		BrandingID:    types.StringValue(application.BrandingID),
		TTL:           types.Int64Value(int64(application.TTL)),
		SigningCertID: types.StringValue(application.SigningCertID),
		SPEntityID:    types.StringValue(application.SpEntityID),
		SPAcsURL:      types.StringValue(application.SpAcsUrl),
		SPSloURL:      types.StringValue(application.SpSloUrl),
		AttrMapping:   types.StringValue(string(attr_mapping)),
		UserSources:   old_user_sources,
	}

	if !reflect.DeepEqual(new_user_sources, old_user_sources) {
		_, err = r.client.UpdateApplicationUserSource(plan.ID.ValueString(), user_source_list)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating application",
				"Could not update application, unexpected error: "+err.Error(),
			)
			return
		}
	}
	plan.UserSources = new_user_sources
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state applicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing application
	if state.ID.ValueString() != "" {
		err := r.client.DeleteApplication(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Application",
				"Could not delete application, unexpected error: "+err.Error(),
			)
		}
		return
	}
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
