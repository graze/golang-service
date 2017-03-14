// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package metrics

import "github.com/DataDog/datadog-go/statsd"

// StatsdClientConf is a configuration struct to create a StatsD client
type StatsdClientConf struct {
	Host, Port, Namespace string
	Tags                  []string
}

// GetStatsd returns a statsd client based on the supplied StatsdClientConf
func GetStatsd(conf StatsdClientConf) (client *statsd.Client, err error) {
	client, err = statsd.New(conf.Host + ":" + conf.Port)
	if err != nil {
		return nil, err
	}

	client.Namespace = conf.Namespace
	client.Tags = append(client.Tags, conf.Tags...)
	return
}
