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
	"io"
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	logContext = New()
	appName    = "LOG_APPLICATION"
	envName    = "ENVIRONMENT"
)

// F is a shorthand for logrus.Fields so less text is required to be typed:
//
// 	log.With(log.F{"k":"v"})
type F logrus.Fields

// Context represents a Logging ContextEntry
type Context interface {
	With(fields F) *ContextEntry
	Err(err error) *ContextEntry
	Add(fields F) *ContextEntry
	Merge(context Context) *ContextEntry
	Get() F

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
	GetLevel() logrus.Level
	SetFormatter(formatter logrus.Formatter)
	AddHook(hook logrus.Hook)
}

// ContextEntry is a logging context that can be passed around
type ContextEntry struct {
	*logrus.Entry
}

// With creates a new `ContextEntry` and adds the fields to it
func (c *ContextEntry) With(fields F) *ContextEntry {
	// type conversion of same type without refection
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}
	entry := c.Entry.WithFields(data)
	return &ContextEntry{entry}
}

// Err adds an error and returns a new `ContextEntry`
func (c *ContextEntry) Err(err error) *ContextEntry {
	entry := c.Entry.WithError(err)
	return &ContextEntry{entry}
}

// Add adds the fields to the current `ContextEntry` and returns itself
func (c *ContextEntry) Add(fields F) *ContextEntry {
	for k, v := range fields {
		c.Entry.Data[k] = v
	}
	return c
}

// Merge will merge the fields in the supplied `context` into this `ContextEntry`
func (c *ContextEntry) Merge(context Context) *ContextEntry {
	return c.Add(context.Get())
}

// Get will return the current fields attached to a context
func (c *ContextEntry) Get() (fields F) {
	fields = make(F, len(c.Entry.Data))
	for k, v := range c.Entry.Data {
		fields[k] = v
	}
	return
}

// SetOutput changes the output of the current context
func (c *ContextEntry) SetOutput(out io.Writer) {
	c.Logger.Out = out
}

// SetFormatter will change the formatter for the current context
func (c *ContextEntry) SetFormatter(formatter logrus.Formatter) {
	c.Logger.Formatter = formatter
}

// SetLevel changes the default logging level of the current context
func (c *ContextEntry) SetLevel(level logrus.Level) {
	c.Logger.Level = level
}

// GetLevel returns the current logging level this context will log at
func (c *ContextEntry) GetLevel() (level logrus.Level) {
	return c.Logger.Level
}

// AddHook will add a hook to the current context
func (c *ContextEntry) AddHook(hook logrus.Hook) {
	c.Logger.Hooks.Add(hook)
}

// New creates a new ContextEntry with a new Logger context (formatter, level, output, hooks)
func New() (context *ContextEntry) {
	base := logrus.New()
	context = &ContextEntry{logrus.NewEntry(base)}
	fields := make(F)
	if app := os.Getenv(appName); app != "" {
		fields["app"] = app
	}
	if env := os.Getenv(envName); env != "" {
		fields["env"] = env
	}
	context.Add(fields)
	return
}
