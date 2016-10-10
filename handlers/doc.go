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

Syslog

Log requests to a syslog server

Environment Variables:
    SYSLOG_NETWORK: The network type of the syslog server (tcp, udp) Leave blank for local syslog
    SYSLOG_HOST: The host of the syslog server. Leave blank for local syslog
    SYSLOG_PORT: The port of the syslog server
    SYSLOG_APPLICATION: The application to report the logs as
    SYSLOG_LEVEL: The level to limit messages to (default: LEVEL6)

Example:
    SYSLOG_NETWORK: udp
    SYSLOG_HOST: app.syslog.local
    SYSLOG_PORT: 1234
    SYSLOG_APPLICATION: app-live

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := handlers.SyslogHandler(r)
    http.ListenAndServe(":1123", loggedRouter)
*/

package handlers
