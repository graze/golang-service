package logging

import (
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	context := New()
	context.Add(F{"test": "test2"})
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
}
