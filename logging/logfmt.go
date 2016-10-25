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
	"os"

	log "github.com/Sirupsen/logrus"
)

// init initialises a logger
func GetLogger() (logger *log.Entry) {
	logger = log.NewEntry(log.New())

	if os.Getenv("LOG_APPLICATION") != "" {
		logger = logger.WithField("application", os.Getenv("LOG_APPLICATION"))
	}
	return
}
