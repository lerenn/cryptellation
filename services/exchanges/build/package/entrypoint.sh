#!/bin/bash

set -eo pipefail

# Wait for SqlDB
/scripts/wait-sqldb.sh

$@
