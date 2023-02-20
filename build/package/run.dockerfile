ARG BUILD_IMAGE=digital-feather/cryptellation:build

# Get build image
FROM ${BUILD_IMAGE} AS BUILD

# Get final base image
FROM alpine

# Get binary
COPY --from=BUILD /go/bin/* /usr/local/bin
