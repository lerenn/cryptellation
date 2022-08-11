#!/bin/bash

set -eo pipefail

# Wait for CockroachDB
./scripts/wait-cockroachdb.sh

# Launch tests
go test ./... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out
