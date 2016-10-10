package handlers

import (
    "github.com/DataDog/datadog-go/statsd"
    "github.com/graze/golang-service/logging"
    "net/http"
    "strings"
    "os"
    "fmt"
)

// statsdHandler returns a logging.StatsdHandler to write request and response informtion to statsd
func StatsdHandler(h http.Handler) http.Handler {
    client, err := GetStatsdFromEnv()
    if err != nil {
        panic(err)
    }
    return logging.StatsdHandler(client, h)
}

// StatsdClientConf is a configuration struct to create a StatsD client
type StatsdClientConf struct {
    host, port, namespace string
    tags []string
}

// GetStatsd returns a statsd client based on the supplied StatsdClientConf
func GetStatsd(conf StatsdClientConf) (client *statsd.Client, err error) {
    client, err = statsd.New(conf.host + ":" + conf.port)
    if err != nil {
        return nil, err
    }

    client.Namespace = conf.namespace
    client.Tags = append(client.Tags, conf.tags...)
    return
}

// GetStatsdFromEnv creates a statsd client based on the environment of the application
//
// Uses the application environments:
//
//  STATSD_HOST: The host of the statsd server
//  STATSD_PORT: The port of the statsd server
//  STATSD_NAMESPACE: The namespace to prefix to every metric name
//  STATSD_TAGS: A comma separared list of tags to apply to every metric reported
//
// Returns a statsd.Client
func GetStatsdFromEnv() (client *statsd.Client, err error) {
    names := []struct{
        name, env string
        req bool
    }{
        {"host", "STATSD_HOST", true},
        {"port", "STATSD_PORT", true},
        {"namespace", "STATSD_NAMESPACE", false},
        {"tags", "STATSD_TAGS", false},
    }

    env := make(map[string]string)

    for _, v := range names {
        env[v.name] = os.Getenv(v.env)
        if v.req && env[v.name] == "" {
            return nil, fmt.Errorf("Unable to get value for setting: %s", v.env)
        }
    }

    tags := make([]string, 0)
    for _, tag := range strings.Split(env["tags"], ",") {
        if tag != "" {
            tags = append(tags, tag)
        }
    }

    conf := StatsdClientConf{env["host"], env["port"], env["namespace"], tags}

    return GetStatsd(conf)
}
