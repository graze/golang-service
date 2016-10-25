// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package logging

import (
	"fmt"
	"log/syslog"
	"os"
	"strconv"
)

// SyslogConf is a configuration struct to be able to create a SysLog writer
type SyslogConf struct {
	network                 string
	host, port, application string
	level                   syslog.Priority
}

// GetSysLog returns a syslog.Writer based on a SyslogConf input
func GetSysLog(conf SyslogConf) (logWriter *syslog.Writer, err error) {
	la := ""
	if conf.host != "" {
		la = conf.host + ":" + conf.port
	}
	logWriter, err = syslog.Dial(
		conf.network,
		la,
		conf.level,
		conf.application)
	return
}

// GetSysLogFromEnv creates a syslog.Writer based on the environment variables supplied
//
// Uses the application environments:
//
//  SYSLOG_NETWORK: The network type of the syslog server (tcp, udp) Leave blank for local syslog
//  SYSLOG_HOST: The host of the syslog server. Leave blank for local syslog
//  SYSLOG_PORT: The port of the syslog server
//  SYSLOG_APPLICATION: The application to report the logs as
//  SYSLOG_LEVEL: The level to limit messages to (default: LEVEL6)
//
// Returns a *syslog.Writer
func GetSysLogFromEnv() (logWriter *syslog.Writer, err error) {
	names := []struct {
		name, env string
		req       bool
	}{
		{"network", "SYSLOG_NETWORK", false},
		{"host", "SYSLOG_HOST", false},
		{"port", "SYSLOG_PORT", false},
		{"application", "SYSLOG_APPLICATION", true},
		{"level", "SYSLOG_LEVEL", false},
	}

	env := make(map[string]string)
	for _, v := range names {
		env[v.name] = os.Getenv(v.env)
		if v.req && env[v.name] == "" {
			return nil, fmt.Errorf("Unable to get value for setting: %s", v.env)
		}
	}

	level := syslog.LOG_LOCAL6 | syslog.LOG_NOTICE
	if env["level"] != "" {
		fmt.Println(env["level"])
		i, err := strconv.Atoi(env["level"])
		if err != nil {
			return nil, err
		}
		level = syslog.Priority(i)
		fmt.Println(level)
	}

	conf := SyslogConf{env["network"], env["host"], env["port"], env["application"], level}
	return GetSysLog(conf)
}
