terraform {
  required_providers {
    note-signer = {
      source = "bobcallaway/note-signer"
    }
  }
}

provider "note-signer" {}

variable "test_name" {
  type        = string
  description = "The name to use for the key pair"
}

ephemeral "note-signer_keypair" "test" {
  name = var.test_name
}
