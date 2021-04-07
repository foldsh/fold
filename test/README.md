# Integration Tests

This top level test directory contains integration tests. They range from tests
that combine packages within a given system to full on end to end tests that use
the SDKs and CLI tools.

Before running them you need to build Docker images for the gateway and runtime
and have them availble in the local Docker runtime. Tests that use the images
will always work against the 'latest' tag of the local build.
