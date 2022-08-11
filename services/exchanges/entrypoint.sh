#!/bin/bash

set -eo pipefail

# Wait for CockroachDB
/scripts/wait-cockroachdb.sh

$@
