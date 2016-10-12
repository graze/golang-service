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
Package logging provides a collection of logging http.Handlers for use by HTTP services.

Healthd

This provides healthd logging (http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)
for health monitoring when using Elastic Beanstalk.

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   w.Write([]byte("This is a catch-all route"))
    })
    loggedRouter := logging.HealthdHandler(os.Stdout, r)
    http.ListenAndServe(":1123", loggedRouter)

Statsd

Log request duration to a statsd host

Usage:
    r := mux.NewRouter()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	   w.Write([]byte("This is a catch-all route"))
    })
    c, err := statsd.New("127.0.0.1:8125")
    loggedRouter := logging.StatsdHandler(c, r)
    http.ListenAndServe(":1123", loggedRouter)
*/
package logging
