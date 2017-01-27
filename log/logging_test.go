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
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	logger := New()
	hook := test.NewLocal(logger.Logger)

	logger.Info("message")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "message", hook.LastEntry().Message)
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)

	hook.Reset()

	logger.With(KV{"variable": 2}).Error("some text")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "some text", hook.LastEntry().Message)
	assert.Equal(t, ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 2, hook.LastEntry().Data["variable"])
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

	logger2, ok := With(KV{}).(*LoggerEntry)
	if !ok {
		t.Error("unable to cast logger to *LoggerEntry")
	}

	assert.Equal(t, os.Stdout, logger2.Logger.Out)
	assert.Equal(t, DebugLevel, logger2.Logger.Level)
	assert.IsType(t, (*logrus.JSONFormatter)(nil), logger2.Logger.Formatter)
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

func TestPassingAroundContext(t *testing.T) {
	ctx := context.Background()

	logger := Ctx(ctx).With(KV{"key": "value"})
	assert.Equal(t, KV{"key": "value"}, logger.Fields())

	ctx = logger.NewContext(ctx)

	logger = Ctx(ctx).With(KV{"key2": "value2"})
	assert.Equal(t, KV{
		"key":  "value",
		"key2": "value2",
	}, logger.Fields())

	other := New().Ctx(ctx)
	assert.Equal(t, KV{"key": "value"}, other.Fields())
}

func TestUsingContextWithGlobalLogWillNotModifyTheGlobalState(t *testing.T) {
	ctx := context.Background()

	ctx = With(KV{"key": "value"}).NewContext(ctx)

	logger1 := Ctx(ctx)
	assert.Equal(t, KV{"key": "value"}, logger1.Fields())

	assert.Equal(t, KV{}, Fields())
}

func TestTheContextDoesNotContainAPointerToTheLogger(t *testing.T) {
	ctx := context.Background()

	logger := With(KV{"key": "value"})
	ctx = logger.NewContext(ctx)

	logger = logger.With(KV{"key2": "value2"})

	assert.NotEqual(t, KV{
		"key":  "value",
		"key2": "value2",
	}, Ctx(ctx).Fields())
}

func TestLevels(t *testing.T) {
	logger := New()
	logger.SetLevel(logrus.DebugLevel)
	hook := test.NewLocal(logger.Logger)

	cases := map[string]struct {
		level logrus.Level
		fn    func()
	}{
		"debug": {
			logrus.DebugLevel,
			func() { logger.Debug("debug") },
		},
		"debugf": {
			logrus.DebugLevel,
			func() { logger.Debugf("debug") },
		},
		"debugln": {
			logrus.DebugLevel,
			func() { logger.Debugln("debug") },
		},
		"info": {
			logrus.InfoLevel,
			func() { logger.Info("info") },
		},
		"infof": {
			logrus.InfoLevel,
			func() { logger.Infof("info") },
		},
		"infoln": {
			logrus.InfoLevel,
			func() { logger.Infoln("info") },
		},
		"print": {
			logrus.InfoLevel,
			func() { logger.Print("print") },
		},
		"printf": {
			logrus.InfoLevel,
			func() { logger.Printf("print") },
		},
		"println": {
			logrus.InfoLevel,
			func() { logger.Println("print") },
		},
		"warn": {
			logrus.WarnLevel,
			func() { logger.Warn("warn") },
		},
		"warnf": {
			logrus.WarnLevel,
			func() { logger.Warnf("Warn") },
		},
		"warnln": {
			logrus.WarnLevel,
			func() { logger.Warnln("Warn") },
		},
		"warning": {
			logrus.WarnLevel,
			func() { logger.Warning("warn") },
		},
		"warningf": {
			logrus.WarnLevel,
			func() { logger.Warningf("Warn") },
		},
		"warningln": {
			logrus.WarnLevel,
			func() { logger.Warningln("Warn") },
		},
		"error": {
			logrus.ErrorLevel,
			func() { logger.Error("error") },
		},
		"errorf": {
			logrus.ErrorLevel,
			func() { logger.Errorf("error") },
		},
		"errorln": {
			logrus.ErrorLevel,
			func() { logger.Errorln("error") },
		},
	}

	for k, tc := range cases {
		tc.fn()
		assert.Equal(t, tc.level, hook.LastEntry().Level, "test: %s", k)
		hook.Reset()
	}
}
