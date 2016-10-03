# Http Service Logging Helpers

### Healthd Logger

- Support the healthd logs from AWS Elastic Beanstalk logs: (http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html)[AWS]

### Statsd Logger

- Output `response_time` and `count` statistics for each request to a statsd host

## Development

### Testing
To run tests, run this on your host machine:

```
$ make install
$ make test
```

# License

- General code: (LICENSE)[MIT License]
- some code: `Copyright (c) 2013 The Gorilla Handlers Authors. All rights reserved.`
