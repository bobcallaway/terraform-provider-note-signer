// Copyright 2025 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package provider implements the note-signer Terraform provider.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure NoteSignerProvider satisfies various provider interfaces.
var _ provider.Provider = &NoteSignerProvider{}
var _ provider.ProviderWithEphemeralResources = &NoteSignerProvider{}

// NoteSignerProvider defines the provider implementation.
type NoteSignerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NoteSignerProviderModel describes the provider data model.
type NoteSignerProviderModel struct {
}

// Metadata returns the provider type name and version.
func (p *NoteSignerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "note-signer"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *NoteSignerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for generating note signing key pairs using the Go note package.",
	}
}

// Configure prepares the provider for use by reading configuration data.
func (p *NoteSignerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NoteSignerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Resources returns the resources implemented by this provider.
func (p *NoteSignerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns the data sources implemented by this provider.
func (p *NoteSignerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// EphemeralResources returns the ephemeral resources implemented by this provider.
func (p *NoteSignerProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewKeyPairEphemeralResource,
	}
}

// New returns a function that creates a new instance of the provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NoteSignerProvider{
			version: version,
		}
	}
}
