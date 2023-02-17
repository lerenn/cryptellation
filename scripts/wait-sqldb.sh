#!/bin/bash

set -eo pipefail

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
