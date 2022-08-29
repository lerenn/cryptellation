#!/bin/bash

set -eo pipefail

# Wait for Redis
# TODO

# Launch tests
go test -p 1 ./... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out
