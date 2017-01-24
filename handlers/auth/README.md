# Authentication Handler

```bash
$ go get github.com/graze/golang-service/handlers/auth
```

Authentication provides a little bit of security to your service.

## Common components

Some common components are available during authentication:

- `finder` (`Finder`) - Asks the application (you) if a supplied set of credentials are valid
- `onError` (`FailHandler`) - When an error occurs during authentication, this is called so the application (you again) can handle it nicely

```go
func finder(creds interface{}, r *http.Request) (interface{}, error) {
    key, ok := creds.(string)
    if !ok {
        return nil, fmt.Error("Invalid credentials format, expecting string")
    }
    user, ok := users[key]
    if !ok {
        return nil, fmt.Errorf("No user found for: %s", key)
    }
    return user, nil
}

func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
    w.WriteHeader(status)
    fmt.Fprintf(w, err.Error())
}
```

## API Key Authentication

Adds authentication to the request using middleware, with the benefit of linking the authentication with a user

```http
GET / HTTP/1.1
Host: service.example.com
Content-Type: application/json
Authorization: Graze hzYAVO9Sg98nsNh81M84O2kyXVy6K1xwHD8
```

```go
keyAuth := auth.NewAPIKey("Graze", auth.FinderFunc(finder), failure.HandlerFunc(onError))

http.Handle("/", keyAuth.Next(router))
```

## X-Api-Key Authentication

The header key: `X-Api-Key` can also be used for authentication by simply providing the key as the value of the header.

```http
GET / HTTP/1.1
Host: service.example.com
Content-Type: application/json
x-api-key: hzYAVO9Sg98nsNh81M84O2kyXVy6K1xwHD8
```

The same Finder and methods can be used

```go
keyAuth := auth.NewXAPIKey(auth.FinderFunc(finder), failure.HandlerFunc(onError))

http.Handle("/", keyAuth.Next(router))
```

### User Retrieval

You can then retrieve the user provided by the `Finder` function within the request handler:

```go
func GetList(w http.ResponseWriter, r *http.Request) {
    user, ok := auth.GetUser(r).(*account.User)
    if !ok {
        w.WriteHeader(http.StatusForbidden)
        return
    }
}
```
