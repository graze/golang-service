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
	context = New()
	AppName = "LOG_APPLICATION"
	EnvName = "ENVIRONMENT"
)

// Logger reduces the number of functions nad changes WithFields to With for
type LogContext interface {
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
}

// Create a new base LogContext
func New() (entry *logrus.Entry) {
	base := logrus.New()
	fields := make(logrus.Fields)
	if app := os.Getenv(AppName); app != "" {
		fields["app"] = app
	}
	if env := os.Getenv(EnvName); env != "" {
		fields["env"] = env
	}
	entry = base.WithFields(fields)
	return
}
