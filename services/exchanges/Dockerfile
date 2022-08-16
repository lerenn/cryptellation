FROM golang:alpine

# Code
WORKDIR /app
COPY ./ /app

# Build from code
ENV CGO_ENABLED 0
RUN --mount=type=cache,mode=0755,target=/root/.cache/go-build \
    --mount=type=cache,mode=0755,target=/root/go \
    go build ./cmd/exchanges-server

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY scripts/ /scripts/
COPY entrypoint.sh .
COPY --from=0 /app/exchanges-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "exchanges-server" ]
