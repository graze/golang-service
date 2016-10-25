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

	"github.com/go-kit/kit/log"
)

var logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

// GetLogger returns a global opinionated logger that outputs to Stderr in Logfmt format
func GetLogger() (logger log.Logger) {
	return logger
}
