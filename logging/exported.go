// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package logging

import (
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
	logrus.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	return logrus.GetLevel()
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

// WithFields returns a standard logger with the context fields
func With(fields F) *Context {
	return logContext.With(fields)
}

// WithError creates an entry from the standard logger and adds an error
// to it, using the value defined in ErrorKey as key.
func Err(err error) *Context {
	return logContext.Err(err)
}

// Add modifies the global context
func Add(fields F) *Context {
	logContext.Add(fields)
	return logContext
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logContext.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logContext.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logContext.Info(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logContext.Error(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logContext.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logContext.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logContext.Infof(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logContext.Errorf(format, args...)
}
