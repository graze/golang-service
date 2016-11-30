// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

/*
Package validate provides a simple interface for validating JSON user input

Example:

	type Item struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

    func (i *Item) Validate(ctx context.Context) error {
		if RuneCountInString(i.Name) == 0 {
			return fmt.Errorf("the field: name must be provided and not empty")
		}
		return nil
	}

	func CreateItem(w http.ResponseWriter, r *http.Request) {
		item := &Item{}
		if err := validate.JsonRequest(ctx, r, item); err != nil {
			w.WriteHeader(400)
			return
		}
	}
*/
package validate

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	// IOError for when we fail to read the stream
	IOError struct{ error }
)

// Validatable items can self validate
type Validatable interface {
	Validate(ctx context.Context) error
}

// JsonRequest takes an http request, decodes the json and validates the input against a Validatable output
// the Validatable variable will get populated with the contents of the body provided by *http.Request
//
// Usage:
// 	type Item struct {
//		Name string `json:"name"`
//		Description string `json:"description"`
//	}
//
//	func (i *Item) Validate(ctx context.Context) error {
//		if RuneCountInString(i.Name) == 0 {
//			return fmt.Errorf("the field: name must be provided and not empty")
//		}
//		return nil
//	}
//
//	func CreateItem(w http.ResponseWriter, r *http.Request) {
//		item := &Item{}
//		if err := validate.JsonRequest(ctx, r, item); err != nil {
//			w.WriteHeader(400)
//			return
//		}
//	}
func JsonRequest(ctx context.Context, r *http.Request, v Validatable) error {
	return Reader(ctx, r.Body, json.Unmarshal, v)
}

// XmlRequest takes an http request docodes the xml into an item and validates the provided item
func XmlRequest(ctx context.Context, r *http.Request, v Validatable) error {
	return Reader(ctx, r.Body, xml.Unmarshal, v)
}

// ReadAndValidate takes a generic io.Reader, an unmarshaller  and validates the input against a Validatable item
// the Validatable variable will get populated with the contents of the body provided by *http.Request
//
// Usage:
// 	type ApiInput struct {
//		Name string `json:"name"`
//		Description string `json:"description"`
//	}
//
//  func (i *ApiInput) Validate(ctx context.Context) error {
//		if RuneCountInString(i.Name) == 0 {
//			return fmt.Errorf("the field: name must be provided and not empty")
//		}
//		return nil
//	}
//
//	func main() {
//		input := &ApiInput{}
//		if err := validate.Reader(ctx, reader, json.Unmarshal, input); err != nil {
//			log.Panic(err)
//		}
//	}
func Reader(ctx context.Context, r io.Reader, unmarshaller func(data []byte, v interface{}) error, v Validatable) error {
	str, err := ioutil.ReadAll(r)
	if err != nil {
		return &IOError{err}
	}
	if err = unmarshaller(str, v); err != nil {
		return err
	}
	return v.Validate(ctx)
}
