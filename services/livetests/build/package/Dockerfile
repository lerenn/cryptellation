FROM golang:alpine

# Code
WORKDIR /app
COPY ./ /app

# Build from code
ENV CGO_ENABLED 0
RUN go build ./cmd/livetests-server

FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash

# Install running files
COPY build/package/entrypoint.sh .
COPY --from=0 /app/livetests-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "livetests-server" ]
