// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package recovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/graze/golang-service/log"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	logger := log.New()
	logger.SetLevel(logrus.DebugLevel)
	hook := test.NewLocal(logger.Logger)

	loggerRecoverer := PanicLogger(logger)

	handler := New(loggerRecoverer, echoRecoverer).Handle(panicHandler)

	rec := httptest.NewRecorder()
	req := newRequest("GET", "http://example.com")

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "oh no!", rec.Body.String())
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "panic occoured", hook.LastEntry().Message)
	assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "critical_error", hook.LastEntry().Data["tag"])
	assert.Equal(t, http.StatusInternalServerError, hook.LastEntry().Data["status"])
}
