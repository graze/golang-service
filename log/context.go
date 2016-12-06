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
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	logEntry  = New()
	appName   = "LOG_APPLICATION"
	envName   = "ENVIRONMENT"
	levelName = "LOG_LEVEL"
)

// key is a type to ensure unique key for context
type key int

// LogKey is the key used for context
const logKey key = iota

// KV is a shorthand for logrus.Fields so less text is required to be typed:
//
// 	log.With(log.KV{"k":"v"})
type KV logrus.Fields

// FieldLogger represents a Logging FieldLogger
type FieldLogger interface {
	Ctx(ctx context.Context) FieldLogger
	NewContext(ctx context.Context) context.Context
	AppendContext(ctx context.Context, fields KV) context.Context

	With(fields KV) FieldLogger
	Err(err error) FieldLogger

	Fields() KV

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
}

// Logger represents a struct that can modify the output of a log
type Logger interface {
	SetOutput(out io.Writer)
	SetLevel(level logrus.Level)
	Level() logrus.Level
	SetFormatter(formatter logrus.Formatter)
	AddHook(hook logrus.Hook)
}

// LoggerEntry is a logging context that can be passed around
type LoggerEntry struct {
	*logrus.Entry
}

// NewContext returns the provided context with this LoggerEntry added
func (c *LoggerEntry) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey, c.Fields())
}

// Ctx will use any logging context contained in context.Context to append to the current log context
func (c *LoggerEntry) Ctx(ctx context.Context) FieldLogger {
	if fields, ok := ctx.Value(logKey).(KV); ok {
		return c.With(fields)
	}
	return c.With(KV{})
}

// AppendContext will create a new context.Context based on ctx with the fields appended
func (c *LoggerEntry) AppendContext(ctx context.Context, fields KV) context.Context {
	return c.Ctx(ctx).With(fields).NewContext(ctx)
}

// With creates a new LoggerEntry and adds the fields to it
func (c *LoggerEntry) With(fields KV) FieldLogger {
	// type conversion of same type without refection
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}
	entry := c.Entry.WithFields(data)
	return &LoggerEntry{entry}
}

// Err adds an error and returns a new LoggerEntry
func (c *LoggerEntry) Err(err error) FieldLogger {
	entry := c.Entry.WithError(err)
	return &LoggerEntry{entry}
}

// Fields will return the current fields attached to a context
func (c *LoggerEntry) Fields() (fields KV) {
	fields = make(KV, len(c.Entry.Data))
	for k, v := range c.Entry.Data {
		fields[k] = v
	}
	return
}

// SetOutput changes the output of the current context
func (c *LoggerEntry) SetOutput(out io.Writer) {
	c.Logger.Out = out
}

// SetFormatter will change the formatter for the current context
func (c *LoggerEntry) SetFormatter(formatter logrus.Formatter) {
	c.Logger.Formatter = formatter
}

// SetLevel changes the default logging level of the current context
func (c *LoggerEntry) SetLevel(level logrus.Level) {
	c.Logger.Level = level
}

// Level returns the current logging level this context will log at
func (c *LoggerEntry) Level() (level logrus.Level) {
	return c.Logger.Level
}

// AddHook will add a hook to the current context
func (c *LoggerEntry) AddHook(hook logrus.Hook) {
	c.Logger.Hooks.Add(hook)
}

// New creates a new FieldLogger with a new Logger (formatter, level, output, hooks)
func New() (entry *LoggerEntry) {
	base := logrus.New()
	logger := &LoggerEntry{logrus.NewEntry(base)}
	fields := make(KV)
	if app := os.Getenv(appName); app != "" {
		fields["app"] = app
	}
	if env := os.Getenv(envName); env != "" {
		fields["env"] = env
	}
	if level := os.Getenv(levelName); level != "" {
		if l, err := logrus.ParseLevel(level); err == nil {
			logger.SetLevel(l)
		} else {
			logger.Err(err).With(KV{
				"module":   "log_initialisation",
				"tag":      "log_new_failed",
				"logLevel": level,
			}).Error("The supplied log level is invalid")
		}
	}
	entry, _ = logger.With(fields).(*LoggerEntry)
	return
}
