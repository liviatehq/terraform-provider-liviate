// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package liviate

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LiviateProvider struct{}

type LiviateProviderModel struct {
	ApiUrl      types.String `tfsdk:"api_url"`
	ApiKey      types.String `tfsdk:"api_key"`
	SecretKey   types.String `tfsdk:"secret_key"`
	Config      types.String `tfsdk:"config"`
	Profile     types.String `tfsdk:"profile"`
	HttpGetOnly types.Bool   `tfsdk:"http_get_only"`
	Timeout     types.Int64  `tfsdk:"timeout"`
}

var _ provider.Provider = (*LiviateProvider)(nil)

func New() provider.Provider {
	return &LiviateProvider{}
}

func (p *LiviateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "liviate"
}

func (p *LiviateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"secret_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"config": schema.StringAttribute{
				Optional: true,
			},
			"profile": schema.StringAttribute{
				Optional: true,
			},
			"http_get_only": schema.BoolAttribute{
				Optional: true,
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

func (p *LiviateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiUrl := os.Getenv("LIVIATE_API_URL")
	if apiUrl == "" {
		apiUrl = "https://console.liviate.com/client/api"
	}
	apiKey := os.Getenv("LIVIATE_API_KEY")
	secretKey := os.Getenv("LIVIATE_SECRET_KEY")
	httpGetOnly, _ := strconv.ParseBool(os.Getenv("LIVIATE_HTTP_GET_ONLY"))
	timeout, _ := strconv.ParseInt(os.Getenv("LIVIATE_TIMEOUT"), 2, 64)

	var data LiviateProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.ApiUrl.ValueString() != "" {
		apiUrl = data.ApiUrl.ValueString()
	}

	if data.ApiKey.ValueString() != "" {
		apiKey = data.ApiKey.ValueString()
	}

	if data.SecretKey.ValueString() != "" {
		secretKey = data.SecretKey.ValueString()
	}

	if data.HttpGetOnly.ValueBool() {
		httpGetOnly = true
	}

	if data.Timeout.ValueInt64() != 0 {
		timeout = data.Timeout.ValueInt64()
	}

	cfg := Config{
		APIURL:      apiUrl,
		APIKey:      apiKey,
		SecretKey:   secretKey,
		HTTPGETOnly: httpGetOnly,
		Timeout:     timeout,
	}

	client, err := cfg.NewClient()

	if err != nil {
		resp.Diagnostics.AddError("liviate", fmt.Sprintf("failed to create client: %T", err))
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *LiviateProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.Conflicting(
			path.MatchRoot("api_url"),
			path.MatchRoot("config"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("api_url"),
			path.MatchRoot("profile"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("api_key"),
			path.MatchRoot("config"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("api_key"),
			path.MatchRoot("profile"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("secret_key"),
			path.MatchRoot("config"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("secret_key"),
			path.MatchRoot("profile"),
		),
	}
}

func (p *LiviateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewserviceOfferingUnconstrainedResource,
		NewserviceOfferingConstrainedResource,
		NewserviceOfferingFixedResource,
	}
}

func (p *LiviateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
