# Building image
FROM golang:alpine

# Disable CGO
ENV CGO_ENABLED 0

# Set the workdir
WORKDIR /go/src/github.com/digital-feather/cryptellation

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Get all remaining code
COPY ./ .

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build go install ./cmd/*

# Get final base image
FROM alpine

# Build arguments
ARG COMMAND

# Environment variables
ENV COMMAND=${COMMAND}

# Get binary
COPY --from=0 /go/bin/${COMMAND} /usr/local/bin

# Command
CMD ${COMMAND} serve
