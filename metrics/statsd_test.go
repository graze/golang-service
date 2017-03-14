// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package metrics

import (
	"net"
	"os"
	"testing"

	"github.com/graze/golang-service/nettest"
	"github.com/stretchr/testify/assert"
)

func TestStatsdLogging(t *testing.T) {
	tests := map[string]struct {
		metric   string
		value    float64
		tags     []string
		expected string
	}{
		"base test": {
			"request.count",
			1,
			[]string{},
			"request.count:1|c",
		},
		"tags test": {
			"request.tags",
			1,
			[]string{"tag1", "tag2:value"},
			"request.tags:1|c|#tag1,tag2:value",
		},
	}

	done := make(chan string)
	addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
	defer srvWg.Wait()
	defer os.Remove(addr.String())
	defer sock.Close()

	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		t.Fatal(err)
	}

	c := StatsdClientConf{
		Host: host,
		Port: port,
	}
	client, err := GetStatsd(c)
	if err != nil {
		t.Fatal(err)
	}

	for k, tc := range tests {
		client.Incr(tc.metric, tc.tags, tc.value)
		assert.Equal(t, tc.expected, <-done, "test: %s", k)
	}
}

func TestStatsdLoggingWithNamespaceAndTags(t *testing.T) {
	tests := map[string]struct {
		metric   string
		value    float64
		tags     []string
		expected string
	}{
		"base test": {
			"request.count",
			1,
			[]string{},
			"service.request.count:1|c|#tag1:value,tag2",
		},
		"extra tags": {
			"request.count",
			1,
			[]string{"tag3,tag4:thing"},
			"service.request.count:1|c|#tag1:value,tag2,tag3,tag4:thing",
		},
	}

	done := make(chan string)
	addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
	defer srvWg.Wait()
	defer os.Remove(addr.String())
	defer sock.Close()

	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		t.Fatal(err)
	}

	c := StatsdClientConf{
		host,
		port,
		"service.",
		[]string{"tag1:value", "tag2"},
	}
	client, err := GetStatsd(c)
	if err != nil {
		t.Fatal(err)
	}

	for k, tc := range tests {
		client.Incr(tc.metric, tc.tags, tc.value)
		assert.Equal(t, tc.expected, <-done, "test: %s", k)
	}
}
