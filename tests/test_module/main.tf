terraform {
  required_providers {
    note-signer = {
      source = "bobcallaway/note-signer"
    }
  }
}

variable "test_name" {
  type        = string
  description = "The name to use for the key pair"
}

ephemeral "note-signer_keypair" "test" {
  name = var.test_name
}

output "private_key" {
  value     = ephemeral.note-signer_keypair.test.private_key
  sensitive = true
}

output "public_key" {
  value = ephemeral.note-signer_keypair.test.public_key
}
