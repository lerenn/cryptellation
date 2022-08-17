#!/bin/bash

set -eo pipefail

# Wait for Redis
# TODO

# Wait for NATS
# TODO

# Launch tests
go test ./... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out
