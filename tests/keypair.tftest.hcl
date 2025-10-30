variables {
  test_name = "default-key"
}

# Basic smoke test to verify the provider loads and can plan the ephemeral resource
run "test_provider_loads" {
  command = plan

  variables {
    test_name = "test-key"
  }

  # Just verify that planning succeeds - we can't assert on ephemeral resource values
  # in terraform test due to their lifecycle. See Go unit tests for detailed validation.
}

run "test_different_name" {
  command = plan

  variables {
    test_name = "production-signer"
  }

  # Verify planning succeeds with different names
}
