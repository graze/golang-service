FROM golang:alpine
RUN apk add --no-cache git

WORKDIR "/go"

RUN go get \
    github.com/DataDog/datadog-go/statsd \
    github.com/gorilla/handlers \
    github.com/stretchr/testify/assert \
    golang.org/x/tools/cmd/godoc
