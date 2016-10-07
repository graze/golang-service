# Golang Service Helpers

[![Build Status](https://travis-ci.org/graze/golang-service.svg?branch=master)](https://travis-ci.org/graze/golang-service)

## Logging

Collection of Logging helpers for use by HTTP services

```bash
$ go get github.com/graze/golang-service/logging
```

### Healthd Logger

- Support the healthd logs from AWS Elastic Beanstalk logs: [AWS](http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)

Usage:
```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   w.Write([]byte("This is a catch-all route"))
})
c, err := statsd.New("127.0.0.1:8125")
loggedRouter := logging.StatsdHandler(c, r)
http.ListenAndServe(":1123", loggedRouter)
```

### Statsd Logger

- Output `response_time` and `count` statistics for each request to a statsd host

Usage:
```go
func statsdHandler(h http.Handler) http.Handler {
    client, err := statsd.New("<ip>:<port>")
    if err != nil {
        panic(err)
    }
    return logging.StatsdHandler(client, h)
}
```

## NetTest

Network helpers when for testing against networks

```bash
$ go get github.com/graze/golang-service/nettest
```

```go
done := make(chan string)
addr, sock, srvWg := nettest.CreateServer(t, "tcp", "localhost:0", done)
defer srvWg.Wait()
defer os.Remove(addr.String())
defer sock.Close()

s, err := net.Dial("tcp", addr.String())
fmt.Fprintf(s, msg + "\n")
if msg = "\n" != <-done {
    panic("message mismatch")
}
```

# Development

## Testing
To run tests, run this on your host machine:

```
$ make build
$ make test
```

# License

- General code: [MIT License](LICENSE)
- some code: `Copyright (c) 2013 The Gorilla Handlers Authors. All rights reserved.`
