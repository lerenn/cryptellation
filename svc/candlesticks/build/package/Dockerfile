# Building image
FROM golang:alpine AS build

# Disable CGO
ENV CGO_ENABLED 0

# Get all remaining code
RUN mkdir -p /go/src/github.com/lerenn/cryptellation
COPY ./ /go/src/github.com/lerenn/cryptellation

# Set the workdir
WORKDIR /go/src/github.com/lerenn/cryptellation/svc/candlesticks

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod go install ./cmd/*

# Get final base image
FROM alpine

# Set environment variables
ENV HEALTH_PORT 9000

# Expose ports
EXPOSE 9000

# Copy files
COPY ./svc/candlesticks/build/package/entrypoint.sh /entrypoint.sh

# Get binary
COPY --from=build /go/bin/* /usr/local/bin

# Set entrypoint
ENTRYPOINT ["sh", "/entrypoint.sh"]