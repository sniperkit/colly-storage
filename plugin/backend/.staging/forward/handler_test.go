package forwardstorage

import (
	"hash/crc32"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sniperkit/xcache"
)

var origin *httptest.Server
var localProxy *httptest.Server
var peerProxy *httptest.Server
var pool *Pool

func TestMain(m *testing.M) {
	setup()
	status := m.Run()
	teardown()
	os.Exit(status)
}

func TestPool(t *testing.T) {
	tests := []struct {
		origin string
		status int
		cached bool
	}{
		{origin.URL + "/jquery-3.1.1.js?buster=123", http.StatusOK, false},
		{origin.URL + "/jquery-3.1.1.js?buster=123", http.StatusOK, true},
		{origin.URL + "/no-found", http.StatusNotFound, false},
		{origin.URL + "/no-found", http.StatusNotFound, false},
	}

	for _, test := range tests {
		res, _ := pool.HTTPClient().Get(test.origin)
		res.Body.Close()

		cached := res.Header.Get("X-From-Cache") == "1"
		if cached != test.cached {
			t.Errorf("expected a different cache hit for %s: got %t, want %t", test.origin, cached, test.cached)
		}

		if res.StatusCode != test.status {
			t.Errorf("unexpected status code for %s: got %d, want %d", test.origin, res.StatusCode, test.status)
		}
	}
}

func TestPoolHeaders(t *testing.T) {
	var got string
	want := "ForwardCacheBot/1.0"

	proxied := origin.Config.Handler
	origin.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		got = req.UserAgent()
		proxied.ServeHTTP(w, req)
	})

	req, _ := http.NewRequest("GET", origin.URL+"/small.js", nil)
	req.Header.Add("User-Agent", want)
	pool.HTTPClient().Do(req)

	if got != want {
		t.Errorf("invalid header sent to origin: got 'User-Agent: %s', want 'User-Agent: %s'", got, want)
	}

	origin.Config.Handler = proxied
}

func ExampleNewPool() {
	pool := NewPool("http://10.0.1.1:3000", httpcache.NewMemoryCache())
	pool.Set("http://10.0.1.1:3000", "http://10.0.1.2:3000")

	// -then-

	http.DefaultTransport = pool
	http.Get("https://...js/1.5.7/angular.min.js")

	// -or-

	http.DefaultClient = pool.HTTPClient()
	http.Get("https://...js/1.5.7/angular.min.js")

	// -or-

	c := pool.HTTPClient()
	c.Get("https://...js/1.5.7/angular.min.js")

	// ...

	http.ListenAndServe(":3000", pool.LocalProxy())
}

func ExampleNewClient() {
	pool := NewClient()
	pool.Set("http://10.0.1.1:3000", "http://10.0.1.2:3000")

	// -then-

	http.DefaultTransport = pool
	http.Get("https://...js/1.5.7/angular.min.js")

	// -or-

	http.DefaultClient = pool.HTTPClient()
	http.Get("https://...js/1.5.7/angular.min.js")

	// -or-

	c := pool.HTTPClient()
	c.Get("https://...js/1.5.7/angular.min.js")
}

func setup() {
	// create an origin server and a pool with 2 members
	origin = httptest.NewServer(http.FileServer(http.Dir("./test")))
	cache := httpcache.NewMemoryCache()
	localProxy = httptest.NewServer(nil)
	peerProxy = httptest.NewServer(nil)
	client := NewClient(WithPath("/fwp"), WithReplicas(100), WithHashFn(crc32.ChecksumIEEE), WithClientTransport(http.DefaultTransport))
	pool = NewPool(localProxy.URL, cache, WithClient(client), WithProxyTransport(http.DefaultTransport))
	pool.Set(localProxy.URL, peerProxy.URL)
	localProxy.Config.Handler = pool.LocalProxy()
	peerProxy.Config.Handler = newProxy("/fwp", cache, http.DefaultTransport)
}

func teardown() {
	origin.Close()
	localProxy.Close()
	peerProxy.Close()
}
