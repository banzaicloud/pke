// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/PuerkitoBio/rehttp"
)

func NewRetryTransport(rt http.RoundTripper) *rehttp.Transport {
	return rehttp.NewTransport(
		rt,
		rehttp.RetryAny(
			rehttp.RetryAll(
				rehttp.RetryHTTPMethods(http.MethodGet),
				rehttp.RetryStatusInterval(400, 600),
			),
			rehttp.RetryAll(
				rehttp.RetryHTTPMethods(http.MethodPost),
				rehttp.RetryStatusInterval(500, 600),
			),
			rehttp.RetryAll(
				rehttp.RetryTemporaryErr(),
			),
			rehttp.RetryAll(
				RetryConnectionRefusedErr(),
			),
		),
		rehttp.ExpJitterDelay(2*time.Second, 60*time.Second),
	)
}

func RetryConnectionRefusedErr() rehttp.RetryFn {
	return func(attempt rehttp.Attempt) bool {
		if operr, ok := attempt.Error.(*net.OpError); ok {
			if syserr, ok := operr.Err.(*os.SyscallError); ok {
				if syserr.Err == syscall.ECONNREFUSED {
					return true
				}
			}
		}

		return false
	}
}
