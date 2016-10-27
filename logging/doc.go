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
Package logging provides a collection of logging helpers that use environment variable to set themselves up

Logging Context

Handle global logging with context. Based on [logrus](https://github.com/Sirupsen/logrus)
with an option to create a global context

    package (
        log "github.com/graze/golang-service/logging"
        "github.com/Sirupsen/logrus"
    )

    // set global properties
    log.SetFormatter(&logrus.TextFormatter()) // default
    log.SetOutput(os.Stderr) // default
    log.SetLevel(logrus.InfoLevel) // default
    log.SetFields(log.F{"service":"super_service"}) // apply `service=super_service` to each log message

    // log using global logger
    log.With(log.F{
        "module": "request_handler",
        "tag":    "received_request"
        "method": "GET",
        "path":   "/path"
    }).Info("Received request");

    // log using context logger
    log_context := log.With(log.F{
        "module": "request_handler"
    })
    log_context.With(log.F{
        "tag":    "received_request",
        "method": "GET",
        "path":   "/path"
    }).Info("Recieved GET /path")
    log_context.WithError(err).Error("Failed to handle input request")

Statsd

Connect to a statsd endpoint using environment variables

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
    logger := logging.GetSysLogFromEnv()
    logger.Write("some message")
*/
package logging
