#!/bin/bash

# Wait for SqlDB
./scripts/wait-sqldb.sh

# Launch tests
go test -p 1 $(go list ./... | grep -e adapters)  -coverprofile cover.out

# Displaying result
go tool cover -func cover.out

# Cleaning up
rm cover.out
