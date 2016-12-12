// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package auth

import "net/http"

// Finder will return an user/account for a given set of credentials
//
// Usage:
// 	type Accounts struct {
//		users map[string]User
//  }
//
//  func (a Accounts) Find(c interface{}, r *http.Request) (interface{}, error) {
// 		key, ok := c.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("The supplied key is in an invald format")
// 		}
// 		user, ok := users[key]
// 		if !ok {
// 			return nil, fmt.Errorf("No user found for: %s", key)
// 		}
// 		return user, nil
// 	}
type Finder interface {
	Find(credentials interface{}, r *http.Request) (interface{}, error)
}

// FinderFunc is a method wrapper around the Finder interface
type FinderFunc func(interface{}, *http.Request) (interface{}, error)

// Find returns the
func (f FinderFunc) Find(c interface{}, r *http.Request) (interface{}, error) {
	return f(c, r)
}
