FROM golang:alpine

# Install required running dependencies
RUN apk update \
    && apk add bash postgresql-client

# Do not use gcc
ENV CGO_ENABLED 0

# Set testing dir/entrypoint
WORKDIR /code/services/exchanges
ENTRYPOINT [ "/bin/bash", "./tests/entrypoint.sh" ]
