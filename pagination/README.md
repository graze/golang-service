
# Pagination

The purpose of pagination is to break up large sets of information (i.e., HTTP responses) into smaller, more manageable ones.

### Pagination Considerations

Pagination adds complexity to both the service and its clients. The decision to paginate a resource should consider the following:

- expected response size
    - e.g. what is the average response size? what are the outliers? how fast are the responses growing?
- complexity of retrieving the resources
    - e.g. what is the performance impact of retrieving and processing X resources?
- how fast is the data changing?
    - is simple (slow) or complex (fast) pagination most appropriate?
- needs of the client(s)
    - e.g. does the client really want/need X records upfront when it is only capable of displaying X/10?
    - e.g. limited or unreliable connectivity (i.e. mobile) - what volume of data can the client handle?

### Implementation

```
r := &http.Request{} // replace with your request object
currentPage := 10 // from user input
userLimit := 100 // from user input
maxLimit := 250

pag, err := pagination.New(currentPage, userLimit, maxLimit, r)
if err != nil {
	panic(err)
}
count := getSomeData()
pag.SetItemsTotal(count)

json := json.Marshal(pag) // return this along with your data to the end user
```

### Implementation of Simple Pagination in a RESTful JSON API

A resource is considered paginated if a `pagination` attribute is present in the JSON response.

```json
{
    "data": { },
    "pagination": { }
}
```

The `pagination` attribute contains metadata about how the resource is paginated, including a list of URLs for the client(s) to navigate between paged resources. The general idea is to limit the amount of work a client has to do to implement a paginated API. e.g. avoid clients having to build URLs


e.g. https://api-example.com/resource?page=3
```json
{
    "data": { },
    "pagination":   
    {
        "page_number": 3,
        "pages_total": 46,

        "items_per_page": 15,        
        "items_per_page_limit": 25,
        "items_total": 689,

        "first_href": "https://api-example.com/resource?page=1",
        "last_href": "https://api-example.com/resource?page=45",
        "next_href": "https://api-example.com/resource?page=4",
        "previous_href": "https://api-example.com/resource?page=2"
    }
}
```
__Pagination Attributes__
- Pages (integers)
    - `page_number` the index of the current page
    - `pages_total`  the total number of pages (i.e. round up the number of items / current page size )
- Items (integers)
    - `items_per_page` the current number of items per page
    - `items_per_page_limit` the maximum amount of items that can be returned per page
    - `items_total` the total number of items
- URLs (strings)
    - `first_href` full URL of the first page of resources
    - `last_href` full URL of the last page of resources
    - `next_href` full URL of the next page of resources. `null` if this is the last page
    - `previous_href` full URL of the previous page of resources. `null` if this is the first page

### The `page` URL Parameter

The `page` URL parameter is an integer that represents the page indices of a given resource. e.g. `https://api.com/resource?page=12`,

The following can be considered synonymous:

1. `https://api.com/resource`
2. `https://api.com/resource?page=`
3. `https://api.com/resource?page=1`

However, all URLs in the `pagination` attribute should contain the fullest representation (3).

### The `limit` URL Parameter

The `limit` URL parameter is an integer that specifies the number of items that should be returned in a page. This should be reflected in the first, last, next, and previous URLs.


e.g. https://api-example.com/resource?page=3&limit=10
```json
{
    "data": { },
    "pagination":   
    {
        "page_number": 3,
        "pages_total": 68,

        "items_per_page": 10,
        "items_per_page_limit": 15,
        "items_total": 675,

        "first_href": "https://api-example.com/resource?page=1&limit=10",
        "last_href": "https://api-example.com/resource?page=68&limit=10",
        "next_href": "https://api-example.com/resource?page=4&limit=10",
        "previous_href": "https://api-example.com/resource?page=2&limit=10"
    }
}
```