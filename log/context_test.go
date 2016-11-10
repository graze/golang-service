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

func TestNewContext(t *testing.T) {
	original := New()
	logger := original.With(KV{"test": "test2"})
	assert.NotEqual(t, original, logger)
	hook := test.NewLocal(logger.Logger)

	logger.Info("test")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "test", hook.LastEntry().Message)
	assert.Equal(t, InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])

	logger.With(KV{"2": 3}).Error("error")
	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, "error", hook.LastEntry().Message)
	assert.Equal(t, ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 3, hook.LastEntry().Data["2"])
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])
}

func TestMergeContext(t *testing.T) {
	logger := New().With(KV{"test": 1})

	logger2 := New().With(KV{"test2": 2})

	assert.Equal(t, KV{"test": 1}, logger.Fields())
	assert.Equal(t, KV{"test2": 2}, logger2.Fields())

	assert.Equal(t, KV{"test": 1, "test2": 2}, logger.With(logger2.Fields()).Fields())
}

func testImplements(t *testing.T) {
	logger := New()
	assert.Implements(t, (*Logger)(nil), logger)
	assert.Implements(t, (*FieldLogger)(nil), logger)

	local := logger.With(KV{"k": "v"})
	assert.Implements(t, (*FieldLogger)(nil), local)

	global := With(KV{"k": "v"})
	assert.Implements(t, (*FieldLogger)(nil), global)
}

func TestNewWithValidLogLevels(t *testing.T) {
	cases := map[string]struct {
		level    string
		expected logrus.Level
	}{
		"lower case": {"info", logrus.InfoLevel},
		"upper case": {"INFO", logrus.InfoLevel},
		"mixed case": {"InFo", logrus.InfoLevel},
		"debug":      {"debug", logrus.DebugLevel},
		"panic":      {"panic", logrus.PanicLevel},
		"warn":       {"warn", logrus.WarnLevel},
		"warning":    {"warning", logrus.WarnLevel},
		"fatal":      {"fatal", logrus.FatalLevel},
		"error":      {"error", logrus.ErrorLevel},
	}

	for k, tc := range cases {
		os.Setenv("LOG_LEVEL", tc.level)
		logger := New()
		assert.Equal(t, tc.expected, logger.Level(), "test: %s", k)
	}
	os.Setenv("LOG_LEVEL", "")
}

func TestNewWithInvalidLogLEvels(t *testing.T) {
	cases := map[string]struct {
		level    string
		expected logrus.Level
	}{
		"plural": {"infos", logrus.InfoLevel},
		"crit":   {"crit", logrus.InfoLevel},
	}

	for k, tc := range cases {
		os.Setenv("LOG_LEVEL", tc.level)
		logger := New()
		assert.Equal(t, tc.expected, logger.Level(), "test: %s", k)
	}
	os.Setenv("LOG_LEVEL", "")
}
