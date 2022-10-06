#!/bin/bash

set -eo pipefail

# Wait for SQLDB
/scripts/wait-sqldb.sh

$@
