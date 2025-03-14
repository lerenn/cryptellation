# Dockerfile arguments
ARG BUILDPLATFORM=linux/amd64

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ARG BUILDGOBASEIMAGE=golang:alpine
ARG BUILDJSBASEIMAGE=node:alpine
ARG TARGETBASEIMAGE=alpine:latest

# Building Javascript image
FROM --platform=${BUILDPLATFORM:-linux/amd64} ${BUILDJSBASEIMAGE} AS build_js

# Build the ui
RUN mkdir -p /ui
COPY ./web/ui /ui
WORKDIR /ui
RUN --mount=type=cache,target=/ui/node_modules \
    npm install && npm run build

# Building Golang image
FROM --platform=${BUILDPLATFORM:-linux/amd64} ${BUILDGOBASEIMAGE} AS build_go

# Disable CGO
ENV CGO_ENABLED=0

# Get all remaining code
RUN mkdir -p /go/src/github.com/lerenn/cryptellation
COPY ./ /go/src/github.com/lerenn/cryptellation

# Set the workdir
WORKDIR /go/src/github.com/lerenn/cryptellation

# Copy the ui
COPY --from=build_js /ui/build /go/src/github.com/lerenn/cryptellation/web/ui/build

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go install ./cmd/*

# Get final base image
FROM --platform=${TARGETOS}/${TARGETARCH} ${TARGETBASEIMAGE} AS final

# Set environment variables
ENV HEALTH_PORT=9000

# Expose ports
EXPOSE 9000

# Get binary
COPY --from=build_go /go/bin/* /usr/local/bin

# Set command
CMD ["worker", "serve"]