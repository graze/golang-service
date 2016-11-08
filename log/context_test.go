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
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	context := New()
	assert.Exactly(t, context, context.Add(KV{"test": "test2"}))
	hook := test.NewLocal(context.Logger)

	context.Info("test")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "test", hook.LastEntry().Message)
	assert.Equal(t, InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])

	context.With(KV{"2": 3}).Error("error")
	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, "error", hook.LastEntry().Message)
	assert.Equal(t, ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 3, hook.LastEntry().Data["2"])
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])
}

func TestMergeContext(t *testing.T) {
	context := New()
	context.Add(KV{"test": 1})

	context2 := New()
	context2.Add(KV{"test2": 2})

	assert.Equal(t, KV{"test": 1}, context.Fields())
	assert.Equal(t, KV{"test2": 2}, context2.Fields())

	assert.Exactly(t, context, context.Merge(context2))
	assert.Equal(t, KV{"test": 1, "test2": 2}, context.Fields())
}

func testImplements(t *testing.T) {
	context := New()
	assert.Implements(t, (*Logger)(nil), context)
	assert.Implements(t, (*Context)(nil), context)

	local := context.With(KV{"k": "v"})
	assert.Implements(t, (*Context)(nil), local)

	global := With(KV{"k": "v"})
	assert.Implements(t, (*Context)(nil), global)
}