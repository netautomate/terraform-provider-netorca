// Copyright (c) HashiCorp, Inc.

package resouces

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"terraform-provider-netorca/internal/netorca"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// -----------------------------------------------------------------------------
// Interface Assertions
// -----------------------------------------------------------------------------

var (
	_ resource.Resource                = (*changeInstancesResource)(nil)
	_ resource.ResourceWithImportState = (*changeInstancesResource)(nil)
	_ resource.ResourceWithConfigure   = (*changeInstancesResource)(nil)
)

// -----------------------------------------------------------------------------
// Constructor and Type Definitions
// -----------------------------------------------------------------------------

// NewChangeInstanceResource returns a new instance of the changeInstancesResource.
func NewChangeInstanceResource() resource.Resource {
	return &changeInstancesResource{}
}

// changeInstancesResource implements the resource.Resource interface.
type changeInstancesResource struct {
	client *netorca.NetOrcaClient
}

// changeInstanceResourceModel defines the schema model for the resource.
type changeInstanceResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	POV          types.String `tfsdk:"pov"`
	State        types.String `tfsdk:"state"`
	DeployedItem types.String `tfsdk:"deployed_item"`
}

// -----------------------------------------------------------------------------
// Resource Interface Methods
// -----------------------------------------------------------------------------

// Metadata sets the resource type name.
func (d *changeInstancesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_change_instances"
}

// Schema defines the schema for the resource.
func (d *changeInstancesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages NetOrca change instances.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "The NetOrca change instance ID. Structured as {pov}/{change_instance_id}",
			},
			"pov": schema.StringAttribute{
				Required:    true,
				Description: "The NetOrca Point Of View (pov) of the change instance",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "Sets the current state of a change instance e.g. APPROVED|ERROR|COMPLETED",
			},
			"deployed_item": schema.StringAttribute{
				Required:    true,
				Description: "An arbitrary json blob used to attach metadata to change instances.",
			},
		},
	}
}

// Configure sets the provider client on the resource.
func (c *changeInstancesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netorca.NetOrcaClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *NetOrca.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	c.client = client
}

// Create updates a change instance and then refreshes the state.
func (c *changeInstancesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan changeInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := netorca.ChangeInstanceUpdateRequest{
		State:        plan.State.ValueString(),
		Description:  "Updated via terraform",
		DeployedItem: plan.DeployedItem.ValueString(),
	}

	err := c.client.ChangeInstancePatch(plan.ID.ValueInt64(), plan.POV.ValueString(), content)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating change instance id: %d", plan.ID.ValueInt64()), err.Error())
		return
	}

	changeInstance, err := c.client.ChangeInstanceGetById(plan.ID.ValueInt64(), plan.POV.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting change instance id: %d", plan.ID.ValueInt64()), err.Error())
		return
	}

	plan.ID = types.Int64Value(plan.ID.ValueInt64())
	plan.POV = types.StringValue("serviceowner")
	plan.State = types.StringValue(changeInstance.State)

	jsonBytes, err := json.Marshal(changeInstance.ServiceItemField.DeployedItem)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error marshalling deployed_item from service_item: %d", plan.ID.ValueInt64()), err.Error())
		return
	}
	plan.DeployedItem = types.StringValue(string(jsonBytes))

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read retrieves the current state of the resource.
func (c *changeInstancesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state changeInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("ID is: %d", state.ID.ValueInt64()))

	changeInstance, err := c.client.ChangeInstanceGetById(state.ID.ValueInt64(), state.POV.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting change instance id: %d", state.ID.ValueInt64()), err.Error())
		return
	}

	state.ID = types.Int64Value(state.ID.ValueInt64())
	state.POV = types.StringValue(state.POV.ValueString())
	state.State = types.StringValue(changeInstance.State)

	deployedItemData, err := json.Marshal(changeInstance.ServiceItemField.DeployedItem)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error marshalling service_item.deployed_item.data from change instance id: %d", state.ID.ValueInt64()), err.Error())
		return
	}
	state.DeployedItem = types.StringValue(string(deployedItemData))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies the change instance if there are any changes.
func (c *changeInstancesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state changeInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.State.Equal(state.State) || !plan.DeployedItem.Equal(state.DeployedItem) {
		content := netorca.ChangeInstanceUpdateRequest{
			State:        plan.State.ValueString(),
			Description:  "Updated via terraform",
			DeployedItem: plan.DeployedItem.ValueString(),
		}

		err := c.client.ChangeInstancePatch(plan.ID.ValueInt64(), plan.POV.ValueString(), content)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error updating change instance id: %s", plan.ID.String()), err.Error())
			return
		}
	}

	changeInstance, err := c.client.ChangeInstanceGetById(plan.ID.ValueInt64(), plan.POV.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error retrieving change instance id: %s", plan.ID.String()), err.Error())
		return
	}

	deployedItemData, err := json.Marshal(changeInstance.ServiceItemField.DeployedItem)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error marshalling service_item.deployed_item.data from change instance id: %d", state.ID.ValueInt64()), err.Error())
		return
	}

	state.ID = types.Int64Value(changeInstance.Id)
	state.State = types.StringValue(changeInstance.State)
	state.DeployedItem = types.StringValue(string(deployedItemData))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete is a no-op since change instances cannot be deleted.
func (d *changeInstancesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state changeInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// ImportState handles importing an existing resource into Terraform state.
func (c *changeInstancesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state changeInstanceResourceModel

	idTokens := strings.Split(req.ID, "/")
	pov := idTokens[0]
	id, err := strconv.ParseInt(idTokens[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error parsing NetOrca change instance ID from terraform ID: %s", req.ID), err.Error())
		return
	}

	config, importDiags := changeInstanceImportFramework(id, pov, c.client)
	resp.Diagnostics.Append(importDiags...)
	if importDiags.HasError() {
		return
	}

	deployedItemData, err := json.Marshal(config.ServiceItemField.DeployedItem)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error marshalling service_item.deployed_item.data from change instance id: %d", state.ID.ValueInt64()), err.Error())
		return
	}

	state.ID = types.Int64Value(state.ID.ValueInt64())
	state.State = types.StringValue(config.State)
	state.POV = types.StringValue(pov)
	state.DeployedItem = types.StringValue(string(deployedItemData))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// -----------------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------------

// changeInstanceImportFramework retrieves a NetOrca change instance for import.
func changeInstanceImportFramework(id int64, pov string, client *netorca.NetOrcaClient) (netorca.ChangeInstance, diag.Diagnostics) {
	var diags diag.Diagnostics

	changeInstance, err := client.ChangeInstanceGetById(id, pov)
	if err != nil {
		diags.AddError(fmt.Sprintf("Error retrieving NetOrca change instance id: %d", id), err.Error())
	}

	return changeInstance, diags
}
