// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package recovery

import (
	"io"
	"net/http"
	"runtime/debug"

	"github.com/graze/golang-service/log"
)

// panicLogger is a local struct to implement the Recoverer interface
type panicLogger struct {
	logger log.FieldLogger
}

// Logger takes a panic event and writes a stack trace to the log
func (l panicLogger) Handle(w io.Writer, r *http.Request, err error, status int) {
	l.logger.Ctx(r.Context()).With(log.KV{
		"tag":    "critical_error",
		"stack":  debug.Stack(),
		"status": status,
	}).Err(err).Error("panic occoured")
}

// PanicLogger creates a logs the provided panic that has been recovered
//
// Usage:
//  logger := log.New()
//
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//     panic("oh-o")
//  })
//
//  outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
//      w.Write([]byte("panic happened, oh dear"))
//  }
//  logPanic := recovery.PanicLogger(logger.With(log.KV{"module":"panic.handler"}))
//  recoverer := recovery.New(logPanic)
//  http.ListenAndServe(":80", recoverer.Handle(r))
func PanicLogger(logger log.FieldLogger) Handler {
	return &panicLogger{logger}
}
