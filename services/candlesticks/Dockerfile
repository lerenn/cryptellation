FROM golang:alpine

WORKDIR /app

# Code
COPY /services/ ./services

# Build from code
WORKDIR /app/services/candlesticks
RUN --mount=type=cache,mode=0755,target=/root/.cache/go-build \
    --mount=type=cache,mode=0755,target=/root/go \
    go build ./cmd/candlesticks-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY services/candlesticks/entrypoint.sh .
COPY services/candlesticks/scripts /scripts/
COPY --from=0 /app/services/candlesticks/candlesticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "candlesticks-server" ]
