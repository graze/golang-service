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

type F logrus.Fields

// LogContext represents a Logging Context
type LogContext interface {
	With(fields F) *Context
	Err(err error) *Context
	Add(fields F) *Context
	Merge(context LogContext) *Context
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

// Context is a logging context that can be passed around
type Context struct {
	*logrus.Entry
}

// With creates a new `Context` and adds the fields to it
func (c *Context) With(fields F) *Context {
	// type conversion of same type without refection
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}
	entry := c.Entry.WithFields(data)
	return &Context{entry}
}

// Err adds an error and returns a new `Context`
func (c *Context) Err(err error) *Context {
	entry := c.Entry.WithError(err)
	return &Context{entry}
}

// Add adds the fields to the current `Context` and returns itself
func (c *Context) Add(fields F) *Context {
	for k, v := range fields {
		c.Entry.Data[k] = v
	}
	return c
}

// Merge will merge the fields in the supplied `context` into this `Context`
func (c *Context) Merge(context LogContext) *Context {
	return c.Add(context.Get())
}

// Get will return the current fields attached to a context
func (c *Context) Get() (fields F) {
	fields = make(F, len(c.Entry.Data))
	for k, v := range c.Entry.Data {
		fields[k] = v
	}
	return
}

// SetOutput changes the output of the current context
func (c *Context) SetOutput(out io.Writer) {
	c.Logger.Out = out
}

// SetFormatter will change the formatter for the current context
func (c *Context) SetFormatter(formatter logrus.Formatter) {
	c.Logger.Formatter = formatter
}

// SetLevel changes the default logging level of the current context
func (c *Context) SetLevel(level logrus.Level) {
	c.Logger.Level = level
}

// GetLevel returns the current logging level this context will log at
func (c *Context) GetLevel() (level logrus.Level) {
	return c.Logger.Level
}

// AddHook will add a hook to the current context
func (c *Context) AddHook(hook logrus.Hook) {
	c.Logger.Hooks.Add(hook)
}

// Create a new Context
func New() (context *Context) {
	base := logrus.New()
	context = &Context{logrus.NewEntry(base)}
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
