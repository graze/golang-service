FROM golang:alpine
RUN apk add --no-cache git curl

WORKDIR "/go"

RUN curl https://glide.sh/get | sh
RUN go get -u \
    golang.org/x/tools/cmd/godoc \
    cmd/gofmt \
    golang.org/x/tools/cmd/goimports \
    github.com/golang/lint/golint
