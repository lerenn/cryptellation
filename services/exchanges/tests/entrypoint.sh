#!/bin/bash

set -eo pipefail

# Wait for SqlDB
./scripts/wait-sqldb.sh

# Launch tests
go test -p 1 ./... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out
