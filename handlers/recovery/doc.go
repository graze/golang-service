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
Package recovery is a http.Handler for handing panics and passing the error to multiple Recoverer handlers

The panic Recovery handler recovers from panics and output a nice format to the client, and handles the error using a variety of handlers.
It will always return an InternalServiceError status code (500) and leaves the contents to the user.

You can create custom handlers to do something when a panic occurs:

Example Handler:

    echoHandler := recovery.HandlerFunc(func (w io.Writer, r *http.Request, err error, status int) {
        w.Write([]byte(err.Error()))
    })

recovery provides an http.Handler for use with http middleware

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        panic("uh-oh")
    })

    recoverer = recovery.New(echoHandler)
    http.ListenAndServe(":80", recoverer)

Logging Panic Handler

The logging Recoverer will log an output of the recovered panic for debugging.

    logPanic := recovery.PanicLogger(log.With(log.KV{"module":"panic.handler"}))

Raygun Panic Handler

To pass panics off to a third party (such as raygun) this handler can be used.

    raygunClient, _ := raygun4go.New(name, key)
    raygunClient.Silent(false)
    raygunClient.Version("1.0")

    recoverer := recovery.New(recovery.Raygun(raygunClient))

Combining Multiple Recovery Handlers

You can supply multiple recovery handlers that will each get called when a panic occurs.

    recoverer := recovery.New(
        recovery.PanicLogger(log.New()),
        recovery.Raygun(raygunClient),
        echoHandler,
    )
*/
package recovery
