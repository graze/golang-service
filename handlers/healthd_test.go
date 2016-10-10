package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "os"
    "time"
    "io/ioutil"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok\n"))
})

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHealthdLogging(t *testing.T) {
    handler := HealthdHandler(okHandler)

    // A typical request with an OK response
    req := newRequest("GET", "http://example.com/")

    rec := httptest.NewRecorder()
    timestamp := time.Now().UTC().Format("2006-01-02-15")
    handler.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)

    file := "/var/log/nginx/healthd/application.log." + timestamp

    fe, err := exists("/var/log/nginx/healthd")
    if err != nil {
        t.Fatal(err)
    }
    assert.True(t, fe)
    fe, err = exists(file)
    if err != nil {
        t.Fatal(err)
    }
    assert.True(t, fe)

    bytes, err := ioutil.ReadFile(file)
    if err != nil {
        t.Fatal(err)
    }

    assert.Regexp(t, `[0-9\.]+"/"200"[0-9\.]+"[0-9\.]+"[0-9\.]*`, string(bytes))
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}
