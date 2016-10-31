// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package handlers

import (
	"net/http"
)

// AllHandlers returns log context, structured, statsd and healthd chained handlers for use with AWS and Graze services
func AllHandlers(h http.Handler) http.Handler {
	return HealthdHandler(StatsdHandler(StructuredHandler(LogContextHandler(h))))
}
