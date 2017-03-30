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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	paginator, err := New(1, 1, 5)

	assert.IsType(t, &Paginator{}, paginator)
	assert.Nil(t, err)
}

func TestInitError(t *testing.T) {
	paginator := &Paginator{}
	err := paginator.Init(0, 500, 5)

	assert.IsType(t, &Paginator{}, paginator)
	assert.NotNil(t, err)
}

func TestInitOk(t *testing.T) {
	paginator := &Paginator{}
	err := paginator.Init(1, 1, 5)

	assert.IsType(t, &Paginator{}, paginator)
	assert.Nil(t, err)
}

func TestNewDefaults(t *testing.T) {
	paginator, err := New(0, 0, 5)

	assert.IsType(t, &Paginator{}, paginator)
	assert.Nil(t, err)
	assert.Equal(t, defaultPageNumber, paginator.PageNumber)
	assert.Equal(t, defaultItemsPerPage, paginator.ItemsPerPage)
}

func TestTooManyItems(t *testing.T) {
	paginator, err := New(1, 10, 5)

	assert.IsType(t, &Paginator{}, paginator)
	assert.Equal(t, err.Error(), "The requested number of items per page (10) is greater than the maximum allowed (5)")
}

func TestOffset(t *testing.T) {
	p1, _ := New(1, 10, 10)
	assert.Equal(t, 0, p1.Offset())

	p2, _ := New(3, 10, 10)
	assert.Equal(t, 20, p2.Offset())

	p3, _ := New(3, 5, 10)
	assert.Equal(t, 10, p3.Offset())
}

func TestSetItemsTotal(t *testing.T) {
	p1, _ := New(1, 10, 10)
	p1.SetItemsTotal(100)
	assert.Equal(t, 100, *p1.ItemsTotal)
	assert.Equal(t, 10, *p1.PagesTotal)

	p2, _ := New(1, 10, 10)
	p2.SetItemsTotal(53)
	assert.Equal(t, 53, *p2.ItemsTotal)
	assert.Equal(t, 6, *p2.PagesTotal)
}

func TestUnknownTotal(t *testing.T) {
	p1, _ := New(1, 10, 10)

	assert.Nil(t, p1.ItemsTotal)
	assert.Nil(t, p1.PagesTotal)
}

// With SetItemsTotal()
func TestPageToHighWithCount(t *testing.T) {
	p, _ := New(1000, 10, 10)
	err := p.SetItemsTotal(10)

	assert.NotNil(t, err)
}
