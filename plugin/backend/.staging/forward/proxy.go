package forwardstorage

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sniperkit/xcache"
)

type key int

const originKey key = 1

// proxy is the forward caching proxy on a peer, it uses
// a cache that conforms to the HTTP RFC (thanks to
// github.com/gregjones/httpcache)
type proxy struct {
	path string
	*httputil.ReverseProxy
}

// newProxy creates a proxy that serves requests on path using the
// specified cache. The proxy handles requests of format:
// /path?q=originUrl where originUrl is the resource being
// requested by the client.
func newProxy(path string, cache httpcache.Cache, transport http.RoundTripper) *proxy {
	return &proxy{
		path: path,
		ReverseProxy: &httputil.ReverseProxy{
			Transport: &httpcache.Transport{
				Cache:               cache,
				MarkCachedResponses: true,
				Transport:           transport,
			},
			Director: director,
		},
	}
}

// ServeHTTP takes the url of the requested resource to be fetched on the
// origin and puts in in the request's context to be used later by the proxy director.
func (p *proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != p.path {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	q := req.URL.Query().Get("q")
	if q == "" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	origin, err := url.Parse(q)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	ctx := context.WithValue(req.Context(), originKey, origin)
	p.ReverseProxy.ServeHTTP(w, req.WithContext(ctx))
}

// director modifies the requested URL to the origin.
func director(req *http.Request) {
	origin := req.Context().Value(originKey).(*url.URL)
	req.URL = origin
	req.Host = origin.Host
}
