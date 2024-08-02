package fortitokencloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ftc_client "terraform-provider-fortitokencloud/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &applicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationsDataSource{}
)

func NewApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

func (d *applicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

// Read refreshes the Terraform state with the latest data.
func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state applicationsDataSourceModel
	apps, err := d.client.GetApplications()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read FTC Applications",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, app := range apps.Apps {
		appState := applicationModel{
			ID:            types.StringValue(app.ID),
			Name:          types.StringValue(app.Name),
			EntityID:      types.StringValue(app.EntityID),
			SsoUrl:        types.StringValue(app.SsoUrl),
			SloUrl:        types.StringValue(app.SloUrl),
			RealmID:       types.StringValue(app.RealmID),
			Type:          types.Int64Value(int64(app.Type)),
			Prefix:        types.StringValue(app.Prefix),
			BrandingID:    types.StringValue(app.BrandingID),
			TTL:           types.Int64Value(int64(app.TTL)),
			SigningCertID: types.StringValue(app.SigningCertID),
			SPEntityID:    types.StringValue(app.SpEntityID),
			SPAcsURL:      types.StringValue(app.SpAcsUrl),
			SPSloURL:      types.StringValue(app.SpSloUrl),
		}

		state.Apps = append(state.Apps, appState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// applicationsDataSource is the data source implementation.
type applicationsDataSource struct {
	client *ftc_client.Client
}

// Configure adds the provider configured client to the data source.
func (d *applicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Schema defines the schema for the data source.
func (d *applicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apps": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
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
							Computed: true,
						},
						"type": schema.Int64Attribute{
							Computed: true,
						},
						"prefix": schema.StringAttribute{
							Computed: true,
						},
						"branding_id": schema.StringAttribute{
							Computed: true,
						},
						"ttl": schema.Int64Attribute{
							Computed: true,
						},
						"signing_cert_id": schema.StringAttribute{
							Computed: true,
						},
						"sp_entity_id": schema.StringAttribute{
							Computed: true,
						},
						"sp_acs_url": schema.StringAttribute{
							Computed: true,
						},
						"sp_slo_url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// applicationsDataSourceModel maps the data source schema data.
type applicationsDataSourceModel struct {
	Apps []applicationModel `tfsdk:"apps"`
}

// applicationModel maps application schema data.
type applicationModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	EntityID      types.String `tfsdk:"entity_id"`
	SsoUrl        types.String `tfsdk:"sso_url"`
	SloUrl        types.String `tfsdk:"slo_url"`
	RealmID       types.String `tfsdk:"realm_id"`
	Type          types.Int64  `tfsdk:"type"`
	Prefix        types.String `tfsdk:"prefix"`
	BrandingID    types.String `tfsdk:"branding_id"`
	TTL           types.Int64  `tfsdk:"ttl"`
	SigningCertID types.String `tfsdk:"signing_cert_id"`
	SPEntityID    types.String `tfsdk:"sp_entity_id"`
	SPAcsURL      types.String `tfsdk:"sp_acs_url"`
	SPSloURL      types.String `tfsdk:"sp_slo_url"`
}
