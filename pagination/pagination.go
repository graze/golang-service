package pagination

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const defaultPageNumber = 1
const defaultItemsPerPage = 10

// Pagination contains all of the JSON response fields as well as a url
// that is later used in the json encoding to populate the hrefs dynamiclly
// some fields are pointers so they can be a null JSON value if empty
type Pagination struct {
	PageNumber        int     `json:"page_number"`
	PagesTotal        *int    `json:"pages_total"`
	ItemsPerPage      int     `json:"items_per_page"`
	ItemsPerPageLimit int     `json:"items_per_page_limit"`
	ItemsTotal        *int    `json:"items_total"`
	FirstHref         string  `json:"first_href"`
	LastHref          *string `json:"last_href"`
	NextHref          *string `json:"next_href"`
	PrevHref          *string `json:"prev_href"`
	request           *http.Request
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
	return fmt.Sprintf("The requested page (%d) is less than 1", int(e))
}

// InvalidItemsPerPageError is the error generated when the requested items per page is invalid
type InvalidItemsPerPageError int

func (e InvalidItemsPerPageError) Error() string {
	return fmt.Sprintf("The requested items per page (%d) is less than 1", int(e))
}

// Offset returns the current offset (i.e. for database queries)
// based on the current page and number of itmes
func (p *Pagination) Offset() int {
	if p.PageNumber == 1 {
		return 0
	}
	return (p.PageNumber - 1) * p.ItemsPerPage
}

// SetItemsTotal will set the ItemsTotal value as well as the dynamic PagesTotal based on ItemsPerPage
func (p *Pagination) SetItemsTotal(i int) {
	p.ItemsTotal = &i
	pt := int(math.Ceil(float64(i) / float64(p.ItemsPerPage)))
	p.PagesTotal = &pt
}

// pageURL is a helper used to build the HREFs based on the current URL
// it really just overrides the currently set page query string with the required value
func (p *Pagination) pageURL(page int) (u *url.URL) {
	u = p.request.URL

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(p.ItemsPerPage))
	u.RawQuery = q.Encode()

	// The URL does not always contain the required information, so if its not there
	// build it from the request info
	if "" == u.Host {
		u.Host = p.request.Host
	}
	if "" == u.Scheme {
		u.Scheme = "http"
		if nil != p.request.TLS {
			u.Scheme = "https"
		}
	}

	return
}

// New returns a configured Pagination struct that can be used
func New(PageNumber int, ItemsPerPage int, ItemsPerPageLimit int, r *http.Request) (p Pagination, err error) {
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

	return Pagination{
		PageNumber:        PageNumber,
		ItemsPerPage:      ItemsPerPage,
		ItemsPerPageLimit: ItemsPerPageLimit,
		request:           r,
	}, err
}

// MarshalJSON is called when ever the Pagination object is json encoded, we
// use this chance to set the correct href values (as they depend on other data being set, like number of results)
func (p Pagination) MarshalJSON() ([]byte, error) {

	// Set our HREFs last minute
	p.FirstHref = p.pageURL(1).String()

	if p.PagesTotal != nil {
		last := p.pageURL(*p.PagesTotal).String()
		p.LastHref = &last
	}

	if p.PagesTotal == nil || p.PageNumber < *p.PagesTotal {
		next := p.pageURL(p.PageNumber + 1).String()
		p.NextHref = &next
	}

	if p.PageNumber > 1 {
		prev := p.pageURL(p.PageNumber - 1).String()
		p.PrevHref = &prev
	}

	// Use an Alias to prevent infinite loop of this function
	type Alias Pagination
	b, err := json.Marshal(&struct {
		Alias
	}{
		(Alias)(p),
	})

	return b, err
}
