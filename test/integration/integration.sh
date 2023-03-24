#!/bin/bash

# Wait for SqlDB
./scripts/wait-sqldb.sh

# Wait for Redis
# TODO

# Wait for NATS
# TODO

# Launch tests
go test -p 1 $(go list ./services/... | grep -e /io/)  -coverprofile cover.out

# Displaying result
go tool cover -func cover.out

# Cleaning up
rm cover.out
