FROM golang:alpine

# Code
WORKDIR /app
COPY ./ /app

# Build from code
ENV CGO_ENABLED 0
RUN go build ./cmd/backtests-server

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY build/package/entrypoint.sh .
COPY --from=0 /app/backtests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "backtests-server" ]
