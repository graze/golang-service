# Log

Handle global logging with context. Based on [logrus](https://github.com/Sirupsen/logrus) incorporating golang's `context.context`

It uses [logfmt](https://brandur.org/logfmt) by default but can also output `json` using a `&logrus.JSONFormatter()`

```bash
$ go get github.com/graze/golang-service/log
```

## Set global properties

Setting these will mean any use of the global logging context or log.New() will use these properties

```go
log.SetFormatter(&logrus.TextFormatter()) // default
log.SetOutput(os.Stderr) // default
log.SetLevel(log.InfoLevel) // default
log.AddFields(log.KV{"service":"super_service"}) // apply `service=super_service` to each log message
```

## logging using the global logger

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

## Log using a local field store

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

## Log using a context

The logger can use golang's context to pass around fields

```go
logger := log.New()
logger.Add(log.KV{"module": "request_handler"})
ctx := logger.NewContext(context.Background())
log.Ctx(ctx).
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
ctx := logger.NewContext(context.Background())

logger2 := log.New()
logger2.SetOutput(os.Stderr)
logger2.Add(log.KV{"tag": "received_request"})
logger2.Ctx(ctx).Info("Received request")
```

```
time="2016-10-28T10:51:32Z" level=info msg="Received request" tag="received_request" module="request_handler"
```

## Modifying a loggers properties

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
