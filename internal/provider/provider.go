// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"

	"terraform-provider-netorca/internal/datasources"
	"terraform-provider-netorca/internal/netorca"
	resouces "terraform-provider-netorca/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// -----------------------------------------------------------------------------
// Interface Assertion
// -----------------------------------------------------------------------------

var _ provider.Provider = (*netOrcaProvider)(nil)

// -----------------------------------------------------------------------------
// Type Definitions
// -----------------------------------------------------------------------------

type netOrcaProvider struct{}

type netorcaProviderConfigModel struct {
	Url    types.String `tfsdk:"url"`
	ApiKey types.String `tfsdk:"apikey"`
}

// New returns a function that creates a new instance of netOrcaProvider - implementing the provider.Provider interface. (required by the Terraform)
func New() func() provider.Provider {
	return func() provider.Provider {
		return &netOrcaProvider{}
	}
}

// -----------------------------------------------------------------------------
// Provider Interface Methods
// -----------------------------------------------------------------------------

// Metadata defines metadata for the provider such as name.
func (p *netOrcaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netorca"
}

// Schema defines the provider-level schema for configuration.
func (p *netOrcaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with NetOrca.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "URL for NetOrca API.",
				Optional:    true,
			},
			"apikey": schema.StringAttribute{
				Description: "Api-Key for NetOrca API authentication. ",
				Optional:    true,
			},
		},
	}
}

// Configure configures the provider by creating a NetOrca client using the provided configuration.
func (p *netOrcaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring NetOrca client")

	// Retrieve provider configuration.
	var config netorcaProviderConfigModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown NetOrca API Host",
			"The provider cannot create the NetOrca API client because the NetOrca API url is unknown. "+
				"Either target apply the source of the value first, set the value statically, or use the NETORCA_URL environment variable.",
		)
	}
	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"NetOrca Api-Key is not set",
			"The provider cannot create the NetOrca API client because the NetOrca API apikey is unknown. "+
				"Either target apply the source of the value first, set the value statically, or use the NETORCA_API_KEY environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Use environment variables as defaults if not set via configuration.
	url := os.Getenv("NETORCA_URL")
	apikey := os.Getenv("NETORCA_API_KEY")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}
	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	// Validate required configuration values.
	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing NetOrca API Host",
			"The provider cannot create the NetOrca API client because the NetOrca API host is missing or empty. "+
				"Set the host in the configuration or use the NETORCA_URL environment variable.",
		)
	}
	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing NetOrca Api-Key",
			"The provider cannot create the NetOrca API client because the NetOrca API apikey is missing or empty. "+
				"Set the apikey in the configuration or use the NETORCA_API_KEY environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "netorca_url", url)
	ctx = tflog.SetField(ctx, "netorca_api_key", apikey)
	tflog.Debug(ctx, "Creating NetOrca client")

	// Create a new NetOrca client.
	client := netorca.NewClient(&url, &apikey, ctx)
	// Make the NetOrca client available to DataSources and Resources.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured NetOrca client", map[string]interface{}{"success": true})
}

// DataSources returns the list of data sources provided by the provider.
func (p *netOrcaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewChangeInstanceDataSource,
		datasources.NewServiceItemDataSource,
	}
}

// Resources returns the list of resources provided by the provider.
func (p *netOrcaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resouces.NewChangeInstanceResource,
	}
}
