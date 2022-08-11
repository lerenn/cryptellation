FROM golang:alpine

WORKDIR /app

# Code
COPY /services/ ./services

# Build from code
WORKDIR /app/services/backtests
RUN --mount=type=cache,mode=0755,target=/root/.cache/go-build \
    --mount=type=cache,mode=0755,target=/root/go \
    go build ./cmd/backtests-server

ENV CGO_ENABLED 0

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY services/backtests/entrypoint.sh .
COPY --from=0 /app/services/backtests/backtests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "backtests-server" ]
