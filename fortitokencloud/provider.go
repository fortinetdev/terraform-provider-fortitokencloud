package fortitokencloud

import (
	"context"
	"os"
	ftc_client "terraform-provider-fortitokencloud/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &fortiTokenCloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &fortiTokenCloudProvider{
			version: version,
		}
	}
}

// fortiTokenCloudProvider is the provider implementation.
type fortiTokenCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Schema defines the provider-level schema for configuration data.
func (p *fortiTokenCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"clientid": schema.StringAttribute{
				Optional: true,
			},
			"clientsecret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// fortiTokenCloudProviderModel maps provider schema data to a Go type.
type fortiTokenCloudProviderModel struct {
	Host         types.String `tfsdk:"host"`
	ClientId     types.String `tfsdk:"clientid"`
	ClientSecret types.String `tfsdk:"clientsecret"`
}

// Metadata returns the provider type name.
func (p *fortiTokenCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fortitokencloud"
	resp.Version = p.version
}

func (p *fortiTokenCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config fortiTokenCloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown FortiTokenCloud API Host",
			"The provider cannot create the FortiTokenCloud API client as there is an unknown configuration value for the FortiTokenCloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown FortiTokenCloud API Username",
			"The provider cannot create the FortiTokenCloud API client as there is an unknown configuration value for the FortiTokenCloud API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown FortiTokenCloud API Password",
			"The provider cannot create the FortiTokenCloud API client as there is an unknown configuration value for the FortiTokenCloud API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("FTC_HOST")
	clientid := os.Getenv("FTC_CLIENTID")
	clientsecret := os.Getenv("FTC_CLIENTSECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ClientId.IsNull() {
		clientid = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientsecret = config.ClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing FortiTokenCloud API Host",
			"The provider cannot create the FortiTokenCloud API client as there is a missing or empty value for the FortiTokenCloud API host. "+
				"Set the host value in the configuration."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientid == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing FortiTokenCloud API Username",
			"The provider cannot create the FortiTokenCloud API client as there is a missing or empty value for the FortiTokenCloud API username. "+
				"Set the username value in the configuration."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientsecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing FortiTokenCloud API Password",
			"The provider cannot create the FortiTokenCloud API client as there is a missing or empty value for the FortiTokenCloud API password. "+
				"Set the password value in the configuration."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new FortiTokenCloud client using the configuration values
	client, err := ftc_client.NewClient(&host, &clientid, &clientsecret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create FortiTokenCloud API Client",
			"An unexpected error occurred when creating the FortiTokenCloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"FortiTokenCloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the FortiTokenCloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *fortiTokenCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationsDataSource,
		NewRealmDataSource,
	}
}

func (p *fortiTokenCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationResource,
		NewUserSourceResource,
		NewDomainResource,
	}
}
