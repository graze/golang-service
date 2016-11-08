// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

/*
Package log provides some helpers for structured contextual logging

Handle global logging with context. Based on [logrus](https://github.com/Sirupsen/logrus)
with an option to create a global context

    package (
        "github.com/graze/golang-service/log"
        "github.com/Sirupsen/logrus"
    )

Global properties that are used whenever `log.<xxx>` is called can be set as such:

    log.SetFormatter(&logrus.TextFormatter{})
    log.SetOutput(os.Stderr)
    log.SetLevel(log.InfoLevel)
    log.Add(log.KV{"service":"super_service"}) // apply `service=super_service` to each log message

You can then log messages using the `log.` commands which will use the above configuration

    log.Add(log.KV{
        "app": "http-service"
    })

    log.With(log.KV{
        "module": "request_handler",
        "tag":    "received_request"
        "method": "GET",
        "path":   "/path"
    }).Info("Received request");

    // app="http-service" module=request_handler tag=received_request method=GET path="/path" level=info
    //   msg="Received request"

It is also possible to create a new logger ignoring the global configuration set above. Calling `log.New` will
create a new instance of a logger which can be passed around and used with other methods

    // create a fresh context using defaults (ignores the global logger properties set above)
    logger := log.New()
    logger.Add(log.KV{
        "module": "request_handler"
    })
    logger.With(log.KV{
        "tag":    "received_request",
        "method": "GET",
        "path":   "/path"
    }).Info("Recieved GET /path")
    logger.Err(err).Error("Failed to handle input request")

    // module=request_handler method=GET tag=received_request path="/path" level=error err="some error"
    //   msg="Failed to handler input request"

When a new logger is created, the format and output can be modified to change the how messages passed to this logger
are logged

    logger := log.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetLevel(log.DebugLevel)
    logger.SetOutput(os.Stdout)

    logger.Debug("some debug output printed")

    // level=debug msg="some debug output printed"

This logger supports golang's `Context`. You can create a new context and use an existing context as such

    logger := log.New()
    logger.Add(log.KV{"key":"value"})
    context := logger.NewContext(context.Background())

    log.Ctx(context).Info("text")

    // key=value level=info msg=text

You can use a logging context stored within a `context.Context` with a second local logger

    logger := log.New()
    logger.Add(log.KV{"key":"value"})
    context := logger.NewContext(context.Background())

    logger2 := log.New()
    logger2.Add(log.KV{"key2":"value2"})
    logger2.Ctx(context).Info("text")

    // key=value key2=value2 level=info msg=text

As the logger is based on logrus you can add Hooks to each logger to send data to multiple outputs.
See: https://github.com/Sirupsen/logrus#hooks
*/
package log
