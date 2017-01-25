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

Common components

All authentication types use a similar Finder to retrieve user information based on the supplied credentials.
For this there are Finder and FinderFunc types which supplies a Find method to retrieve a user based on some credentials

    type Finder interface {
        Find(interface{}, *http.Request) (interace{}, error)
    }

The FinderFunc converts a function to a Finder interface

The Finder and FailHandler are common to all authentication types

    func finder(key interface{}, r *http.Request) (interface{}, error) {
        k, ok := key.(string)
        if !ok {
            return nil, fmt.Errorf("The supplied key is in an invald format")
        }
        user, ok := users[k]
        if !ok {
            return nil, fmt.Errorf("No user found for: %s", key)
        }
        return user, nil
    }

    func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
        w.WriteHeader(status)
        fmt.Fprintf(w, err.Error())
    }

Authorization Bearer Api Key Auth

For a basic api key based authentication. It directly passes the apiKey as a the credentials to the Finder.Func method

Usage:
    keyAuth := auth.NewAPIKey("Graze", auth.FinderFunc(finder), failure.HandlerFunc(onError))

    http.Handle("/", keyAuth.Next(router))

X-Api-Key Authorization

Almost identical to the Authorization header, is using the X-Api-Key header to simply provide just they key to handle.
It uses the same Finder and onError.

Usage:
    keyAuth := auth.NewXApiKey(auth.FinderFunc(finder), failure.HandlerFunc(onError))

    http.Handle("/", keyAuth.Next(router))

Usage

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
