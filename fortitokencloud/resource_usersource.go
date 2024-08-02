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
	_            resource.Resource                = &userSourceResource{}
	_            resource.ResourceWithConfigure   = &userSourceResource{}
	_            resource.ResourceWithImportState = &userSourceResource{}
	type_int_map                                  = map[int64]string{
		1: "saml",
		2: "oidc",
	}
	type_str_map = map[string]int64{
		"saml": 1,
		"oidc": 2,
	}
)

// NewUserSourceResource is a helper function to simplify the provider implementation.
func NewUserSourceResource() resource.Resource {
	return &userSourceResource{}
}

// UserSourceResource is the resource implementation.
type userSourceResource struct {
	client *ftc_client.Client
}

func formatUsObj(plan userSourceResourceModel, create bool) (*map[string]interface{}, map[string][]string) {
	obj := make(map[string]interface{})
	domain_list := make(map[string][]string)
	obj["name"] = plan.Name.ValueString()
	if create {
		obj["realm_id"] = plan.RealmID.ValueString()
		obj["type"] = type_str_map[strings.ToLower(plan.Type.ValueString())]
	}
	if plan.UsernameAssertion.ValueString() == "" {
		obj["username_assertion"] = nil
	} else {
		obj["username_assertion"] = plan.UsernameAssertion.ValueString()
	}
	if plan.LoginHint.ValueString() == "" {
		obj["login_hint"] = nil
	} else {
		obj["login_hint"] = plan.LoginHint.ValueString()
	}
	if plan.AttrMapping.ValueString() == "" {
		obj["attr_mapping"] = map[string]interface{}{}
	} else {
		var attr_mapping map[string]interface{}
		err := json.Unmarshal([]byte(plan.AttrMapping.ValueString()), &attr_mapping)
		if err != nil {
			return nil, domain_list
		}
		obj["attr_mapping"] = attr_mapping
	}
	if strings.ToLower(plan.Type.ValueString()) == "saml" {
		saml_obj := map[string]interface{}{}
		if plan.EntityID.ValueString() == "" {
			saml_obj["entity_id"] = nil
		} else {
			saml_obj["entity_id"] = plan.EntityID.ValueString()
		}
		if plan.LogoutUrl.ValueString() == "" {
			saml_obj["login_url"] = nil
		} else {
			saml_obj["login_url"] = plan.LoginUrl.ValueString()
		}
		if plan.LogoutUrl.ValueString() == "" {
			saml_obj["logout_url"] = nil
		} else {
			saml_obj["logout_url"] = plan.LogoutUrl.ValueString()
		}
		saml_obj["post_binding"] = plan.PostBinding.ValueBool()
		saml_obj["include_subject"] = plan.IncludeSubject.ValueBool()
		obj["saml_params"] = saml_obj
	} else {
		odic_obj := map[string]interface{}{}
		if plan.AuthUri.ValueString() == "" {
			odic_obj["auth_uri"] = nil
		} else {
			odic_obj["auth_uri"] = plan.AuthUri.ValueString()
		}
		if plan.TokenUri.ValueString() == "" {
			odic_obj["token_uri"] = nil
		} else {
			odic_obj["token_uri"] = plan.TokenUri.ValueString()
		}
		if plan.UserInfoUri.ValueString() == "" {
			odic_obj["userinfo_uri"] = nil
		} else {
			odic_obj["userinfo_uri"] = plan.UserInfoUri.ValueString()
		}
		if plan.LogoutUri.ValueString() == "" {
			odic_obj["logout_uri"] = nil
		} else {
			odic_obj["logout_uri"] = plan.LogoutUri.ValueString()
		}
		if plan.Issuer.ValueString() == "" {
			odic_obj["issuer"] = nil
		} else {
			odic_obj["issuer"] = plan.Issuer.ValueString()
		}
		if plan.ClientID.ValueString() == "" {
			odic_obj["client_id"] = nil
		} else {
			odic_obj["client_id"] = plan.ClientID.ValueString()
		}
		if plan.ClientSecret.ValueString() == "" {
			odic_obj["client_secret"] = nil
		} else {
			odic_obj["client_secret"] = plan.ClientSecret.ValueString()
		}
		obj["oidc_params"] = odic_obj
	}

	var domain_ids []string
	for _, domain_id := range plan.Domains {
		domain_ids = append(domain_ids, domain_id.ValueString())
	}
	domain_list["domain_ids"] = domain_ids
	return &obj, domain_list
}

// Metadata returns the resource type name.
func (r *userSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usersource"
}

// Configure adds the provider configured client to the resource.
func (r *userSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *userSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"type": schema.StringAttribute{
				Required: true,
			},
			"entity_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"login_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"logout_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"auth_uri": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"token_uri": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"userinfo_uri": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"logout_uri": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"issuer": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"client_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"client_secret": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"realm_id": schema.StringAttribute{
				Required: true,
			},
			"prefix": schema.StringAttribute{
				Computed: true,
			},
			"post_binding": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_subject": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"username_assertion": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("username"),
			},
			"login_hint": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"attr_mapping": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"domain_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"proxy_entity_id": schema.StringAttribute{
				Computed: true,
			},
			"proxy_acs_url": schema.StringAttribute{
				Computed: true,
			},
			"proxy_slo_url": schema.StringAttribute{
				Computed: true,
			},
			"proxy_sso_url": schema.StringAttribute{
				Computed: true,
			},
			"proxy_callback_url": schema.StringAttribute{
				Computed: true,
			},
			"proxy_post_logout_redirect_uri": schema.StringAttribute{
				Computed: true,
			},
			"proxy_oidc_login_url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// userSourceResourceModel maps the resource schema data.
type userSourceResourceModel struct {
	ID                         types.String   `tfsdk:"id"`
	Name                       types.String   `tfsdk:"name"`
	Type                       types.String   `tfsdk:"type"`
	EntityID                   types.String   `tfsdk:"entity_id"`
	LoginUrl                   types.String   `tfsdk:"login_url"`
	LogoutUrl                  types.String   `tfsdk:"logout_url"`
	AuthUri                    types.String   `tfsdk:"auth_uri"`
	TokenUri                   types.String   `tfsdk:"token_uri"`
	UserInfoUri                types.String   `tfsdk:"userinfo_uri"`
	LogoutUri                  types.String   `tfsdk:"logout_uri"`
	Issuer                     types.String   `tfsdk:"issuer"`
	ClientID                   types.String   `tfsdk:"client_id"`
	ClientSecret               types.String   `tfsdk:"client_secret"`
	RealmID                    types.String   `tfsdk:"realm_id"`
	Prefix                     types.String   `tfsdk:"prefix"`
	PostBinding                types.Bool     `tfsdk:"post_binding"`
	IncludeSubject             types.Bool     `tfsdk:"include_subject"`
	UsernameAssertion          types.String   `tfsdk:"username_assertion"`
	LoginHint                  types.String   `tfsdk:"login_hint"`
	Domains                    []types.String `tfsdk:"domain_ids"`
	ProxyEntityID              types.String   `tfsdk:"proxy_entity_id"`
	ProxyAcsUrl                types.String   `tfsdk:"proxy_acs_url"`
	ProxySloUrl                types.String   `tfsdk:"proxy_slo_url"`
	ProxySSoUrl                types.String   `tfsdk:"proxy_sso_url"`
	ProxyCallbackUrl           types.String   `tfsdk:"proxy_callback_url"`
	ProxyPostLogoutRedirectUri types.String   `tfsdk:"proxy_post_logout_redirect_uri"`
	ProxyOidcLoginUrl          types.String   `tfsdk:"proxy_oidc_login_url"`
	AttrMapping                types.String   `tfsdk:"attr_mapping"`
}

// Create a new resource.
func (r *userSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	var new_dms = plan.Domains
	if resp.Diagnostics.HasError() {
		return
	}

	obj, domain_ids := formatUsObj(plan, true)

	// Create new user source
	usersource, err := r.client.CreateUserSource(obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user source",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	usersource, err = r.client.GetUserSource(usersource.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading user source",
			"Could not read user source ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	attr_mapping, _ := json.Marshal(usersource.AttrMapping)
	plan = userSourceResourceModel{
		ID:                         types.StringValue(usersource.ID),
		Name:                       types.StringValue(usersource.Name),
		Type:                       types.StringValue(type_int_map[int64(usersource.Type)]),
		EntityID:                   types.StringValue(usersource.EntityID),
		LoginUrl:                   types.StringValue(usersource.LoginUrl),
		LogoutUrl:                  types.StringValue(usersource.LogoutUrl),
		AuthUri:                    types.StringValue(usersource.AuthUri),
		TokenUri:                   types.StringValue(usersource.TokenUri),
		UserInfoUri:                types.StringValue(usersource.UserInfoUri),
		LogoutUri:                  types.StringValue(usersource.LogoutUri),
		Issuer:                     types.StringValue(usersource.Issuer),
		ClientID:                   types.StringValue(usersource.ClientID),
		ClientSecret:               types.StringValue(plan.ClientSecret.ValueString()),
		RealmID:                    types.StringValue(usersource.RealmID),
		Prefix:                     types.StringValue(usersource.Prefix),
		PostBinding:                types.BoolValue(usersource.PostBinding),
		IncludeSubject:             types.BoolValue(usersource.IncludeSubject),
		UsernameAssertion:          types.StringValue(usersource.UsernameAssertion),
		LoginHint:                  types.StringValue(usersource.LoginHint),
		ProxyEntityID:              types.StringValue(usersource.ProxySP.EntityID),
		ProxyAcsUrl:                types.StringValue(usersource.ProxySP.AcsUrl),
		ProxySloUrl:                types.StringValue(usersource.ProxySP.SloUrl),
		ProxySSoUrl:                types.StringValue(usersource.ProxySP.SsoUrl),
		ProxyCallbackUrl:           types.StringValue(usersource.ProxySP.CallbackUrl),
		ProxyPostLogoutRedirectUri: types.StringValue(usersource.ProxySP.PostLogoutRedirectUrl),
		ProxyOidcLoginUrl:          types.StringValue(usersource.ProxySP.OidcLoginUrl),
		AttrMapping:                types.StringValue(string(attr_mapping)),
		Domains:                    make([]types.String, 0),
	}

	// Set state to fully populated data
	resp.State.Set(ctx, plan)

	_, err = r.client.UpdateUserSourceDomains(usersource.ID, domain_ids)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app user source mapping",
			"Could not create app user source mapping, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Domains = new_dms
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *userSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if resource has been deleted upstream, reset state
	if state.ID.ValueString() == "" {
		state = userSourceResourceModel{}
	} else {
		// Get refreshed user source
		usersource, err := r.client.GetUserSource(state.ID.ValueString())
		if err != nil {
			if !strings.Contains(err.Error(), "status: 404") {
				resp.Diagnostics.AddError(
					"Error Reading user source",
					"Could not read User Source ID "+state.ID.ValueString()+": "+err.Error(),
				)
				return
			}
			state = userSourceResourceModel{}
		} else {
			// Overwrite items with refreshed state
			attr_mapping, _ := json.Marshal(usersource.AttrMapping)
			domain_list := make([]types.String, 0)
			for _, domain := range usersource.Domains {
				domain_list = append(domain_list, types.StringValue(domain.ID))
			}
			new_secret := state.ClientSecret
			if state.ClientSecret.ValueString() == "" {
				new_secret = types.StringValue(usersource.ClientSecret)
			}
			state = userSourceResourceModel{
				ID:                         types.StringValue(usersource.ID),
				Name:                       types.StringValue(usersource.Name),
				Type:                       types.StringValue(type_int_map[int64(usersource.Type)]),
				EntityID:                   types.StringValue(usersource.EntityID),
				LoginUrl:                   types.StringValue(usersource.LoginUrl),
				LogoutUrl:                  types.StringValue(usersource.LogoutUrl),
				AuthUri:                    types.StringValue(usersource.AuthUri),
				TokenUri:                   types.StringValue(usersource.TokenUri),
				UserInfoUri:                types.StringValue(usersource.UserInfoUri),
				LogoutUri:                  types.StringValue(usersource.LogoutUri),
				Issuer:                     types.StringValue(usersource.Issuer),
				ClientID:                   types.StringValue(usersource.ClientID),
				ClientSecret:               new_secret,
				RealmID:                    types.StringValue(usersource.RealmID),
				Prefix:                     types.StringValue(usersource.Prefix),
				PostBinding:                types.BoolValue(usersource.PostBinding),
				IncludeSubject:             types.BoolValue(usersource.IncludeSubject),
				UsernameAssertion:          types.StringValue(usersource.UsernameAssertion),
				LoginHint:                  types.StringValue(usersource.LoginHint),
				ProxyEntityID:              types.StringValue(usersource.ProxySP.EntityID),
				ProxyAcsUrl:                types.StringValue(usersource.ProxySP.AcsUrl),
				ProxySloUrl:                types.StringValue(usersource.ProxySP.SloUrl),
				ProxySSoUrl:                types.StringValue(usersource.ProxySP.SsoUrl),
				ProxyCallbackUrl:           types.StringValue(usersource.ProxySP.CallbackUrl),
				ProxyPostLogoutRedirectUri: types.StringValue(usersource.ProxySP.PostLogoutRedirectUrl),
				ProxyOidcLoginUrl:          types.StringValue(usersource.ProxySP.OidcLoginUrl),
				AttrMapping:                types.StringValue(string(attr_mapping)),
				Domains:                    domain_list,
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

func (r *userSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userSourceResourceModel
	var state userSourceResourceModel
	var usersource *ftc_client.UserSource
	var err error
	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	old_domains := state.Domains
	new_domains := plan.Domains

	obj, domain_ids := formatUsObj(plan, false)

	if plan.ID.ValueString() == "" {
		// Create new user source if not found
		obj, domain_ids = formatUsObj(plan, true)
		usersource, err = r.client.CreateUserSource(obj)
	} else {
		// Update existing user source
		usersource, err = r.client.UpdateUserSource(plan.ID.ValueString(), obj)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating user source",
			"Could not update application, unexpected error: "+err.Error(),
		)
		return
	}

	usersource, err = r.client.GetUserSource(usersource.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading user source",
			"Could not read user source ID "+usersource.ID+": "+err.Error(),
		)
		return
	}

	attr_mapping, _ := json.Marshal(usersource.AttrMapping)
	plan = userSourceResourceModel{
		ID:                         types.StringValue(usersource.ID),
		Name:                       types.StringValue(usersource.Name),
		Type:                       types.StringValue(type_int_map[int64(usersource.Type)]),
		EntityID:                   types.StringValue(usersource.EntityID),
		LoginUrl:                   types.StringValue(usersource.LoginUrl),
		LogoutUrl:                  types.StringValue(usersource.LogoutUrl),
		AuthUri:                    types.StringValue(usersource.AuthUri),
		TokenUri:                   types.StringValue(usersource.TokenUri),
		UserInfoUri:                types.StringValue(usersource.UserInfoUri),
		LogoutUri:                  types.StringValue(usersource.LogoutUri),
		Issuer:                     types.StringValue(usersource.Issuer),
		ClientID:                   types.StringValue(usersource.ClientID),
		ClientSecret:               types.StringValue(plan.ClientSecret.ValueString()),
		RealmID:                    types.StringValue(usersource.RealmID),
		Prefix:                     types.StringValue(usersource.Prefix),
		PostBinding:                types.BoolValue(usersource.PostBinding),
		IncludeSubject:             types.BoolValue(usersource.IncludeSubject),
		UsernameAssertion:          types.StringValue(usersource.UsernameAssertion),
		LoginHint:                  types.StringValue(usersource.LoginHint),
		ProxyEntityID:              types.StringValue(usersource.ProxySP.EntityID),
		ProxyAcsUrl:                types.StringValue(usersource.ProxySP.AcsUrl),
		ProxySloUrl:                types.StringValue(usersource.ProxySP.SloUrl),
		ProxySSoUrl:                types.StringValue(usersource.ProxySP.SsoUrl),
		ProxyCallbackUrl:           types.StringValue(usersource.ProxySP.CallbackUrl),
		ProxyPostLogoutRedirectUri: types.StringValue(usersource.ProxySP.PostLogoutRedirectUrl),
		ProxyOidcLoginUrl:          types.StringValue(usersource.ProxySP.OidcLoginUrl),
		AttrMapping:                types.StringValue(string(attr_mapping)),
		Domains:                    old_domains,
	}

	if !reflect.DeepEqual(new_domains, old_domains) {
		_, err = r.client.UpdateUserSourceDomains(plan.ID.ValueString(), domain_ids)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating application",
				"Could not update application, unexpected error: "+err.Error(),
			)
			return
		}
		plan.Domains = new_domains
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing user source
	if state.ID.ValueString() != "" {
		err := r.client.DeleteUserSource(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting UserSource",
				"Could not delete users ource, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *userSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
