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
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	original := New("", "", "")
	logger := original.With(KV{"test": "test2"})
	assert.NotEqual(t, original, logger)

	base, ok := logger.(*LoggerEntry)
	if !ok {
		t.Error("unable to convert logger to *LoggerEntry")
	}
	hook := test.NewLocal(base.Logger)

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
	logger := New("", "", "").With(KV{"test": 1})

	logger2 := New("", "", "").With(KV{"test2": 2})

	assert.Equal(t, KV{"test": 1}, logger.Fields())
	assert.Equal(t, KV{"test2": 2}, logger2.Fields())

	assert.Equal(t, KV{"test": 1, "test2": 2}, logger.With(logger2.Fields()).Fields())
}

func testImplements(t *testing.T) {
	logger := New("", "", "")
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
		"blank":      {"", logrus.InfoLevel},
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
		logger := New("", "", tc.level)
		assert.Equal(t, tc.expected, logger.Level(), "test: %s", k)
	}
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
		logger := New("", "", tc.level)
		assert.Equal(t, tc.expected, logger.Level(), "test: %s", k)
	}
}

func testAppendContext(t *testing.T) {
	ctx := With(KV{"key": "value"}).NewContext(context.Background())

	ctx = AppendContext(ctx, KV{"key2": "value2"})

	assert.Equal(t, KV{"key": "value", "key2": "value2"}, Ctx(ctx).Fields())

	logger := New("", "", "").With(KV{"key": "value"})
	ctx = logger.NewContext(context.Background())

	logger2 := New("", "", "")

	ctx = logger2.AppendContext(ctx, KV{"key2": "value2"})

	assert.Equal(t, KV{"key": "value", "key2": "value2"}, logger2.Ctx(ctx).Fields())
	assert.Equal(t, KV{}, logger2.Fields())
}

type newKey int

const (
	keyOne newKey = iota
	keyTwo
)

func TestNewContextKeepsOldContextValues(t *testing.T) {
	ctx := context.WithValue(context.Background(), keyOne, "bar")
	ctx = context.WithValue(ctx, keyTwo, "foo")

	logger := New("", "", "")
	ctx = logger.With(KV{"key": "value"}).NewContext(ctx)

	assert.Equal(t, "bar", ctx.Value(keyOne))
	assert.Equal(t, "foo", ctx.Value(keyTwo))
	logger2 := New("", "", "")
	assert.Equal(t, KV{"key": "value"}, logger2.Ctx(ctx).Fields())
}
