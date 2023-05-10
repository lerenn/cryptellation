# Building image
FROM golang:alpine

# Disable CGO
ENV CGO_ENABLED 0

# Set the workdir
WORKDIR /go/src/github.com/lerenn/cryptellation

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Get all remaining code
COPY ./ .

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build go install ./cmd/*


