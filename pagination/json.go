package pagination

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// JSON extends a basic Paginator but adds the request object
type JSON struct {
	Paginator
	r *http.Request
}

// NewJSON creates a new JSON paginator, which extends a basic Paginator but adds the request object
func NewJSON(pageNumber int, itemsPerPage int, itemsPerPageLimit int, r *http.Request) (j *JSON, err error) {
	j = new(JSON)
	err = j.Init(pageNumber, itemsPerPage, itemsPerPageLimit)
	j.r = r
	return
}

// JSONFields is a seperate to JSON that defines the JSON formatting
type JSONFields struct {
	PageNumber        int     `json:"page_number"`
	PagesTotal        *int    `json:"pages_total"`
	ItemsPerPage      int     `json:"items_per_page"`
	ItemsPerPageLimit int     `json:"items_per_page_limit"`
	ItemsTotal        *int    `json:"items_total"`
	FirstHref         string  `json:"first_href"`
	LastHref          *string `json:"last_href"`
	NextHref          *string `json:"next_href"`
	PrevHref          *string `json:"prev_href"`
}

// MarshalJSON is called when ever the Pagination object is json encoded, we
// use this chance to set the correct href values (as they depend on other data being set, like number of results)
func (j JSON) MarshalJSON() ([]byte, error) {

	jf := JSONFields{
		PageNumber:        j.PageNumber,
		PagesTotal:        j.PagesTotal,
		ItemsPerPage:      j.ItemsPerPage,
		ItemsPerPageLimit: j.ItemsPerPageLimit,
		ItemsTotal:        j.ItemsTotal,
		FirstHref:         j.pageURL(1).String(),
	}

	if j.PagesTotal != nil {
		last := j.pageURL(*j.PagesTotal).String()
		jf.LastHref = &last
	}

	if j.PagesTotal == nil || j.PageNumber < *j.PagesTotal {
		next := j.pageURL(j.PageNumber + 1).String()
		jf.NextHref = &next
	}

	if j.PageNumber > 1 {
		prev := j.pageURL(j.PageNumber - 1).String()
		jf.PrevHref = &prev
	}

	return json.Marshal(jf)
}

// pageURL is a helper used to build the HREFs based on the current URL
// it really just overrides the currently set page query string with the required value
func (j *JSON) pageURL(page int) (u *url.URL) {
	u = j.r.URL

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("limit", strconv.Itoa(j.ItemsPerPage))
	u.RawQuery = q.Encode()

	// The URL does not always contain the required information, so if its not there
	// build it from the request info
	if "" == u.Host {
		u.Host = j.r.Host
	}

	// use a secure default, this can be set in the URL directly
	// or by the X-Forwarded-Proto and RFC7239 headers
	if "" == u.Scheme {
		u.Scheme = "https"
	}

	return
}
