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
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/DataDog/datadog-go/statsd"
    "../nettest"
    "time"
    "os"
    "net/http"
    "net/url"
)

func TestStatsdLogging(t *testing.T) {
    getDuration := func(dur string) (duration time.Duration) {
        duration, err := time.ParseDuration(dur)
        if err != nil {
            t.Fatal(err)
        }
        return
    }

    cases := map[string]struct{
        request *http.Request
        timestamp time.Time
        duration time.Duration
        expected []string
    }{
        "basic": {
            newRequest("GET", "http://example.com"),
            time.Now().UTC(),
            getDuration("0.302s"),
            []string{
                "service.logging.live.request.response_time:302.000000|ms|#test,endpoint:/,statusCode:200,method:GET,protocol:HTTP/1.1",
                "service.logging.live.request.count:1|c|#test,endpoint:/,statusCode:200,method:GET,protocol:HTTP/1.1",
            },
        },
        "post path": {
            newRequest("POST", "http://example.com/path/here"),
            time.Now().UTC(),
            getDuration("0.102s"),
            []string{
                "service.logging.live.request.response_time:102.000000|ms|#test,endpoint:/path/here,statusCode:200,method:POST,protocol:HTTP/1.1",
                "service.logging.live.request.count:1|c|#test,endpoint:/path/here,statusCode:200,method:POST,protocol:HTTP/1.1",
            },
        },
        "strips params off method": {
            newRequest("GET", "http://example.com/token/1/test?apid=1&thing=2"),
            time.Now().UTC(),
            getDuration("0.927s"),
            []string{
                "service.logging.live.request.response_time:927.000000|ms|#test,endpoint:/token/1/test,statusCode:200,method:GET,protocol:HTTP/1.1",
                "service.logging.live.request.count:1|c|#test,endpoint:/token/1/test,statusCode:200,method:GET,protocol:HTTP/1.1",
            },
        },
        "connect http2 test": {
            &http.Request{
                Method:     "CONNECT",
                Proto:      "HTTP/2.0",
                ProtoMajor: 2,
                ProtoMinor: 0,
                URL:        &url.URL{Host: "www.example.com:443"},
                Host:       "www.example.com:443",
                RemoteAddr: "192.168.100.5",
            },
            time.Now().UTC(),
            getDuration("0.927s"),
            []string{
                "service.logging.live.request.response_time:927.000000|ms|#test,endpoint:www.example.com:443,statusCode:200,method:CONNECT,protocol:HTTP/2.0",
                "service.logging.live.request.count:1|c|#test,endpoint:www.example.com:443,statusCode:200,method:CONNECT,protocol:HTTP/2.0",
            },
        },
    }

    done := make(chan string)
    addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
    defer srvWg.Wait()
    defer os.Remove(addr.String())
    defer sock.Close()

	client, err := statsd.New(addr.String())
	if err != nil {
		t.Fatal(err)
	}
    client.Tags = append(client.Tags, "test")
    client.Namespace = "service.logging.live."

    for k, tc := range cases {
        writeStatsdLog(client, tc.request, *tc.request.URL, tc.timestamp, tc.duration, http.StatusOK, 100)
        for _, message := range tc.expected {
            assert.Equal(t, message, <-done, "test: %s", k)
        }
    }
}
