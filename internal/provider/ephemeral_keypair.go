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

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/mod/sumdb/note"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ ephemeral.EphemeralResource = &KeyPairEphemeralResource{}

// NewKeyPairEphemeralResource creates a new instance of the key pair ephemeral resource.
func NewKeyPairEphemeralResource() ephemeral.EphemeralResource {
	return &KeyPairEphemeralResource{}
}

// KeyPairEphemeralResource defines the ephemeral resource implementation.
type KeyPairEphemeralResource struct{}

// KeyPairEphemeralResourceModel describes the ephemeral resource data model.
type KeyPairEphemeralResourceModel struct {
	Name       types.String `tfsdk:"name"`
	PrivateKey types.String `tfsdk:"private_key"`
	PublicKey  types.String `tfsdk:"public_key"`
}

// Metadata returns the ephemeral resource type name.
func (e *KeyPairEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypair"
}

// Schema defines the schema for the ephemeral resource.
func (e *KeyPairEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a public/private key pair for signing notes using the Go sumdb/note package.",
		Description:         "Generates a public/private key pair for signing notes using the Go sumdb/note package.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name to associate with the key pair. This is embedded in the public key.",
				Description:         "The name to associate with the key pair. This is embedded in the public key.",
				Required:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "The generated private key in note format. This value is sensitive and should be protected.",
				Description:         "The generated private key in note format. This value is sensitive and should be protected.",
				Computed:            true,
				Sensitive:           true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "The generated public key in note format.",
				Description:         "The generated public key in note format.",
				Computed:            true,
			},
		},
	}
}

// Configure prepares the ephemeral resource for use.
func (e *KeyPairEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Provider configuration is not needed for this resource
}

// Open generates a new key pair and returns it in the response.
func (e *KeyPairEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data KeyPairEphemeralResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate the key pair using the note.GenerateKey function
	privateKey, publicKey, err := note.GenerateKey(nil, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate key pair",
			fmt.Sprintf("Could not generate note signing key pair: %s", err),
		)
		return
	}

	// Set the generated keys in the model
	data.PrivateKey = types.StringValue(privateKey)
	data.PublicKey = types.StringValue(publicKey)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

// Close performs any cleanup needed when the ephemeral resource is closed.
func (e *KeyPairEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	// No cleanup needed for this ephemeral resource
}
