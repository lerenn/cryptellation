FROM golang:alpine

WORKDIR /app

# Code
COPY /services/ ./services
COPY go.* .

# Build from code
WORKDIR /app/services/exchanges
RUN --mount=type=cache,mode=0755,target=/root/.cache/go-build \
    --mount=type=cache,mode=0755,target=/root/go \
    go build ./cmd/exchanges-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/exchanges/entrypoint.sh .
COPY services/exchanges/scripts /scripts/
COPY --from=0 /app/services/exchanges/exchanges-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "exchanges-server" ]
