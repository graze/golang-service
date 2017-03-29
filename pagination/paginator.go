package pagination

import (
	"fmt"
	"math"
)

const defaultPageNumber = 1
const defaultItemsPerPage = 10

// PaginatorInterface defines what makes a paginator
type PaginatorInterface interface {
	Init(PageNumber int, ItemsPerPage int, ItemsPerPageLimit int) (err error)
	Offset() int
	SetItemsTotal(i int)
	Validate() (err error)
}

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
func New(PageNumber int, ItemsPerPage int, ItemsPerPageLimit int) (p *Paginator, err error) {
	p = &Paginator{}
	err = p.Init(PageNumber, ItemsPerPage, ItemsPerPageLimit)
	return
}

// Init configs the Paginator struct, it is used so that we can have child structs
// inherit Pagination - see JSON for an example
func (p *Paginator) Init(PageNumber int, ItemsPerPage int, ItemsPerPageLimit int) (err error) {
	if ItemsPerPage > ItemsPerPageLimit {
		err = &TooManyItemsPerPageError{ItemsPerPage, ItemsPerPageLimit}
		return
	}

	if PageNumber < 0 {
		err = InvalidPageNumberError(PageNumber)
		return
	}

	if ItemsPerPage < 0 {
		err = InvalidItemsPerPageError(ItemsPerPage)
		return
	}

	// defaults
	if 0 == PageNumber {
		PageNumber = defaultPageNumber
	}

	if 0 == ItemsPerPage {
		ItemsPerPage = defaultItemsPerPage
	}

	p.PageNumber = PageNumber
	p.ItemsPerPage = ItemsPerPage
	p.ItemsPerPageLimit = ItemsPerPageLimit

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
func (p *Paginator) SetItemsTotal(i int) {
	pt := int(math.Ceil(float64(i) / float64(p.ItemsPerPage)))
	p.PagesTotal = &pt
	p.ItemsTotal = &i
}

// Validate provides validation to make sure that given its current state, everything is ok
func (p *Paginator) Validate() (err error) {
	// Total pages is optional, but if its set and the current page
	// is greater than available, fail validation
	if nil != p.PagesTotal && p.PageNumber > *p.PagesTotal {
		err = InvalidPageNumberError(p.PageNumber)
		return
	}

	return
}
