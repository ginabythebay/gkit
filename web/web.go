package web

import (
	"net/http"

	"golang.org/x/net/context"
)

// Context returns a context either for app engine, or for a more modern system
func Context(r *http.Request) context.Context {
	return getContext(r)
}

// Transport returns an http transport.
func Transport(ctx context.Context) http.RoundTripper {
	return getTransport(ctx)
}
