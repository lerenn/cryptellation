#!/bin/bash
set -e

# Install tool
if ! command -v golangci-lint &> /dev/null; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

for DIRECTORY in services/* ; do
    cd $DIRECTORY
	echo -e "\e[94m[Linting $DIRECTORY]\e[0m"
    golangci-lint run
    cd -
done
