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

## Pre-build handlers

This is a collection of pre-made handlers that use environment variables to quickly add to your application

```bash
$ go get github.com/graze/golang-service/handlers
```

### Syslog

Output to syslog using the following environment variables

Environment Variables:
```
    SYSLOG_NETWORK: The network type of the syslog server (tcp, udp) Leave blank for local syslog
    SYSLOG_HOST: The host of the syslog server. Leave blank for local syslog
    SYSLOG_PORT: The port of the syslog server
    SYSLOG_APPLICATION: The application to report the logs as
    SYSLOG_LEVEL: The level to limit messages to (default: LEVEL6)
```
Example:
```
    SYSLOG_NETWORK: udp
    SYSLOG_HOST: app.syslog.local
    SYSLOG_PORT: 1234
    SYSLOG_APPLICATION: app-live
```
Usage:
```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   w.Write([]byte("This is a catch-all route"))
})
loggedRouter := handlers.SyslogHandler(r)
http.ListenAndServe(":1123", loggedRouter)
```

### Healthd

This will create the application.log files in `/var/log/nginx/healthd/` and rotate them for you

```go
loggedRouter := handlers.HealthdHandler(r)
http.ListenAndServe(":1123", loggedRouter)
```

### StatsD

This will output to StatsD using the following variables

Environment Variables:
```
    STATSD_HOST: The host of the statsd server
    STATSD_PORT: The port of the statsd server
    STATSD_NAMESPACE: The namespace to prefix to every metric name
    STATSD_TAGS: A comma separared list of tags to apply to every metric reported
```
Example:
```
    STATSD_HOST: localhost
    STATSD_PORT: 8125
    STATSD_NAMESPACE: app.live.
    STATSD_TAGS: env:live,version:1.0.2
```
Usage:
```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
   w.Write([]byte("This is a catch-all route"))
})
loggedRouter := handlers.StatsdHandler(r)
http.ListenAndServe(":1123", loggedRouter)
```

### Combining

These handlers can be combined in any combination or you can use the `handlers.AllHandlers(r)` as a short cut for all
of them

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
$ make install
$ make test
```

# License

- General code: [MIT License](LICENSE)
- some code: `Copyright (c) 2013 The Gorilla Handlers Authors. All rights reserved.`
