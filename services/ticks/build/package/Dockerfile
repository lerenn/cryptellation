FROM golang:alpine

# Code
WORKDIR /app
COPY ./ /app

# Build from code
ENV CGO_ENABLED 0
RUN go build ./cmd/ticks-server

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY build/package/entrypoint.sh .
COPY --from=0 /app/ticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "ticks-server" ]
