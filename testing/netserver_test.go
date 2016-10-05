// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// @license https://github.com/graze/golang-service-logging/blob/master/LICENSE
// @link    https://github.com/graze/golang-service-logging

package testing

import (
    "github.com/stretchr/testify/assert"
    "testing"
    "net"
    "os"
    "fmt"
)

func TestCreateTCPServer(t *testing.T) {
    cases := map[string]struct{msg, la string}{
        "base": { "test", "localhost:"},
        "specied port": { "test", "localhost:1201"},
        "ip addresss": { "test", "127.0.0.1:"},
    }

    for k, tc := range cases {
        done := make(chan string)
        addr, sock, srvWg := CreateTCPServer(t, tc.la, done)
        defer srvWg.Wait()
        defer os.Remove(addr.String())
        defer sock.Close()

        s, err := net.Dial("tcp", addr.String())
		if err != nil {
			t.Fatalf("%s: Dial() failed: %v", k, err)
		}
		fmt.Fprintf(s, tc.msg + "\n")
		assert.Equal(t, tc.msg + "\n", <-done, "test: %s", k)
		s.Close()
    }
}
