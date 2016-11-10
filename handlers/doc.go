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
Package handlers provides a collection of logging http.Handlers for use by HTTP services that take in configuration
from environment variables

Combining Handlers

We can combine the following handlers automatically or manually

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.AllHandlers(r)
    http.ListenAndServe(":1123", loggedRouter)

They can also be manually chained together
    loggedRouter := handlers.StatsdHandler(handlers.HealthdHandler(r))

Logging Context

This creates a logging context to be passed into the handling function with information about the request

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Ctx(r.Context()).Info("log a message with the context")
        w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.LogContextHandler(r)
    http.ListenAndServe(":1123", loggedRouter)

Healthd

This provides healthd logging (http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)
for health monitoring when using Elastic Beanstalk.

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.HealthdHandler(r)
    http.ListenAndServe(":1123", loggedRouter)

Statsd

Log request duration to a statsd host

Environment Variables:
    STATSD_HOST: The host of the statsd server
    STATSD_PORT: The port of the statsd server
    STATSD_NAMESPACE: The namespace to prefix to every metric name
    STATSD_TAGS: A comma separared list of tags to apply to every metric reported

Example:
    STATSD_HOST: localhost
    STATSD_PORT: 8125
    STATSD_NAMESPACE: app.live.
    STATSD_TAGS: env:live,version:1.0.2

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.StatsdHandler(r)
    http.ListenAndServe(":1123", loggedRouter)

Structured

Log requests using a structured format for handling with json/logfmt

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.StructuredLogHandler(
        log.With(log.KV{"module":"request.handler"}),
        r)
    http.ListenAndServe(":1123", loggedRouter)

Default Output:
    time="2016-10-28T10:51:32Z" level=info msg="GET / HTTP/1.1" dur=0.003200881 http.bytes=80 http.host="localhost:1123" http.method=GET http.path="/" http.protocol="HTTP/1.1" http.ref= http.status=200 http.uri="/" http.user= module=request.handler tag="request_handled" ts="2016-10-28T10:51:31.542424381Z"
*/
package handlers
