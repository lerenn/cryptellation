#!/bin/bash

set -eo pipefail

# Wait for Redis
# TODO

# Wait for NATS
# TODO

# Launch tests
go test -p 1 ./internal/adapters/... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out

# Cleaning up
rm cover.out
