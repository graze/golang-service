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

	"github.com/MindscapeHQ/raygun4go"
	"github.com/graze/golang-service/handlers/failure"
	"github.com/graze/golang-service/log"
)

// raygunClient allows us to mock a client for testing porpoises
type raygunClient interface {
	Request(*http.Request) *raygun4go.Client
	CustomData(interface{}) *raygun4go.Client
	CreateError(string) error
}

// loggerRecoverer is a local struct to implement the Recoverer interface
type raygunRecoverer struct {
	client raygunClient
}

// Recover creates a new raygun client each time as the details of each error will change per request
func (l raygunRecoverer) Handle(w http.ResponseWriter, r *http.Request, err error, status int) {
	l.client.Request(r)
	l.client.CustomData(log.Ctx(r.Context()).Fields())
	l.client.CreateError(err.Error())
}

// Raygun creates a Raygun Recoverer given the details
func Raygun(client raygunClient) failure.Handler {
	return &raygunRecoverer{client}
}
