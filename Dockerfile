FROM golang:alpine
RUN apk add --no-cache git curl

WORKDIR "/go"

RUN curl https://glide.sh/get | sh
RUN go get golang.org/x/tools/cmd/godoc
