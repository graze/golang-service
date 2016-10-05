// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// @license https://github.com/graze/golang-service-logging/blob/master/LICENSE
// @link    https://github.com/graze/golang-service-logging

package nettest

import (
    "github.com/stretchr/testify/assert"
    "testing"
    "net"
    "os"
    "fmt"
)

func TestCreateServer(t *testing.T) {
    cases := map[string]struct{
        net string
        msgs []string
        la string
    }{
        "base": { "tcp", []string{"test"}, "localhost:"},
        "specied port": { "tcp", []string{"test"}, "localhost:1201"},
        "ip addresss": { "tcp", []string{"test"}, "127.0.0.1:"},
        "multiple messages": { "tcp", []string{"test", "test2"}, "127.0.0.1:"},
        "udp": {"udp", []string{"test"}, "localhost:"},
        "udp multiple messages": {"udp", []string{"test", "test2"}, "localhost:0"},
    }

    for k, tc := range cases {
        done := make(chan string)
        addr, sock, srvWg := CreateServer(t, tc.net, tc.la, done)
        defer srvWg.Wait()
        defer os.Remove(addr.String())
        defer sock.Close()

        assert.Equal(t, tc.net, addr.Network())

        s, err := net.Dial(tc.net, addr.String())
		if err != nil {
			t.Fatalf("%s: Dial() failed: %v", k, err)
		}
        for _, msg := range tc.msgs {
            fmt.Fprintf(s, msg + "\n")
    		assert.Equal(t, msg + "\n", <-done, "test: %s", k)
        }
		s.Close()
    }
}
