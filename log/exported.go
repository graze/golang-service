// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package log

import (
	"context"
	"io"

	"github.com/Sirupsen/logrus"
)

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// key is a type to ensure unique key for context
type key int

// LogKey is the key used for context
const logKey key = 0

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logEntry.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logEntry.SetFormatter(formatter)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	logEntry.SetLevel(level)
}

// Level returns the standard logger level.
func Level() logrus.Level {
	return logEntry.Level()
}

// AddHook adds a new hook to the global logging context
func AddHook(hook logrus.Hook) {
	logEntry.AddHook(hook)
}

// With returns a new LoggerEntry with the supplied fields
func With(fields KV) *LoggerEntry {
	return logEntry.With(fields)
}

// Err creates a new LoggerEntry from the standard logger and adds an error
// to it, using the value defined in ErrorKey as key.
func Err(err error) *LoggerEntry {
	return logEntry.Err(err)
}

// Fields will return the current set of fields in the global context
func Fields() KV {
	return logEntry.Fields()
}

// Ctx will use the provided context with its logs if applicable
func Ctx(ctx context.Context) *LoggerEntry {
	return logEntry.Ctx(ctx)
}

// NewContext adds the current `logEntry` into `ctx`
func NewContext(ctx context.Context) context.Context {
	return logEntry.NewContext(ctx)
}

// AppendContext creates a new context.Context from the supplied ctx with the fields appended to the end
func AppendContext(ctx context.Context, fields KV) context.Context {
	return logContext.AppendContext(ctx, fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logEntry.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logEntry.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logEntry.Info(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logEntry.Error(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logEntry.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logEntry.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logEntry.Infof(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logEntry.Errorf(format, args...)
}
