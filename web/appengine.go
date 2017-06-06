// +build appengine

package web

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func getContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func getTransport(ctx context.Context) http.RoundTripper {
	return &urlfetch.Transport{Context: ctx}
}
