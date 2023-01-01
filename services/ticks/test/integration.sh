#!/bin/bash

# Wait for Redis
# TODO

# Wait for SqlDB
./scripts/wait-sqldb.sh

# Launch tests
go test -p 1 ./internal/infrastructure/... -coverprofile cover.out

# Displaying result
go tool cover -func cover.out

# Cleaning up
rm cover.out
