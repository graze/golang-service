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

Usage
 	r := mux.NewRouter()
 	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
 	   panic("oh-o")
 	})

 	outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
 		w.Write([]byte("panic happened, oh dear"))
 	}
 	recoverer := recovery.New(r, recovery.Logger(log.New()), raygun.New(raygunClient), recovery.RecovererFunc(format))
 	http.ListenAndServe(":80", recoverer)

Logging Panic Handler

The logging Recoverer will log an output of the recovered panic for debugging.

Usage:
    logger := log.New()

    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       panic("oh-o")
    })

    outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
        w.Write([]byte("panic happened, oh dear"))
    }
    logPanic := recovery.PanicLogger(logger.With(log.KV{"module":"panic.handler"}))
    recoverer := recovery.New(r, logPanic)
    http.ListenAndServe(":80", recoverer)

Raygun Panic Handler

To pass panics off to a third party (such as raygun) this handler can be used.

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       panic("oh-o")
    })

    outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
        w.Write([]byte("panic happened, oh dear"))
    }

    raygunClient, _ := raygun4go.New(name, key)
    raygunClient.Silent(false)
    raygunClient.Version("1.0")

    raygunHandler = raygun.New(raygunClient)
    recoverer := recovery.New(r, raygunHandler)
    http.ListenAndServe(":80", recoverer)
*/
package recovery
