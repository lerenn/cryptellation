# Dockerfile arguments
ARG BUILDPLATFORM=linux/amd64

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ARG BUILDBASEIMAGE=golang:alpine
ARG TARGETBASEIMAGE=alpine:latest

# Building image
FROM --platform=${BUILDPLATFORM:-linux/amd64} ${BUILDBASEIMAGE} AS build

# Disable CGO
ENV CGO_ENABLED=0

# Get all remaining code
RUN mkdir -p /go/src/github.com/lerenn/cryptellation
COPY ./ /go/src/github.com/lerenn/cryptellation

# Set the workdir
WORKDIR /go/src/github.com/lerenn/cryptellation

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go install ./cmd/worker

# Get final base image
FROM --platform=${TARGETOS}/${TARGETARCH} ${TARGETBASEIMAGE}

# Set environment variables
ENV HEALTH_PORT=9000

# Expose ports
EXPOSE 9000

# Get binary
COPY --from=build /go/bin/* /usr/local/bin

# Set command
CMD ["worker", "serve"]