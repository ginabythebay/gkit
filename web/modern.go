// +build !appengine

package web

import (
	"net/http"

	"golang.org/x/net/context"
)

func getContext(r *http.Request) context.Context {
	return r.Context()
}

func getTransport(ctx context.Context) http.RoundTripper {
	return http.DefaultTransport
}
