package handlers

import (
    "github.com/graze/golang-service/logging"
    "net/http"
    "time"
    "os"
)

type healthdFileHandler struct {
    path      string
    base      http.Handler
    handler   http.Handler
    timestamp string
}

// ServerHTTP for the healthdFileHandler automatically rotates the log files based on the hour
func (h healthdFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    now := time.Now().UTC().Format("2006-01-02-15")
    if h.timestamp != now || h.base == nil {
        h.timestamp = now
        file := h.path + "application.log." + now
        logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            panic(err)
        }
        defer logFile.Close()
        h.base = logging.HealthdHandler(logFile, h.handler)
    }
    h.base.ServeHTTP(w, req)
}

// healthdHandler returns a logging.HealthdHandler and outputs healthd formatted log files to the appropriate path
func HealthdHandler(h http.Handler) http.Handler {
    path := "/var/log/nginx/healthd/"
    err := os.MkdirAll(path, 0666)
    if err != nil {
        panic(err)
    }

    return healthdFileHandler{path, nil, h, ""}
}
