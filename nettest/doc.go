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
Package nettest provides a set of network helpers for when unit testing networks

Uses a channel log the messages recieved by the server

Usage:
    done := make(chan string)
    addr, sock, srvWg := CreateServer(t, tc.net, tc.la, done)
    defer srvWg.Wait()
    defer os.Remove(addr.String())
    defer sock.Close()

    s, err := net.Dial(tc.net, addr.String())
    defer s.Close()
    fmt.Fprintf(s, "test message\n")
    if "test message\n" != <-done {
        t.Error("message not recieved")
    }
*/
package nettest
