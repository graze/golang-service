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
Package auth provides a collection of authentication http.Handlers for use by HTTP services

Api Key Auth

For a basic api key based authentication

Usage:
    func finder(key string, r *http.Request) (interface{}, error) {
        user, ok := users[key]
        if !ok {
            return nil, fmt.Errorf("No user found for: %s", key)
        }
        return user, nil
    }

    func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
        w.WriteHeader(status)
        fmt.Fprintf(w, err.Error())
    }

    keyAuth := auth.ApiKey{
        Provider: "Graze",
        Finder: finder,
        OnError: onError,
    }

    http.Handle("/", keyAuth.Next(router))

Authentication can be added to a handler chain too:

    handlers := first(second(keyAuth.Then(fourth(r)))

Or using something like `alice`:

    chain := alice.New(first, second, keyAuth.Handler, fourth)

User Retrieval

The authentication also adds the user field returned by the finder to the
context object on the request. It can be retireved using the `auth.GetUser` method.

    func GetList(w http.ResponseWriter, r *http.Request) {
        user, ok := auth.GetUser(r).(*account.User)
        if !ok {
            w.WriteHeader(403)
            return
        }
    }
*/
package auth
