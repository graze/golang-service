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

set global logger properties

    log.SetFormatter(&logrus.TextFormatter{})
    log.SetOutput(os.Stderr)
    log.SetLevel(log.InfoLevel)
    log.Add(log.KV{"service":"super_service"}) // apply `service=super_service` to each log message

logging using the global logger

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

Log using context logger

    // create a fresh context using defaults (ignores the global logger properties set above)
    context := log.New()
    context.Add(log.KV{
        "module": "request_handler"
    })
    context.With(log.KV{
        "tag":    "received_request",
        "method": "GET",
        "path":   "/path"
    }).Info("Recieved GET /path")
    context.Err(err).Error("Failed to handle input request")

    // module=request_handler method=GET tag=received_request path="/path" level=error err="some error"
    //   msg="Failed to handler input request"

Modify a new context logger

    context := log.New()
    context.SetFormatter(&logrus.JSONFormatter{})
    context.SetLevel(log.DebugLevel)
    context.SetOutput(os.Stdout)

    context.Debug("some debug output printed")

    // level=debug msg="some debug output printed"

Using context.Context with the logger

    logger := log.New()
    logger.Add(log.KV{"key":"value"})
    context := logger.NewContext(ctx)

    log.Ctx(context).Info("text")

    // key=value level=info msg=text
*/
package log
