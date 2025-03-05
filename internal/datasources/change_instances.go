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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// -----------------------------------------------------------------------------
// Type Definitions
// -----------------------------------------------------------------------------

type changeInstanceDataSource struct {
	client *netorca.NetOrcaClient
}

type changeInstanceDataSourceData struct {
	Pov                 types.String `tfsdk:"pov"`
	ChangeInstanceCount types.Int64  `tfsdk:"change_instance_count"`
	ChangeInstances     types.List   `tfsdk:"change_instances"`
	Filters             types.Object `tfsdk:"filters"`

	// internal field to hold the parsed filters.
	filters *changeInstanceDataSourceFiltersData `tfsdk:"-"`
}

type changeInstanceDataSourceFiltersData struct {
	pov                types.String `tfsdk:"pov"`
	ApplicationId      types.Int64  `tfsdk:"application_id"`
	ChangeType         types.String `tfsdk:"change_type"`
	CommitId           types.String `tfsdk:"commit_id"`
	ConsumerTeamId     types.Int64  `tfsdk:"consumer_team_id"`
	Limit              types.Int64  `tfsdk:"limit"`
	Offset             types.Int64  `tfsdk:"offset"`
	Ordering           types.String `tfsdk:"ordering"`
	ServiceId          types.Int64  `tfsdk:"service_id"`
	ServiceItemId      types.Int64  `tfsdk:"service_item_id"`
	ServiceName        types.String `tfsdk:"service_name"`
	ServiceOwnerTeamId types.Int64  `tfsdk:"service_owner_team_id"`
	State              types.String `tfsdk:"state"`
	SubmissionId       types.Int64  `tfsdk:"submission_id"`
}

// -----------------------------------------------------------------------------
// Interface Assertions - Ensure the data source implements the necessary interfaces for Terraform
// -----------------------------------------------------------------------------

var (
	_ datasource.DataSourceWithConfigure = &changeInstanceDataSource{}
)

func NewChangeInstanceDataSource() datasource.DataSource {
	return &changeInstanceDataSource{}
}

// -----------------------------------------------------------------------------
// DataSourceWithConfigure Interface Methods
// -----------------------------------------------------------------------------

// Metadata sets the data source type name.
func (c *changeInstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_change_instances"
}

// Schema defines the schema for the data source.
func (c *changeInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data provider to return a list of change instances.",
		Attributes: map[string]schema.Attribute{
			"pov": schema.StringAttribute{
				Description: "The POV from which to make the request (serviceowner|consumer).",
				Required:    true,
			},
			"change_instance_count": schema.Int64Attribute{
				Description: "The number of change instances that the request has matched.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"change_instances": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"created": schema.StringAttribute{
							Computed: true,
						},
						"modified": schema.StringAttribute{
							Computed: true,
						},
						"owner": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: changeInstanceOwnerAttrType,
						},
						"consumer_team": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: changeInstanceConsumerTeamAttrType,
						},
						"submission": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: changeInstanceSubmissionAttrType,
						},
						"service_item": schema.ObjectAttribute{
							Computed:       true,
							AttributeTypes: serviceItemAttrType,
						},
					},
				},
			},
			"filters": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"application_id": schema.Int64Attribute{
						Optional: true,
					},
					"change_type": schema.StringAttribute{
						Optional: true,
					},
					"commit_id": schema.StringAttribute{
						Optional: true,
					},
					"consumer_team_id": schema.Int64Attribute{
						Optional: true,
					},
					"limit": schema.Int64Attribute{
						Optional: true,
					},
					"offset": schema.Int64Attribute{
						Optional: true,
					},
					"ordering": schema.StringAttribute{
						Optional: true,
					},
					"service_id": schema.Int64Attribute{
						Optional: true,
					},
					"service_item_id": schema.Int64Attribute{
						Optional: true,
					},
					"service_name": schema.StringAttribute{
						Optional: true,
					},
					"service_owner_team_id": schema.Int64Attribute{
						Optional: true,
					},
					"state": schema.StringAttribute{
						Optional: true,
					},
					"submission_id": schema.Int64Attribute{
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider-configured client to the data source.
func (c *changeInstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (c *changeInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data changeInstanceDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract filter configuration if provided.
	if !data.Filters.IsNull() {
		diags = data.extractFilters(ctx)
		resp.Diagnostics.Append(diags...)
	}
	data.filters.pov = data.Pov

	// Build query map based on filters.
	queryMap := make(map[string]interface{})
	if data.filters != nil {
		queryMap["pov"] = data.filters.pov.ValueString()
		queryMap["application_id"] = data.filters.ApplicationId.ValueInt64()
		queryMap["change_type"] = data.filters.ChangeType.ValueString()
		queryMap["commit_id"] = data.filters.CommitId.ValueString()
		queryMap["consumer_team_id"] = data.filters.ConsumerTeamId.ValueInt64()
		queryMap["limit"] = data.filters.Limit.ValueInt64()
		queryMap["offset"] = data.filters.Offset.ValueInt64()
		queryMap["ordering"] = data.filters.Ordering.ValueString()
		queryMap["service_id"] = data.filters.ServiceId.ValueInt64()
		queryMap["service_item_id"] = data.filters.ServiceItemId.ValueInt64()
		queryMap["service_name"] = data.filters.ServiceName.ValueString()
		queryMap["service_owner_team_id"] = data.filters.ServiceOwnerTeamId.ValueInt64()
		queryMap["state"] = data.filters.State.ValueString()
		queryMap["submission_id"] = data.filters.SubmissionId.ValueInt64()
	}

	query, err := netorca.NewChangeInstanceQuery(queryMap)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintln("Error creating change instance query"), err.Error())
		return
	}

	changeInstancesRaw, err := c.client.ChangeInstanceGet(query)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintln("Error getting change instances"), err.Error())
		return
	}

	changeInstances, err := getTerraformChangeInstances(changeInstancesRaw.Results, resp)
	if err != nil {
		return
	}

	data.ChangeInstances = changeInstances
	data.ChangeInstanceCount = types.Int64Value(int64(changeInstancesRaw.Count))
	tflog.Trace(ctx, "Read a data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// -----------------------------------------------------------------------------
// Helper Methods and Functions
// -----------------------------------------------------------------------------

// extractFilters extracts the filter values from the configuration.
func (c *changeInstanceDataSourceData) extractFilters(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics
	c.filters = &changeInstanceDataSourceFiltersData{}
	diags = c.Filters.As(ctx, c.filters, basetypes.ObjectAsOptions{})
	return diags
}

// getTerraformChangeInstances converts netorca change instances into a Terraform list.
func getTerraformChangeInstances(changeInstances []netorca.ChangeInstance, resp *datasource.ReadResponse) (types.List, error) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: changeInstanceAttrTypes}
	elems := []attr.Value{}

	for _, v := range changeInstances {
		ownerObj := map[string]attr.Value{
			"id":   types.Int64Value(v.Owner.Id),
			"name": types.StringValue(v.Owner.Name),
		}
		ownerObjVal, ownerDiags := types.ObjectValue(changeInstanceOwnerAttrType, ownerObj)
		diags.Append(ownerDiags...)

		submissionObj := map[string]attr.Value{
			"id":        types.Int64Value(v.Submission.Id),
			"commit_id": types.StringValue(v.Submission.CommitId),
		}
		submissionObjVal, submissionDiags := types.ObjectValue(changeInstanceSubmissionAttrType, submissionObj)
		diags.Append(submissionDiags...)

		metadata, err := json.Marshal(v.ConsumerTeam.Metadata)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling consumer team metadata"), err.Error())
			return types.ListNull(elemType), err
		}
		consumerTeamObj := map[string]attr.Value{
			"id":       types.Int64Value(v.ConsumerTeam.Id),
			"name":     types.StringValue(v.ConsumerTeam.Name),
			"metadata": types.StringValue(string(metadata)),
		}
		consumerTeamObjVal, consumerTeamDiags := types.ObjectValue(changeInstanceConsumerTeamAttrType, consumerTeamObj)
		diags.Append(consumerTeamDiags...)

		// Service Item related conversions.
		serviceItemApplicationMetadata, err := json.Marshal(v.ServiceItemField.Application.Metadata)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service_item.application.metadata"), err.Error())
			return types.ListNull(elemType), err
		}
		serviceItemApplicationObj := map[string]attr.Value{
			"id":       types.Int64Value(int64(v.ServiceItemField.Application.Id)),
			"name":     types.StringValue(v.ServiceItemField.Application.Name),
			"metadata": types.StringValue(string(serviceItemApplicationMetadata)),
			"owner":    types.Int64Value(int64(v.ServiceItemField.Application.Owner)),
		}
		serviceItemApplicationObjVal, serviceItemApplicationDiags := types.ObjectValue(serviceItemApplicationAttrType, serviceItemApplicationObj)
		diags.Append(serviceItemApplicationDiags...)

		serviceItemDeclaration, err := json.Marshal(v.ServiceItemField.Declaration)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintln("Error Marshalling service item declaration"), err.Error())
			return types.ListNull(elemType), err
		}
		serviceItemOwnerObj := map[string]attr.Value{
			"id":   types.Int64Value(int64(v.ServiceItemField.Service.Owner.Id)),
			"name": types.StringValue(v.ServiceItemField.Service.Owner.Name),
		}
		serviceItemOwnerObjVal, serviceItemOwnerDiags := types.ObjectValue(serviceItemOwnerAttrType, serviceItemOwnerObj)
		diags.Append(serviceItemOwnerDiags...)

		serviceItemServiceObj := map[string]attr.Value{
			"id":          types.Int64Value(int64(v.ServiceItemField.Service.Id)),
			"name":        types.StringValue(v.ServiceItemField.Service.Name),
			"owner":       types.Object(serviceItemOwnerObjVal),
			"healthcheck": types.BoolValue(v.ServiceItemField.Service.HealthCheck),
		}
		serviceItemServiceObjVal, serviceItemServiceDiags := types.ObjectValue(serviceItemServiceAttrType, serviceItemServiceObj)
		diags.Append(serviceItemServiceDiags...)

		deployedItemData, err := json.Marshal(v.ServiceItemField.DeployedItem)
		if err != nil {
			return types.ListNull(elemType), err
		}

		serviceItemObj := map[string]attr.Value{
			"id":            types.Int64Value(int64(v.ServiceItemField.Id)),
			"url":           types.StringValue(v.ServiceItemField.Url),
			"name":          types.StringValue(v.ServiceItemField.Name),
			"created":       types.StringValue(v.ServiceItemField.Created),
			"modified":      types.StringValue(v.ServiceItemField.Modified),
			"runtime_state": types.StringValue(v.ServiceItemField.RuntimeState),
			"change_state":  types.StringValue(v.ServiceItemField.ChangeState),
			"service":       types.Object(serviceItemServiceObjVal),
			"application":   types.Object(serviceItemApplicationObjVal),
			"declaration":   types.StringValue(string(serviceItemDeclaration)),
			"deployed_item": types.StringValue(string(deployedItemData)),
		}
		serviceItemObjVal, serviceItemDiags := types.ObjectValue(serviceItemAttrType, serviceItemObj)
		diags.Append(serviceItemDiags...)

		obj := map[string]attr.Value{
			"id":            types.Int64Value(v.Id),
			"url":           types.StringValue(v.Url),
			"state":         types.StringValue(v.State),
			"created":       types.StringValue(v.Created),
			"modified":      types.StringValue(v.Modified),
			"owner":         types.Object(ownerObjVal),
			"consumer_team": types.Object(consumerTeamObjVal),
			"submission":    types.Object(submissionObjVal),
			"service_item":  types.Object(serviceItemObjVal),
		}
		objVal, d := types.ObjectValue(changeInstanceAttrTypes, obj)
		diags.Append(d...)
		elems = append(elems, objVal)
	}

	listVal, d := types.ListValue(elemType, elems)
	diags.Append(d...)
	resp.Diagnostics.Append(d...)
	return listVal, nil
}

// -----------------------------------------------------------------------------
// Global Attribute Type Definitions
// -----------------------------------------------------------------------------

var changeInstanceAttrTypes = map[string]attr.Type{
	"id":       types.Int64Type,
	"url":      types.StringType,
	"state":    types.StringType,
	"created":  types.StringType,
	"modified": types.StringType,
	"owner": types.ObjectType{
		AttrTypes: changeInstanceOwnerAttrType,
	},
	"consumer_team": types.ObjectType{
		AttrTypes: changeInstanceConsumerTeamAttrType,
	},
	"submission": types.ObjectType{
		AttrTypes: changeInstanceSubmissionAttrType,
	},
	"service_item": types.ObjectType{
		AttrTypes: serviceItemAttrType,
	},
}

var changeInstanceConsumerTeamAttrType = map[string]attr.Type{
	"id":       types.Int64Type,
	"name":     types.StringType,
	"metadata": types.StringType,
}

var changeInstanceOwnerAttrType = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
}

var changeInstanceSubmissionAttrType = map[string]attr.Type{
	"id":        types.Int64Type,
	"commit_id": types.StringType,
}
