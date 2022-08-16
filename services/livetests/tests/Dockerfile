FROM golang:alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Do not use gcc
ENV CGO_ENABLED 0

# Set entrypoint
WORKDIR /code
ENTRYPOINT [ "/bin/bash", "./tests/entrypoint.sh" ]
