# Validation

A very simple validation for user input

```bash
$ go get github.com/graze/golang-service/validate
```

It uses the interface: `Validator` containing the method `Validate` which will return an error if the item is not valid

```go
type Validator interface {
    Validate(ctx context.Context) error
}
```

```go
type Item struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Winning     bool   `json:"winning"`
}
func (i Item) Validate(ctx context.Context) error {
    if utf8.RuneCountInString(i.Name) <= 3 {
        return fmt.Errorf("field: name must have more than 3 characters")
    }
    return nil
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    item := &Item{}
    if err := validate.JSONRequest(ctx, r, item); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // item.Name, item.Description, item.Winning etc...
}
```
