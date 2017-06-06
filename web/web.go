// Package web provides a layer of indirection to make it easier to
// write web apps that can run on app engine or on more modern go
// implementations.
package web

import (
	"net/http"

	"golang.org/x/net/context"
)

// Context returns a context for the current request.
func Context(r *http.Request) context.Context {
	return getContext(r)
}

// Transport returns an http transport for making requests.
func Transport(ctx context.Context) http.RoundTripper {
	return getTransport(ctx)
}
