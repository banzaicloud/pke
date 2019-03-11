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
