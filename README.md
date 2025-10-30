# Terraform Provider: note-signer

This Terraform provider generates cryptographic key pairs for signing notes using the Go `sumdb/note` package. The keys are compatible with the Go module checksum database format.

## Features

- **Ephemeral Resource**: Generates public/private key pairs that exist only during the Terraform execution
- **Note Signing**: Uses the official Go `sumdb/note` package for key generation
- **Secure by Default**: Private keys are marked as sensitive in Terraform

## Usage

```terraform
terraform {
  required_providers {
    note-signer = {
      source = "bobcallaway/note-signer"
    }
  }
}

provider "note-signer" {}

ephemeral "note-signer_keypair" "example" {
  name = "my-signing-key"
}

output "public_key" {
  value = ephemeral.note-signer_keypair.example.public_key
}

output "private_key" {
  value     = ephemeral.note-signer_keypair.example.private_key
  sensitive = true
}
```

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

### Building The Provider

```shell
go build -o terraform-provider-note-signer
```

### Running Tests

```shell
terraform test
```

### Generating Documentation

Documentation is generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```shell
go generate
```

**Note**: After running `go generate`, commit any generated or modified files (especially in the `docs/` directory). CI will fail if there are uncommitted changes after code generation.

## License

This provider is available under the Apache 2.0 License.
