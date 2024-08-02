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
var (
	_ datasource.DataSource              = &realmDataSource{}
	_ datasource.DataSourceWithConfigure = &realmDataSource{}
)

func NewRealmDataSource() datasource.DataSource {
	return &realmDataSource{}
}

func (d *realmDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm"
}

// Read refreshes the Terraform state with the latest data.
func (d *realmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data realmDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	realm, err := d.client.GetRealmByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read realm",
			err.Error(),
		)
		return
	}

	data = realmDataSourceModel{
		ID:   types.StringValue(realm.ID),
		Name: types.StringValue(realm.Name),
	}

	// Set state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// realmDataSource is the data source implementation.
type realmDataSource struct {
	client *ftc_client.Client
}

// Configure adds the provider configured client to the data source.
func (d *realmDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *realmDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// realmDataSourceModel maps coffees schema data.
type realmDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
