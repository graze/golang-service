// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package pagination

import (
	"encoding/json"
	"net/http"
	"net/url"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJson(t *testing.T) {
	p, err := NewJSON(1, 3, 5, &http.Request{})

	assert.IsType(t, &JSON{}, p)
	assert.Nil(t, err)

	assert.Equal(t, 1, p.PageNumber)
	assert.Equal(t, 3, p.ItemsPerPage)
	assert.Equal(t, 5, p.ItemsPerPageLimit)
}

func TestPageUrl(t *testing.T) {
	p, _ := NewJSON(1, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})

	assert.Equal(t, "http://example.com:80/path/to/data/?foo=bar&limit=10&page=1&test=test", p.pageURL(1).String())
}

func TestToJsonFirstPage(t *testing.T) {
	p, _ := NewJSON(1, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})
	p.SetItemsTotal(100)

	json, err := json.Marshal(p)

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":1,"pages_total":10,"items_per_page":10,"items_per_page_limit":10,"items_total":100,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=10\u0026test=test","next_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=2\u0026test=test","prev_href":null}`, string(json))
}

func TestToJsonLastPage(t *testing.T) {
	p, _ := NewJSON(10, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})
	p.SetItemsTotal(100)

	json, err := json.Marshal(p)

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":10,"pages_total":10,"items_per_page":10,"items_per_page_limit":10,"items_total":100,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=10\u0026test=test","next_href":null,"prev_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=9\u0026test=test"}`, string(json))
}

func TestToJsonNoTotal(t *testing.T) {
	p, _ := NewJSON(10, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})

	json, err := json.Marshal(p)

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":10,"pages_total":null,"items_per_page":10,"items_per_page_limit":10,"items_total":null,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":null,"next_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=11\u0026test=test","prev_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=9\u0026test=test"}`, string(json))
}
