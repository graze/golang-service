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
    logger := log.NewLogfmtLogger(os.Stdout)
    loggedRouter := logging.StructuredHandler(log.NewContext(log).With("component", "handler"), r)
    http.ListenAndServe(":1123", loggedRouter)
*/
package handlers
