terraform {
  required_providers {
    note-signer = {
      source = "bobcallaway/note-signer"
    }
  }
}

provider "note-signer" {}

# Generate a key pair for signing notes
ephemeral "note-signer_keypair" "example" {
  name = "my-signing-key"
}

# Example: Use the keys with a local_file resource
# Note: Ephemeral resources can only be used within the same configuration
# They cannot be exported as root module outputs
resource "local_file" "public_key" {
  content  = ephemeral.note-signer_keypair.example.public_key
  filename = "${path.module}/public_key.txt"
}

# The private key should be handled with care and stored securely
resource "local_sensitive_file" "private_key" {
  content  = ephemeral.note-signer_keypair.example.private_key
  filename = "${path.module}/private_key.txt"
}
