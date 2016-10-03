# Http Service Logging Helpers

### Healthd Logger

- Support the healthd logs from AWS Elastic Beanstalk logs: [AWS](http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)

### Statsd Logger

- Output `response_time` and `count` statistics for each request to a statsd host

## Usage

```bash
$ go get github.com/graze/golang-service/logging
```

```go
func statsdHandler(h http.Handler) http.Handler {
    client, err := stats.New("<ip>:<port>")
    if err != nil {
        panic(err)
    }
    return logging.StatsdHandler(client, h)
}
```

## Development

### Testing
To run tests, run this on your host machine:

```
$ make install
$ make test
```

# License

- General code: [MIT License](LICENSE)
- some code: `Copyright (c) 2013 The Gorilla Handlers Authors. All rights reserved.`
