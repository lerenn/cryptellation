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

# Build from code
RUN go build ./cmd/${BINARY}

# Get base image for the final image
FROM alpine

# Get binary
COPY --from=0 /app/${BINARY} /usr/local/bin

# Set command
CMD [ "${BINARY}" ]
