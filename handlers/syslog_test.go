// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package handlers

import (
    "testing"
    "github.com/graze/golang-service/nettest"
    "github.com/stretchr/testify/assert"
    "net"
    "net/http"
    "net/http/httptest"
    "os"
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

    handler := SyslogHandler(okHandler)

    // A typical request with an OK response
    req := newRequest("GET", "http://example.com/")

    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)

    assert.Regexp(t, `<176>[0-9\-]+T[0-9:]+Z \w+ service\.application\[\d+\]:  - - \[\d{2}\/\w+\/\d{4}:\d{2}:\d{2}:\d{2} [\+\-]+\d+\] "GET \/ HTTP\/1\.1" 200 \d+ "" ""`, <-done)
}
