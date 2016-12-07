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
func With(fields KV) FieldLogger {
	return logEntry.With(fields)
}

// Err creates a new LoggerEntry from the standard logger and adds an error
// to it, using the value defined in ErrorKey as key.
func Err(err error) FieldLogger {
	return logEntry.Err(err)
}

// Fields will return the current set of fields in the global context
func Fields() KV {
	return logEntry.Fields()
}

// Ctx will use the provided context with its logs if applicable
func Ctx(ctx context.Context) FieldLogger {
	return logEntry.Ctx(ctx)
}

// NewContext adds the current `logEntry` into `ctx`
func NewContext(ctx context.Context) context.Context {
	return logEntry.NewContext(ctx)
}

// AppendContext creates a new context.Context from the supplied ctx with the fields appended to the end
func AppendContext(ctx context.Context, fields KV) context.Context {
	return logEntry.AppendContext(ctx, fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logEntry.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logEntry.Info(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logEntry.Print(args...)
}

// Warn logs a message at level Warning on the standard logger.
func Warn(args ...interface{}) {
	logEntry.Warn(args...)
}

// Warning logs a message at level Warning on the standard logger.
func Warning(args ...interface{}) {
	logEntry.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logEntry.Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger
func Fatal(args ...interface{}) {
	logEntry.Fatal(args)
}

// Panic logs a message at level Panic on the standard logger
func Panic(args ...interface{}) {
	logEntry.Panic(args)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logEntry.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logEntry.Infof(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logEntry.Printf(format, args...)
}

// Warnf logs a message at level Warning on the standard logger.
func Warnf(format string, args ...interface{}) {
	logEntry.Warnf(format, args...)
}

// Warningf logs a message at level Warning on the standard logger.
func Warningf(format string, args ...interface{}) {
	logEntry.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logEntry.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger
func Fatalf(format string, args ...interface{}) {
	logEntry.Fatal(format, args)
}

// Panicf logs a message at level Panic on the standard logger
func Panicf(format string, args ...interface{}) {
	logEntry.Panic(format, args)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logEntry.Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logEntry.Infoln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logEntry.Println(args...)
}

// Warnln logs a message at level Warning on the standard logger.
func Warnln(args ...interface{}) {
	logEntry.Warnln(args...)
}

// Warningln logs a message at level Warning on the standard logger.
func Warningln(args ...interface{}) {
	logEntry.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logEntry.Errorln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger
func Fatalln(args ...interface{}) {
	logEntry.Fatalln(args)
}

// Panicln logs a message at level Panic on the standard logger
func Panicln(args ...interface{}) {
	logEntry.Panicln(args)
}
