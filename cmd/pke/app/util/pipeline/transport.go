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

package pipeline

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type transportLogger struct {
	roundTripper http.RoundTripper
	output       io.Writer
}

func (t *transportLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), "requestTS", time.Now())
	req = req.WithContext(ctx)

	_, _ = fmt.Fprintf(t.output, "%s --> %s %q\n", req.Proto, req.Method, req.URL)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	ctx = resp.Request.Context()
	if ts, ok := ctx.Value("requestTS").(time.Time); ok {
		_, _ = fmt.Fprintf(t.output, "%s <-- %d %q %s\n", resp.Proto, resp.StatusCode, resp.Request.URL, time.Now().Sub(ts))
	} else {
		_, _ = fmt.Fprintf(t.output, "%s <-- %d %q\n", resp.Proto, resp.StatusCode, resp.Request.URL)
		if resp.StatusCode/100 != 2 {
			if b, err := httputil.DumpResponse(resp, true); err != nil {
				_, _ = fmt.Fprintf(t.output, "%s\n", b)
			}

		}
	}

	return resp, err

}
func (t *transportLogger) transport() http.RoundTripper {
	if t.roundTripper != nil {
		return t.roundTripper
	}

	return http.DefaultTransport
}
