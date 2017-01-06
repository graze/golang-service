// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

/*
Package raygun sends recovered panic errors to raygun from within an http.Handler

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       panic("oh-o")
    })

    outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
        w.Write([]byte("panic happened, oh dear"))
    }

    raygunClient, _ := raygun4go.New(name, key)
    raygunClient.Silent(false)
    raygunClient.Version("1.0")

    raygunHandler = raygun.New(raygunClient)
    recoverer := recovery.New(r, raygunHandler)
    http.ListenAndServe(":80", recoverer)
*/
package raygun

import (
	"io"
	"net/http"

	"github.com/MindscapeHQ/raygun4go"
	"github.com/graze/golang-service/handlers/recovery"
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
func (l raygunRecoverer) Handle(w io.Writer, r *http.Request, err error, status int) {
	l.client.Request(r)
	l.client.CustomData(log.Ctx(r.Context()).Fields())
	l.client.CreateError(err.Error())
}

// New creates a Raygun Recoverer given the details
func New(client raygunClient) recovery.Handler {
	return &raygunRecoverer{client}
}
