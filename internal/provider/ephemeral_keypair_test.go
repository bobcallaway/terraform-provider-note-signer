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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"golang.org/x/mod/sumdb/note"
)

func TestKeyPairEphemeralResource_Schema(t *testing.T) {
	ctx := context.Background()
	resource := NewKeyPairEphemeralResource()

	schemaReq := ephemeral.SchemaRequest{}
	schemaResp := &ephemeral.SchemaResponse{}

	resource.Schema(ctx, schemaReq, schemaResp)

	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResp.Diagnostics)
	}

	if schemaResp.Schema.Attributes == nil {
		t.Fatal("Schema has no attributes")
	}

	// Verify required attributes exist
	if _, ok := schemaResp.Schema.Attributes["name"]; !ok {
		t.Error("Schema missing 'name' attribute")
	}
	if _, ok := schemaResp.Schema.Attributes["private_key"]; !ok {
		t.Error("Schema missing 'private_key' attribute")
	}
	if _, ok := schemaResp.Schema.Attributes["public_key"]; !ok {
		t.Error("Schema missing 'public_key' attribute")
	}
}

func TestKeyPairEphemeralResource_Metadata(t *testing.T) {
	ctx := context.Background()
	resource := NewKeyPairEphemeralResource()

	req := ephemeral.MetadataRequest{
		ProviderTypeName: "note-signer",
	}
	resp := &ephemeral.MetadataResponse{}

	resource.Metadata(ctx, req, resp)

	expected := "note-signer_keypair"
	if resp.TypeName != expected {
		t.Errorf("TypeName = %v, want %v", resp.TypeName, expected)
	}
}

// TestNoteKeyGeneration tests the underlying note.GenerateKey function
func TestNoteKeyGeneration(t *testing.T) {
	tests := []struct {
		name     string
		keyName  string
	}{
		{
			name:    "simple name",
			keyName: "test-key",
		},
		{
			name:    "production name",
			keyName: "production-signer",
		},
		{
			name:    "short name",
			keyName: "k",
		},
		{
			name:    "name with dash",
			keyName: "my-signing-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, publicKey, err := note.GenerateKey(nil, tt.keyName)
			if err != nil {
				t.Fatalf("GenerateKey() error = %v", err)
			}

			// Verify private key is not empty
			if privateKey == "" {
				t.Error("Private key should not be empty")
			}

			// Verify public key is not empty
			if publicKey == "" {
				t.Error("Public key should not be empty")
			}

			// Verify private key format
			privateKeyPattern := regexp.MustCompile(`^PRIVATE\+KEY\+`)
			if !privateKeyPattern.MatchString(privateKey) {
				t.Errorf("Private key format incorrect: %s", privateKey)
			}

			// Verify public key contains the name
			publicKeyPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(tt.keyName) + `\+`)
			if !publicKeyPattern.MatchString(publicKey) {
				t.Errorf("Public key should start with name %s, got: %s", tt.keyName, publicKey)
			}

			// Verify keys are different
			if privateKey == publicKey {
				t.Error("Private key and public key should be different")
			}

			// Verify we can use the generated keys to sign and verify
			signer, err := note.NewSigner(privateKey)
			if err != nil {
				t.Fatalf("NewSigner() error = %v", err)
			}

			verifiers, err := note.NewVerifier(publicKey)
			if err != nil {
				t.Fatalf("NewVerifier() error = %v", err)
			}

			// Create a note and sign it
			text := "test message\n"
			signed, err := note.Sign(&note.Note{Text: text}, signer)
			if err != nil {
				t.Fatalf("Sign() error = %v", err)
			}

			// Verify the signature
			_, err = note.Open(signed, note.VerifierList(verifiers))
			if err != nil {
				t.Fatalf("Open() error = %v, signed = %s", err, signed)
			}
		})
	}
}
