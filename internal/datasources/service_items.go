// Copyright (c) HashiCorp, Inc.

package datasources

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-netorca/internal/netorca"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// -----------------------------------------------------------------------------
// Type Definitions
// -----------------------------------------------------------------------------

type serviceItemDataSource struct {
	client *netorca.NetOrcaClient
}

type serviceItemDataSourceData struct {
	ServiceItemCount types.Int64  `tfsdk:"service_item_count"`
	Pov              types.String `tfsdk:"pov"`
	ServiceItems     types.List   `tfsdk:"service_items"`
	Filters          types.Object `tfsdk:"filters"`

	// internal field for parsed filter values.
	filters *serviceItemDataSourceFiltersData `tfsdk:"-"`
}

type serviceItemDataSourceFiltersData struct {
	ApplicationId      types.Int64  `tfsdk:"application_id"`
	ChangeState        types.String `tfsdk:"change_state"`
	ConsumerTeamId     types.Int64  `tfsdk:"consumer_team_id"`
	Limit              types.Int64  `tfsdk:"limit"`
	Name               types.String `tfsdk:"name"`
	Offset             types.Int64  `tfsdk:"offset"`
	Ordering           types.String `tfsdk:"ordering"`
	RuntimeState       types.String `tfsdk:"runtime_state"`
	ServiceOwnerId     types.Int64  `tfsdk:"service_owner_id"`
	ServiceOwnerTeamId types.Int64  `tfsdk:"service_owner_team_id"`
	ServiceName        types.String `tfsdk:"service_name"`
}

// -----------------------------------------------------------------------------
// Interface Assertions - Ensure the data source implements the necessary interfaces for Terraform
// -----------------------------------------------------------------------------

var (
	_ datasource.DataSourceWithConfigure = &serviceItemDataSource{}
)

// NewServiceItemDataSource returns a new instance of serviceItemDataSource.
func NewServiceItemDataSource() datasource.DataSource {
	return &serviceItemDataSource{}
}

// -----------------------------------------------------------------------------
// DataSourceWithConfigure Interface Methods
// -----------------------------------------------------------------------------

// Metadata sets the data source type name.
func (e *serviceItemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_items"
}

// Schema defines the schema for the data source.
func (e *serviceItemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data provider to return a list of change instances.",
		Attributes: map[string]schema.Attribute{
			"pov": schema.StringAttribute{
				MarkdownDescription: "The POV from which to make the request (serviceowner|consumer)",
				Required:            true,
			},
			"service_item_count": schema.Int64Attribute{
				MarkdownDescription: "The number of change instances returned as a part of this query",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"application_id": schema.Int64Attribute{
						MarkdownDescription: "Returns only service items matching specified application_id.",
						Optional:            true,
					},
					"change_state": schema.StringAttribute{
						MarkdownDescription: "Returns only service items matching specified change state. (CHANGES_PENDING|CHANGES_ERRORED|CHANGES_REJECTED|CHANGES_APPROVED)",
						Optional:            true,
					},
					"consumer_team_id": schema.Int64Attribute{
						MarkdownDescription: "Returns only service items matching specified consumer team id.",
						Optional:            true,
					},
					"limit": schema.Int64Attribute{
						MarkdownDescription: "Returns a number of results up to the configured limit number.",
						Optional:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Returns a specific service item with the given name.",
						Optional:            true,
					},
					"offset": schema.Int64Attribute{
						MarkdownDescription: "The initial index from which to return results.",
						Optional:            true,
					},
					"ordering": schema.StringAttribute{
						MarkdownDescription: "The name of the field to use when ordering results.",
						Optional:            true,
					},
					"runtime_state": schema.StringAttribute{
						MarkdownDescription: "The current state of the service.",
						Optional:            true,
					},
					"service_owner_id": schema.Int64Attribute{
						MarkdownDescription: "Returns service items of a given service owner.",
						Optional:            true,
					},
					"service_owner_team_id": schema.Int64Attribute{
						MarkdownDescription: "Returns service items of a given service owner team.",
						Optional:            true,
					},
					"service_name": schema.StringAttribute{
						MarkdownDescription: "Name of a service to filter by.",
						Optional:            true,
					},
				},
			},
			"service_items": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"created": schema.StringAttribute{
							Computed: true,
						},
						"modified": schema.StringAttribute{
							Computed: true,
						},
						"runtime_state": schema.StringAttribute{
							Computed: true,
						},
						"change_state": schema.StringAttribute{
							Computed: true,
						},
						"service": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: serviceItemServiceAttrType,
						},
						"application": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: serviceItemApplicationAttrType,
						},
						"deployed_item": schema.StringAttribute{
							Computed: true,
						},
						"consumer_team": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: serviceItemConsumerTeamAttrTypes,
						},
						"service_owner_team": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: serviceItemServiceOwnerTeamAttrTypes,
						},
						"declaration": schema.StringAttribute{
							Computed: true,
						},
						"healthcheck_status": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (c *serviceItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netorca.NetOrcaClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NetOrca.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	c.client = client
}

// Read is called when Terraform needs to read the state of the data source.
func (c *serviceItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serviceItemDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Filters.IsNull() {
		diags = data.extractFilters(ctx)
		resp.Diagnostics.Append(diags...)
	}

	// Set values for query parameters.
	serviceItemQuery := make(map[string]interface{})
	serviceItemQuery["pov"] = data.Pov.ValueString()
	if data.filters != nil {
		serviceItemQuery["application_id"] = data.filters.ApplicationId.ValueInt64()
		serviceItemQuery["change_state"] = data.filters.ChangeState.ValueString()
		serviceItemQuery["consumer_team_id"] = data.filters.ConsumerTeamId.ValueInt64()
		serviceItemQuery["limit"] = data.filters.Limit.ValueInt64()
		serviceItemQuery["name"] = data.filters.Name.ValueString()
		serviceItemQuery["offset"] = data.filters.Offset.ValueInt64()
		serviceItemQuery["ordering"] = data.filters.Ordering.ValueString()
		serviceItemQuery["runtime_state"] = data.filters.RuntimeState.ValueString()
		serviceItemQuery["service_owner_id"] = data.filters.ServiceOwnerId.ValueInt64()
		serviceItemQuery["service_owner_team_id"] = data.filters.ServiceOwnerTeamId.ValueInt64()
		serviceItemQuery["service_name"] = data.filters.ServiceName.ValueString()
	}

	query, err := netorca.NewServiceItemQuery(serviceItemQuery)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintln("Error with service item definition"), err.Error())
		return
	}

	serviceItems, err := c.client.ServiceItemsGet(query)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintln("Error getting service items"), err.Error())
		return
	}

	data.ServiceItemCount = types.Int64Value(int64(serviceItems.Count))
	data.ServiceItems, err = getTerraformServiceItems(serviceItems.Results, resp)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintln("Error serialising netorca service items into terraform objects"), err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// -----------------------------------------------------------------------------
// Helper Methods and Functions
// -----------------------------------------------------------------------------

// extractFilters extracts filter information from the data source configuration.
func (c *serviceItemDataSourceData) extractFilters(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	c.filters = &serviceItemDataSourceFiltersData{}
	diags = c.Filters.As(ctx, c.filters, basetypes.ObjectAsOptions{})
	return diags
}

// getTerraformServiceItems converts netorca service items into a Terraform list.
func getTerraformServiceItems(serviceItems []netorca.ServiceItem, resp *datasource.ReadResponse) (types.List, error) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: serviceItemAttrTypes}
	elems := []attr.Value{}

	for _, v := range serviceItems {
		applicationMetadata, err := json.Marshal(v.Application.Metadata)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.application.metadata"), err.Error())
			return types.ListNull(elemType), err
		}

		netorcaApplicationObj := map[string]attr.Value{
			"id":       types.Int64Value(int64(v.Application.Id)),
			"name":     types.StringValue(v.Application.Name),
			"metadata": types.StringValue(string(applicationMetadata)),
			"owner":    types.Int64Value(int64(v.Application.Owner)),
		}
		netorcaApplicationObjVal, netorcaApplicationOwnerDiags := types.ObjectValue(serviceItemApplicationAttrType, netorcaApplicationObj)
		diags.Append(netorcaApplicationOwnerDiags...)

		serviceItemOwnerObj := map[string]attr.Value{
			"id":   types.Int64Value(int64(v.Service.Owner.Id)),
			"name": types.StringValue(v.Service.Owner.Name),
		}
		serviceItemOwnerObjVal, serviceItemOwnerDiags := types.ObjectValue(serviceItemOwnerAttrType, serviceItemOwnerObj)
		diags.Append(serviceItemOwnerDiags...)

		serviceItemServiceObj := map[string]attr.Value{
			"id":          types.Int64Value(int64(v.Service.Id)),
			"name":        types.StringValue(v.Service.Name),
			"owner":       types.Object(serviceItemOwnerObjVal),
			"healthcheck": types.BoolValue(v.Service.HealthCheck),
		}
		serviceItemServiceObjVal, serviceItemServiceDiags := types.ObjectValue(serviceItemServiceAttrType, serviceItemServiceObj)
		diags.Append(serviceItemServiceDiags...)

		consumerTeamMetadata, err := json.Marshal(v.ConsumerTeam.Metadata)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.deployed_item.consumer_team.metadata"), err.Error())
			return types.ListNull(elemType), err
		}

		serviceItemConsumerTeamObj := map[string]attr.Value{
			"id":       types.Int64Value(int64(v.ConsumerTeam.Id)),
			"name":     types.StringValue(v.ConsumerTeam.Name),
			"metadata": types.StringValue(string(consumerTeamMetadata)),
		}

		serviceItemServiceOwnerTeamObj := map[string]attr.Value{
			"id":   types.Int64Value(int64(v.ServiceOwnerTeam.Id)),
			"name": types.StringValue(v.ServiceOwnerTeam.Name),
		}

		deployedItemData, err := json.Marshal(v.DeployedItem)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.deployed_item.data"), err.Error())
			return types.ListNull(elemType), err
		}

		consumerTeamMetadata, err = json.Marshal(v.ConsumerTeam.Metadata)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.consumer_team.metadata"), err.Error())
			return types.ListNull(elemType), err
		}
		serviceItemConsumerTeamObj = map[string]attr.Value{
			"id":       types.Int64Value(int64(v.ConsumerTeam.Id)),
			"name":     types.StringValue(v.ConsumerTeam.Name),
			"metadata": types.StringValue(string(consumerTeamMetadata)),
		}
		serviceItemConsumerTeamObjVal, serviceItemConsumerTeamDiags := types.ObjectValue(serviceItemConsumerTeamAttrTypes, serviceItemConsumerTeamObj)
		diags.Append(serviceItemConsumerTeamDiags...)

		serviceItemServiceOwnerTeamObj = map[string]attr.Value{
			"id":   types.Int64Value(int64(v.ServiceOwnerTeam.Id)),
			"name": types.StringValue(v.ServiceOwnerTeam.Name),
		}
		serviceItemServiceOwnerTeamObjVal, serviceItemServiceOwnerTeamDiags := types.ObjectValue(serviceItemServiceOwnerTeamAttrTypes, serviceItemServiceOwnerTeamObj)
		diags.Append(serviceItemServiceOwnerTeamDiags...)

		declarationValues, err := json.Marshal(v.Declaration)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.declaration"), err.Error())
			return types.ListNull(elemType), err
		}

		obj := map[string]attr.Value{
			"id":                 types.Int64Value(v.Id),
			"url":                types.StringValue(v.Url),
			"name":               types.StringValue(v.Name),
			"created":            types.StringValue(v.Created),
			"modified":           types.StringValue(v.Modified),
			"runtime_state":      types.StringValue(v.RuntimeState),
			"change_state":       types.StringValue(v.ChangeState),
			"service":            types.Object(serviceItemServiceObjVal),
			"application":        types.Object(netorcaApplicationObjVal),
			"deployed_item":      types.StringValue(string(deployedItemData)),
			"consumer_team":      types.Object(serviceItemConsumerTeamObjVal),
			"service_owner_team": types.Object(serviceItemServiceOwnerTeamObjVal),
			"declaration":        types.StringValue(string(declarationValues)),
			"healthcheck_status": types.Int64PointerValue(v.HealthcheckStatus),
		}
		objVal, d := types.ObjectValue(serviceItemAttrTypes, obj)
		diags.Append(d...)
		elems = append(elems, objVal)
	}

	listVal, d := types.ListValue(elemType, elems)
	diags.Append(d...)
	resp.Diagnostics.Append(d...)
	return listVal, nil
}

// -----------------------------------------------------------------------------
// Global Variables (Attribute Type Definitions)
// -----------------------------------------------------------------------------

var serviceItemAttrType = map[string]attr.Type{
	"id":            types.Int64Type,
	"url":           types.StringType,
	"name":          types.StringType,
	"created":       types.StringType,
	"modified":      types.StringType,
	"runtime_state": types.StringType,
	"change_state":  types.StringType,
	"service": types.ObjectType{
		AttrTypes: serviceItemServiceAttrType,
	},
	"application": types.ObjectType{
		AttrTypes: serviceItemApplicationAttrType,
	},
	"declaration":   types.StringType,
	"deployed_item": types.StringType,
}

var serviceItemServiceAttrType = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
	"owner": types.ObjectType{
		AttrTypes: serviceItemOwnerAttrType,
	},
	"healthcheck": types.BoolType,
}

var serviceItemOwnerAttrType = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
}

var serviceItemApplicationAttrType = map[string]attr.Type{
	"id":       types.Int64Type,
	"name":     types.StringType,
	"metadata": types.StringType,
	"owner":    types.Int64Type,
}

var serviceItemAttrTypes = map[string]attr.Type{
	"id":            types.Int64Type,
	"url":           types.StringType,
	"name":          types.StringType,
	"created":       types.StringType,
	"modified":      types.StringType,
	"runtime_state": types.StringType,
	"change_state":  types.StringType,
	"service": types.ObjectType{
		AttrTypes: serviceItemServiceAttrType,
	},
	"application": types.ObjectType{
		AttrTypes: serviceItemApplicationAttrType,
	},
	"deployed_item": types.StringType,
	"consumer_team": types.ObjectType{
		AttrTypes: serviceItemConsumerTeamAttrTypes,
	},
	"service_owner_team": types.ObjectType{
		AttrTypes: serviceItemServiceOwnerTeamAttrTypes,
	},
	"declaration":        types.StringType,
	"healthcheck_status": types.Int64Type,
}

var serviceItemConsumerTeamAttrTypes = map[string]attr.Type{
	"id":       types.Int64Type,
	"name":     types.StringType,
	"metadata": types.StringType,
}

var serviceItemServiceOwnerTeamAttrTypes = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
}
