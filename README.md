# Golang Service Helpers

[![Build Status](https://travis-ci.org/graze/golang-service.svg?branch=master)](https://travis-ci.org/graze/golang-service)

- [Log](#log) Structured Contextual logging
- [Handlers](#handlers) http request middleware to add logging (healthd, context, statsd, structured logs)
- [Metrics](#metrics) send monitoring metrics to collectors (currently: stats)
- [NetTest](#nettest) helpers for use when testing networks

## Log

Handle global logging with context. Based on [logrus](https://github.com/Sirupsen/logrus)
with an option to create a global context.

It uses [logfmt](https://brandur.org/logfmt) by default but can also output `json` using a `&logrus.JSONFormatter()`

```bash
$ go get github.com/graze/golang-service/log
```

### Set global properties

Setting these will mean any use of the global logging context or log.New() will use these properties

```go
log.SetFormatter(&logrus.TextFormatter()) // default
log.SetOutput(os.Stderr) // default
log.SetLevel(log.InfoLevel) // default
log.AddFields(log.KV{"service":"super_service"}) // apply `service=super_service` to each log message
```

### logging using the global logger

```go
log.With(log.KV{
    "module": "request_handler",
    "tag":    "received_request"
    "method": "GET",
    "path":   "/path"
}).Info("Received request");
```

Example:
```
time="2016-10-28T10:51:32Z" level=info msg="Received request" module="request_handler" tag="received_request" method=GET path=/path service="super_service"
```

### Log using a local field store

```go
logger := log.New()
logger.Add(log.KV{
    "module": "request_handler"
})
logger.With(log.KV{
    "tag":    "received_request",
    "method": "GET",
    "path":   "/path"
}).Info("Received GET /path")
logger.Err(err).Error("Failed to handle input request")
```

```
time="2016-10-28T10:51:32Z" level=info msg="Received GET /path" tag="received_request" method=GET path=/path module="request_handler"
```

### Log using a context

The logger can use golang's context to pass around fields

```go
logger := log.New()
logger.Add(log.KV{"module": "request_handler"})
context := logger.NewContext(context.Background())
log.Ctx(context).
    With(log.KV{"tag": "received_request"}).
    Info("Received request")
```

```
time="2016-10-28T10:51:32Z" level=info msg="Received request" tag="received_request" module="request_handler"
```

The context can be applied to another local logger

```go
logger := log.New()
logger.Add(log.KV{"module": "request_handler"})
context := logger.NewContext(context.Background())

logger2 := log.New()
logger2.SetOutput(os.Stderr)
logger2.Add(log.KV{"tag": "received_request"})
logger2.Ctx(context).Info("Received request")
```

```
time="2016-10-28T10:51:32Z" level=info msg="Received request" tag="received_request" module="request_handler"
```

### Modifying a loggers properties

```go
logger := log.New()
logger.SetFormatter(&logrus.JSONFormatter{})
logger.SetLevel(log.DebugLevel)
logger.SetOutput(os.Stdout)

logger.Debug("some debug output printed")
```

`logger` implements the `log.Logger` interface which includes `SetFormatter`, `SetLevel`, `SetOutput`, `Level` and `AddHook`

```
{"time":"2016-10-28T10:51:32Z","level":"debug","msg":"some debug output printed"}
```

## Handlers

Collection of middleware handlers for use by HTTP services

```bash
$ go get github.com/graze/golang-service/handlers
```

### Context Adder

`log` a logging context is stored within the request context.

```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    log.Ctx(ctx).With(log.KV{"module":"get"}).Info("logging GET")
}

http.ListenAndServe(":1234", handlers.LogContextHandler(r))
```

Output:
```
time="2016-10-28T10:51:32Z" level=info msg="Logging GET" dur=0.00881 http.host="localhost:1234" http.method=GET http.path="/" http.protocol="HTTP/1.1" http.uri="/" module=get transaction=8ba382cc-5c42-441c-8f48-11029d806b9a
```

### Healthd Logger

- Support the healthd logs from AWS Elastic Beanstalk logs: [AWS](http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)

By default it writes entries to the location: `/var/log/nginx/healthd/application.log.<year>-<month>-<day>-<hour>`. Using the `handlers.HealthdIoHandler(io.Writer, http.Handler)` will write to a custom path

see http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html

Example:

```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a catch-all route"))
})
loggedRouter := handlers.HealthdHandler(r)
http.ListenAndServe(":1123", loggedRouter)
```

### Statsd Logger

- Output `response_time` and `count` statistics for each request to a statsd host

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

To use a manually created statsd client:

```go
c, _ := statsd.New("127.0.0.1:8125")
loggedRouter := handlers.StatsdHandler(c, r)
```

### Structured Request Logger

This outputs a structured log entry for each request send to the http server

```go
r := mux.NewRouter()
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a catch-all route"))
})
loggedRouter := handlers.StructuredLogHandler(
    log.With(log.KV{"module":"request.handler"}),
    r)
http.ListenAndServe(":1123", loggedRouter)
```

Default Output:
```
time="2016-10-28T10:51:32Z" level=info msg="GET / HTTP/1.1" dur=0.003200881 http.bytes=80 http.host="localhost:1123" http.method=GET http.path="/" http.protocol="HTTP/1.1" http.ref= http.status=200 http.uri="/" http.user= module=request.handler tag="request_handled" ts="2016-10-28T10:51:31.542424381Z"
```

## Metrics

Provides statsd metrics sending

Manually create a client
```go
client, _ := metrics.GetStatsd(StatdClientConf{host,port,namespace,tags})
client.Incr("metric", []string{}, 1)
```

Create a client using environment variables. [see above](#statsd-logger)
```go
client, _ := metrics.GetStatsdFromEnv()
client.Incr("metric", []string{}, 1)
```

## NetTest

Network helpers when for testing against networks

```bash
$ go get github.com/graze/golang-service/nettest
```

```go
done := make(chan string)
addr, sock, srvWg := nettest.CreateServer(t, "tcp", ":0", done)
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
