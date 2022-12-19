FROM golang:alpine

# Install required running dependencies
RUN apk update \
    && apk add bash postgresql-client

# Do not use gcc
ENV CGO_ENABLED 0

# Set entrypoint
WORKDIR /code
