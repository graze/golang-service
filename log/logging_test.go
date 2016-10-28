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
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	logger, hook := test.NewNullLogger()

	logger.Info("message")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "message", hook.LastEntry().Message)
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)

	hook.Reset()

	logger.WithFields(logrus.Fields{
		"variable": 2,
	}).Error("some text")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "some text", hook.LastEntry().Message)
	assert.Equal(t, ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 2, hook.LastEntry().Data["variable"])
}

func TestEnvironment(t *testing.T) {
	os.Setenv("LOG_APPLICATION", "some_app")
	os.Setenv("ENVIRONMENT", "test")

	logger := New()
	hook := test.NewLocal(logger.Logger)

	logger.Info("some text")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "some text", hook.LastEntry().Message)
	assert.Equal(t, "some_app", hook.LastEntry().Data["app"])
	assert.Equal(t, "test", hook.LastEntry().Data["env"])
}

func TestGlobalConfiguration(t *testing.T) {
	SetOutput(os.Stdout)
	SetLevel(DebugLevel)
	SetFormatter(&logrus.JSONFormatter{})

	logger := New()

	// New() uses the default settings
	assert.Equal(t, os.Stderr, logger.Logger.Out)
	assert.Equal(t, InfoLevel, logger.Logger.Level)
	assert.IsType(t, (*logrus.TextFormatter)(nil), logger.Logger.Formatter)

	context := With(F{})

	assert.Equal(t, os.Stdout, context.Logger.Out)
	assert.Equal(t, DebugLevel, context.Logger.Level)
	assert.IsType(t, (*logrus.JSONFormatter)(nil), context.Logger.Formatter)
}

func TestModificationOfContextLogger(t *testing.T) {
	logger := New()

	// New() uses the default settings
	assert.Equal(t, os.Stderr, logger.Logger.Out)
	assert.Equal(t, InfoLevel, logger.Logger.Level)
	assert.IsType(t, (*logrus.TextFormatter)(nil), logger.Logger.Formatter)

	logger.SetOutput(os.Stdout)
	logger.SetLevel(DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	assert.Equal(t, os.Stdout, logger.Logger.Out)
	assert.Equal(t, DebugLevel, logger.Logger.Level)
	assert.IsType(t, (*logrus.JSONFormatter)(nil), logger.Logger.Formatter)
}
