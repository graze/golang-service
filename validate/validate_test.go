// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package validate

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TypesStruct struct {
	Number int               `json:"int" XMLName:"int"`
	Text   string            `json:"string" XMLName:"string"`
	Arr    []string          `json:"array" XMLName:"array"`
	Obj    map[string]string `json:"object" XMLName:"object"`
}

func (t *TypesStruct) Validate(ctx context.Context) error {
	return nil
}

type XMLStruct struct {
	Name   string `XmlName:"name"`
	Drinks int    `XmlName:"drinks"`
}

type emptyStruct struct{}

func (s *emptyStruct) Validate(ctx context.Context) error {
	return nil
}

func (x *XMLStruct) Validate(ctx context.Context) error {
	return nil
}

type FailureStruct map[string]interface{}

func (s *FailureStruct) Validate(ctx context.Context) error {
	return errors.New("this is not valid")
}

func newRequest(t *testing.T, method, path, body string) *http.Request {
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		t.Errorf("Failed to create a request")
		t.Fail()
	}

	return req
}

func TestReader(t *testing.T) {
	cases := map[string]struct {
		Data         string
		Unmarshaller func(data []byte, v interface{}) error
		Struct       Validatable
		Errored      bool
		Expected     string
	}{
		"base": {
			`
{
"int": 1,
"string": "thing",
"array": ["a","b"],
"object": {
"child": "thing"
}
}`,
			json.Unmarshal,
			&TypesStruct{},
			false,
			"",
		},
		"invalid json": {
			`{"some":"text",}`,
			json.Unmarshal,
			&TypesStruct{},
			true,
			"invalid character '}' looking for beginning of object key string",
		},
		"fails validation": {
			`{"some":"text"}`,
			json.Unmarshal,
			&FailureStruct{},
			true,
			"this is not valid",
		},
		"xml validation": {
			`<XmlStruct><Name>bob</Name><Drinks>12</Drinks></XmlStruct>`,
			xml.Unmarshal,
			&XMLStruct{},
			false,
			"",
		},
	}

	for k, tc := range cases {
		err := Reader(context.Background(), strings.NewReader(tc.Data), tc.Unmarshaller, tc.Struct)
		assert.Equal(t, tc.Errored, err != nil, "test: %s", k)
		if tc.Errored {
			assert.Equal(t, tc.Expected, err.Error(), "test: %s", k)
		}
	}
}

func TestJsonRequest(t *testing.T) {
	cases := map[string]struct {
		Request  *http.Request
		Struct   Validatable
		Errored  bool
		Expected string
	}{
		"base": {
			newRequest(t, "POST", "/thing", `
{
"int": 1,
"string": "thing",
"array": ["a","b"],
"object": {
"child": "thing"
}
}`),
			&TypesStruct{},
			false,
			"",
		},
		"empty": {
			newRequest(t, "POST", "/thing", ``),
			&emptyStruct{},
			false,
			"",
		},
	}

	for k, tc := range cases {
		err := JSONRequest(context.Background(), tc.Request, tc.Struct)
		assert.Equal(t, tc.Errored, err != nil, "test: %s", k)
		if tc.Errored {
			assert.Equal(t, tc.Expected, err.Error(), "test: %s", k)
		}
	}
}

func TestXmlRequest(t *testing.T) {
	cases := map[string]struct {
		Request  *http.Request
		Struct   Validatable
		Errored  bool
		Expected string
	}{
		"xml validation": {
			newRequest(t, "POST", "/thing", `<XmlStruct><Name>bob</Name><Drinks>12</Drinks></XmlStruct>`),
			&XMLStruct{},
			false,
			"",
		},
	}

	for k, tc := range cases {
		err := XMLRequest(context.Background(), tc.Request, tc.Struct)
		assert.Equal(t, tc.Errored, err != nil, "test: %s", k)
		if tc.Errored {
			assert.Equal(t, tc.Expected, err.Error(), "test: %s", k)
		}
	}
}

func TestIOError(t *testing.T) {
	err := IOError{fmt.Errorf("some error")}
	assert.Equal(t, "some error", err.Error())
}
