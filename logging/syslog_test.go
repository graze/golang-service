// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package logging

import (
	"net"
	"os"
	"testing"

	"github.com/graze/golang-service/nettest"
	"github.com/stretchr/testify/assert"
)

func TestRSyslogLogging(t *testing.T) {
	done := make(chan string)
	addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
	defer srvWg.Wait()
	defer os.Remove(addr.String())
	defer sock.Close()

	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("SYSLOG_NETWORK", "udp")
	os.Setenv("SYSLOG_HOST", host)
	os.Setenv("SYSLOG_PORT", port)
	os.Setenv("SYSLOG_APPLICATION", "service.application")
	os.Setenv("SYSLOG_LEVEL", "181")

	logWriter, err := GetSysLogFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	defer logWriter.Close()

	logWriter.Write([]byte("some text"))

	assert.Regexp(t, `<181>[0-9\-]+T[0-9:]+Z \w+ service\.application\[\d+\]: some text`+"\n", <-done)
}
