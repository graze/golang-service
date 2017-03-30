package pagination

import (
	"fmt"
	"math"
)

const defaultPageNumber = 1
const defaultItemsPerPage = 10

// Paginator contains all of the basic pagination fields
type Paginator struct {
	PageNumber        int
	PagesTotal        *int
	ItemsPerPage      int
	ItemsPerPageLimit int
	ItemsTotal        *int
}

// TooManyItemsPerPageError is the error generated when the requested items per page is greater than the max
type TooManyItemsPerPageError struct {
	PerPage, MaxPerPage int
}

func (e *TooManyItemsPerPageError) Error() string {
	return fmt.Sprintf("The requested number of items per page (%d) is greater than the maximum allowed (%d)", e.PerPage, e.MaxPerPage)
}

// InvalidPageNumberError is the error generated when the requested page number is invalid
type InvalidPageNumberError int

func (e InvalidPageNumberError) Error() string {
	return fmt.Sprintf("The requested page (%d) is not available", int(e))
}

// InvalidItemsPerPageError is the error generated when the requested items per page is invalid
type InvalidItemsPerPageError int

func (e InvalidItemsPerPageError) Error() string {
	return fmt.Sprintf("The requested items per page (%d) is less than 1", int(e))
}

// New returns a new Paginator whilst calling init
func New(pageNumber int, itemsPerPage int, itemsPerPageLimit int) (p *Paginator, err error) {
	p = &Paginator{}
	err = p.Init(pageNumber, itemsPerPage, itemsPerPageLimit)
	return
}

// Init configs the Paginator struct, it is used so that we can have child structs
// inherit Pagination - see JSON for an example
func (p *Paginator) Init(pageNumber int, itemsPerPage int, itemsPerPageLimit int) (err error) {
	if itemsPerPage > itemsPerPageLimit {
		err = &TooManyItemsPerPageError{itemsPerPage, itemsPerPageLimit}
		return
	}

	if pageNumber < 0 {
		err = InvalidPageNumberError(pageNumber)
		return
	}

	if itemsPerPage < 0 {
		err = InvalidItemsPerPageError(itemsPerPage)
		return
	}

	// defaults
	if 0 == pageNumber {
		pageNumber = defaultPageNumber
	}

	if 0 == itemsPerPage {
		itemsPerPage = defaultItemsPerPage
	}

	p.PageNumber = pageNumber
	p.ItemsPerPage = itemsPerPage
	p.ItemsPerPageLimit = itemsPerPageLimit

	return
}

// Offset returns the current offset (i.e. for database queries)
// based on the current page and number of itmes
func (p *Paginator) Offset() int {
	if p.PageNumber == 1 {
		return 0
	}
	return (p.PageNumber - 1) * p.ItemsPerPage
}

// SetItemsTotal will set the ItemsTotal value as well as the dynamic PagesTotal based on ItemsPerPage
func (p *Paginator) SetItemsTotal(i int) (err error) {
	pt := int(math.Ceil(float64(i) / float64(p.ItemsPerPage)))

	if p.PageNumber > 1 && p.PageNumber > pt {
		err = InvalidPageNumberError(p.PageNumber)
		return
	}

	p.PagesTotal = &pt
	p.ItemsTotal = &i
	return
}
