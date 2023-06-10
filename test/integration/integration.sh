#!/bin/bash

# Wait for SqlDB
if [ ! -z ${SQLDB_HOST+x} ]; then
  export PGHOST=$SQLDB_HOST
  export PGPASSWORD=$SQLDB_PASSWORD
  export PGPORT=$SQLDB_PORT
  export PGUSER=$SQLDB_USER
  export PGDATABASE=defaultdb

  until psql -c '\q' &> /dev/null; do
    echo "SqlDB is unavailable - sleeping"
    sleep 3
  done
fi

# Wait for Redis
# TODO

# Wait for NATS
# TODO

# Launch tests
go test -p 1 $(go list ./internal/infra/...)  -coverprofile cover.out

# Displaying result
go tool cover -func cover.out

# Cleaning up
rm cover.out
