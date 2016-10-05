// Package golang-service is a set of packages that provide many tools for helping create services in golang
//
// golang-service contains the following packages:
//
// The logging package provides a set of http.Handler logging handlers to write specific logs about requests
//
// The testing package provides a set of helpers for use when testing
package golangservice

import (
    _ "github.com/graze/golang-service/logging"
    _ "github.com/graze/golang-service/testing"
)
