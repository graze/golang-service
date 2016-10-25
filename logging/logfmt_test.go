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
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogfmtInit(t *testing.T) {
	logger, hook := test.NewNullLogger()

	logger.Info("message")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "message", hook.LastEntry().Message)
	assert.Equal(t, log.InfoLevel, hook.LastEntry().Level)

	hook.Reset()

	logger.WithFields(log.Fields{
		"variable": 2,
	}).Error("some text")
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "some text", hook.LastEntry().Message)
	assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, 2, hook.LastEntry().Data["variable"])
}
