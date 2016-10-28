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
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	logContext = New()
	AppName    = "LOG_APPLICATION"
	EnvName    = "ENVIRONMENT"
)

type F logrus.Fields

// Logger reduces the number of functions nad changes WithFields to With for
type LogContext interface {
	With(fields F) *Context
	Err(err error) *Context
	Add(fields F)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
}

type Context struct {
	*logrus.Entry
}

// With adds fields
func (c *Context) With(fields F) *Context {
	// type conversion of same type without refection
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}
	entry := c.Entry.WithFields(data)
	return &Context{entry}
}

// Err adds an error and returns a new context
func (c *Context) Err(err error) *Context {
	entry := c.Entry.WithError(err)
	return &Context{entry}
}

// Add will add the fields specified to the current context for future use
func (c *Context) Add(fields F) {
	for k, v := range fields {
		c.Entry.Data[k] = v
	}
}

// Create a new Context
func New() (context *Context) {
	base := logrus.New()
	context = &Context{logrus.NewEntry(base)}
	fields := make(F)
	if app := os.Getenv(AppName); app != "" {
		fields["app"] = app
	}
	if env := os.Getenv(EnvName); env != "" {
		fields["env"] = env
	}
	context.Add(fields)
	return
}
