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

func TestStatsdLogging(t *testing.T) {
    tests := map[string]struct{
        request *http.Request
        expected []string
    }{
        "simple get": { newRequest("GET", "http://example.com"), []string{
            "request.response_time:0.000000|ms|#endpoint:/,statusCode:200,method:GET,protocol:HTTP/1.1",
            "request.count:1|c|#endpoint:/,statusCode:200,method:GET,protocol:HTTP/1.1",
        }},
        "post removes fields": { newRequest("POST", "http://example.com/token?apid=1"), []string{
            "request.response_time:0.000000|ms|#endpoint:/token,statusCode:200,method:POST,protocol:HTTP/1.1",
            "request.count:1|c|#endpoint:/token,statusCode:200,method:POST,protocol:HTTP/1.1",
        }},
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
    os.Setenv("STATSD_HOST", host)
    os.Setenv("STATSD_PORT", port)

    handler := StatsdHandler(okHandler)

    for k, tc := range tests {
        rec := httptest.NewRecorder()
        handler.ServeHTTP(rec, tc.request)

        assert.Equal(t, http.StatusOK, rec.Code)

        for _, message := range tc.expected {
            assert.Equal(t, message, <-done, "test: %s", k)
        }
    }
}

func TestStatsdLoggingWithNamespaceAndTags(t *testing.T) {
    done := make(chan string)
    addr, sock, srvWg := nettest.CreateServer(t, "udp", "localhost:", done)
    defer srvWg.Wait()
    defer os.Remove(addr.String())
    defer sock.Close()

    host, port, err := net.SplitHostPort(addr.String())
    if err != nil {
        t.Fatal(err)
    }
    os.Setenv("STATSD_HOST", host)
    os.Setenv("STATSD_PORT", port)
    os.Setenv("STATSD_NAMESPACE", "service.")
    os.Setenv("STATSD_TAGS", "tag1:value,tag2")

    handler := StatsdHandler(okHandler)

	// A typical request with an OK response
	req := newRequest("POST", "http://example.com/some/path")

    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)

    expected := []string{
        "service.request.response_time:0.000000|ms|#tag1:value,tag2,endpoint:/some/path,statusCode:200,method:POST,protocol:HTTP/1.1",
        "service.request.count:1|c|#tag1:value,tag2,endpoint:/some/path,statusCode:200,method:POST,protocol:HTTP/1.1",
    }

    for _, message := range expected {
        assert.Equal(t, message, <-done)
    }
}
