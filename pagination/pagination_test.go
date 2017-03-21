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
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	pagination, err := New(1, 1, 5, &http.Request{})

	assert.IsType(t, Pagination{}, pagination)
	assert.Nil(t, err)
}

func TestNewDefaults(t *testing.T) {
	pagination, err := New(0, 0, 5, &http.Request{})

	assert.IsType(t, Pagination{}, pagination)
	assert.Nil(t, err)
	assert.Equal(t, defaultPageNumber, pagination.PageNumber)
	assert.Equal(t, defaultItemsPerPage, pagination.ItemsPerPage)
}

func TestTooManyItems(t *testing.T) {
	pagination, err := New(1, 10, 5, &http.Request{})

	assert.IsType(t, Pagination{}, pagination)
	assert.Equal(t, err.Error(), "The requested number of items per page (10) is greater than the maximum allowed (5)")
}

func TestOffset(t *testing.T) {
	p1, _ := New(1, 10, 10, &http.Request{})
	assert.Equal(t, 0, p1.Offset())

	p2, _ := New(3, 10, 10, &http.Request{})
	assert.Equal(t, 20, p2.Offset())

	p3, _ := New(3, 5, 10, &http.Request{})
	assert.Equal(t, 10, p3.Offset())
}

func TestSetItemsTotal(t *testing.T) {
	p1, _ := New(1, 10, 10, &http.Request{})
	p1.SetItemsTotal(100)
	assert.Equal(t, 100, *p1.ItemsTotal)
	assert.Equal(t, 10, *p1.PagesTotal)

	p2, _ := New(1, 10, 10, &http.Request{})
	p2.SetItemsTotal(53)
	assert.Equal(t, 53, *p2.ItemsTotal)
	assert.Equal(t, 6, *p2.PagesTotal)
}

func TestUnknownTotal(t *testing.T) {
	p1, _ := New(1, 10, 10, &http.Request{})

	assert.Nil(t, p1.ItemsTotal)
	assert.Nil(t, p1.PagesTotal)
}

func TestPageUrl(t *testing.T) {
	p, _ := New(1, 10, 10, &http.Request{
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
	p, _ := New(1, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})
	p.SetItemsTotal(100)

	json, err := p.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":1,"pages_total":10,"items_per_page":10,"items_per_page_limit":10,"items_total":100,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=10\u0026test=test","next_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=2\u0026test=test","prev_href":null}`, string(json))
}

func TestToJsonLastPage(t *testing.T) {
	p, _ := New(10, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})
	p.SetItemsTotal(100)

	json, err := p.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":10,"pages_total":10,"items_per_page":10,"items_per_page_limit":10,"items_total":100,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=10\u0026test=test","next_href":null,"prev_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=9\u0026test=test"}`, string(json))
}

func TestToJsonNoTotal(t *testing.T) {
	p, _ := New(10, 10, 10, &http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "example.com:80", // host or host:port
			Path:     "path/to/data/",
			RawQuery: "test=test&foo=bar", // encoded query values, without '?'
		},
	})

	json, err := p.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, `{"page_number":10,"pages_total":null,"items_per_page":10,"items_per_page_limit":10,"items_total":null,"first_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=1\u0026test=test","last_href":null,"next_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=11\u0026test=test","prev_href":"http://example.com:80/path/to/data/?foo=bar\u0026limit=10\u0026page=9\u0026test=test"}`, string(json))
}
