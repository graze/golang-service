package logging

import (
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	context := New()
	assert.Exactly(t, context, context.Add(F{"test": "test2"}))
	hook := test.NewLocal(context.Logger)

	context.Info("test")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "test", hook.LastEntry().Message)
	assert.Equal(t, InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])

	context.With(F{"2": 3}).Error("error")
	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, "error", hook.LastEntry().Message)
	assert.Equal(t, ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 3, hook.LastEntry().Data["2"])
	assert.Equal(t, "test2", hook.LastEntry().Data["test"])
}

func TestMergeContext(t *testing.T) {
	context := New()
	context.Add(F{"test": 1})

	context2 := New()
	context2.Add(F{"test2": 2})

	assert.Equal(t, F{"test": 1}, context.Get())
	assert.Equal(t, F{"test2": 2}, context2.Get())

	assert.Exactly(t, context, context.Merge(context2))
	assert.Equal(t, F{"test": 1, "test2": 2}, context.Get())
}
