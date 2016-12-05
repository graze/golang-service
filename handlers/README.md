# Handlers

Collection of middleware handlers for use by HTTP services

```bash
$ go get github.com/graze/golang-service/handlers
```

## Context Adder

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

## Healthd Logger

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

## Statsd Logger

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

## Structured Request Logger

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

## Authentication Handler

```bash
$ go get github.com/graze/golang-service/handlers/auth
```

Adds authentication to the request using middleware, with the benefit of linking the authentication with a user

```go
func finder(key string, r *http.Request) (interface{}, error) {
    user, ok := users[key]
    if !ok {
        return nil, fmt.Errorf("No user found for: %s", key)
    }
    return user, nil
}

func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
    w.WriteHeader(status)
    fmt.Fprintf(w, err.Error())
}

keyAuth := auth.ApiKey{
    Provider: "Graze",
    Finder: finder,
    OnError: onError,
}

http.Handle("/", keyAuth.Next(router))
```

You can then retrieve the user within the request handler:

```go
func GetList(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUser(r).(*account.User)
    if !ok {
        w.WriteHeader(403)
        return
    }
}
```
