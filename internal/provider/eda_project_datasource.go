package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &EdaProjectDataSource{}

func NewEdaProjectDataSource() datasource.DataSource {
	return &EdaProjectDataSource{}
}

type EdaProjectDataSource struct {
	client *providerClient
}

type EdaProjectDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	URL            types.String `tfsdk:"url"`
	SCMBranch      types.String `tfsdk:"scm_branch"`
	OrganizationID types.Int64  `tfsdk:"organization_id"`
}

func (d *EdaProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eda_project"
}

func (d *EdaProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get EDA project datasource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "EDA Project ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "EDA Project name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "EDA Project description.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The SCM URL for the EDA project.",
				Computed:    true,
			},
			"scm_branch": schema.StringAttribute{
				Description: "The SCM branch for the EDA project.",
				Computed:    true,
			},
			"organization_id": schema.Int64Attribute{
				Description: "The organization ID for the EDA project.",
				Computed:    true,
			},
		},
	}
}

func (d *EdaProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*providerClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = configureData
}

func (d *EdaProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EdaProjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int.",
			fmt.Sprintf("Unable to convert id: %v. ", data.ID.ValueString()))
		return
	}

	// Get EDA endpoint and build URL
	url := fmt.Sprintf("eda/api/v1/projects/%d/", id)
	body, statusCode, err := d.client.GenericAPIRequest(ctx, http.MethodGet, url, nil, []int{200, 404}, "gateway")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making API http request",
			fmt.Sprintf("Error was: %s.", err.Error()))
		return
	}

	if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var responseData EdaProjectAPIModel

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshal response body into object",
			fmt.Sprintf("Error =  %v.", err.Error()))
		return
	}

	idAsString := strconv.FormatInt(responseData.ID, 10)
	data.ID = types.StringValue(idAsString)
	data.Name = types.StringValue(responseData.Name)
	data.URL = types.StringValue(responseData.URL)
	data.OrganizationID = types.Int64Value(responseData.OrganizationID)

	if responseData.Description != "" {
		data.Description = types.StringValue(responseData.Description)
	}

	if responseData.SCMBranch != "" {
		data.SCMBranch = types.StringValue(responseData.SCMBranch)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
