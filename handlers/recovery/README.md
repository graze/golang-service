# Panic Recovery Handler

```bash
$ go get github.com/graze/golang-service/handlers/recovery
```

The panic Recovery handler recovers from panics and output a nice format to the client, and handles the error using a variety of handlers.
It will always return an InternalServiceError status code (500) and leaves the contents to the user.

You can create custom handlers to do something when a panic occurs:

**Example Handler:**

A simple handler would be to return the error text to the user in the body of the response:

```go
echoHandler := failure.HandlerFunc(func (w io.Writer, r *http.Request, err error, status int) {
    w.Write([]byte(err.Error()))
})
```

**Usage as an http handler:**

The recovery middleware can then be used as an http.Handler as `recovery.New` will return a `func(h http.Handler) http.Handler` allowing it to be chained together with other handlers.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	panic("uh-oh")
})

recoverer = recovery.New(echoHandler)
http.ListenAndServe(":80", recoverer)
```

When calling each handler, recovery will pass http.StatusInternalServerError as the status, but it
won't write the header to allow other handlers to write their own headers. You'll need do something like:

```go
fiveHundredHandler = failure.HandlerFunc(func (w http.ResponseWriter, r *http.Request, err error, status int) {
    w.writeHeader(status)
}))
```

You can of course write a different status should you so choose.

## Logging Panic Handler:

The logging Recoverer will log an output of the recovered panic for debugging.

```go
logPanic := recovery.PanicLogger(log.With(log.KV{"module":"panic.handler"}))
```

## Raygun Panic Handler

To pass panics off to a third party (such as raygun) this handler can be used.

```go
raygunClient, _ := raygun4go.New(name, key)
raygunClient.Silent(false)
raygunClient.Version("1.0")

recoverer := recovery.New(recovery.Raygun(raygunClient))
```

## Combining Multiple Recovery Handlers

You can supply multiple recovery handlers that will each get called when a panic occurs.

```go
recoverer := recovery.New(
    recovery.PanicLogger(log.New()),
    recovery.Raygun(raygunClient),
    echoHandler,
)
```
