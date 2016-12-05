# Metrics

Provides statsd metrics sending

Manually create a client
```go
client, _ := metrics.GetStatsd(StatdClientConf{host,port,namespace,tags})
client.Incr("metric", []string{}, 1)
```

Create a client using environment variables. [see above](#statsd-logger)
```go
client, _ := metrics.GetStatsdFromEnv()
client.Incr("metric", []string{}, 1)
```
