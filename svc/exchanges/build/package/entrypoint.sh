#!/bin/sh

# Fail explicitly if needed
set -euxo pipefail

# Execute migrations
data migrations migrate

# Serve api
api serve