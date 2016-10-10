package handlers

import (
    "net/http"
)

// LoggingHandlers returns rsyslog, statsd and healthd chained handlers for use with AWS and Graze services
func AllHandlers(h http.Handler) http.Handler {
    return HealthdHandler(StatsdHandler(SyslogHandler(h)))
}
