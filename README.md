# Golang Service Helpers

[![Build Status](https://travis-ci.org/graze/golang-service.svg?branch=master)](https://travis-ci.org/graze/golang-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/graze/golang-service)](https://goreportcard.com/report/github.com/graze/golang-service)
[![GoDoc](https://godoc.org/github.com/graze/golang-service?status.svg)](https://godoc.org/github.com/graze/golang-service)

- [Log](log/README.md) Structured logging
- [Handlers](handlers/README.md) http request middleware to add logging (auth, healthd, log context, statsd, structured logs)
- [Metrics](metrics/README.md) send monitoring metrics to collectors (currently: stats)
- [NetTest](nettest/README.md) helpers for use when testing networks
- [Validation](validate/README.md) to ensure the user input is correct

[Godoc Documentation](https://godoc.org/github.com/graze/golang-service)

# Development

The development for this repository is done using docker.

## Testing

To run tests, run this on your host machine:

```
$ make install
$ make test
```

# License

- General code: [MIT License](LICENSE)
- some code: `Copyright (c) 2013 The Gorilla Handlers Authors. All rights reserved.`
