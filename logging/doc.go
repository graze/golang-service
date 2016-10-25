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
