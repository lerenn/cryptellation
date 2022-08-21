FROM golang:alpine

# Code
WORKDIR /app
COPY ./ /app

# Build from code
ENV CGO_ENABLED 0
RUN go build ./cmd/candlesticks-server


FROM alpine

# Install required running dependencies
RUN apk update \
    && apk add bash \
    && apk add postgresql-client

# Install running files
COPY scripts/ /scripts/
COPY build/package/entrypoint.sh .
COPY --from=0 /app/candlesticks-server /usr/local/bin

ENTRYPOINT [ "/bin/bash", "/entrypoint.sh" ]
CMD [ "candlesticks-server" ]
