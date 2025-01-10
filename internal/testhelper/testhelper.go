package testhelper

import (
	"log"
	"net/http"
	"net/http/httptest"
)

// TestHTTPServer is a http server for testing
type TestHTTPServer struct {
	srv      *httptest.Server
	handlers map[string]http.HandlerFunc
}

// NewTestHTTPServer creates a new TestHTTPServer
func NewTestHTTPServer() *TestHTTPServer {
	return &TestHTTPServer{handlers: make(map[string]http.HandlerFunc)}
}

// RegisterHandler register a http handler
func (ts *TestHTTPServer) RegisterHandler(path string, h http.HandlerFunc) {
	ts.handlers[path] = h
}

// URL returns the url of the http server
func (ts *TestHTTPServer) URL() string {
	return ts.srv.URL
}

// Start starts the http server
func (ts *TestHTTPServer) Start() {
	ts.srv = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("received request at path %q\n", r.URL.Path)

		// check auth
		if r.Header.Get("User-Agent") != "wikipedia (https://github.com/majdus/go-wikipedia/)" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h, ok := ts.handlers[r.URL.Path]
		if !ok {
			http.Error(w, "the resource path doesn't exist", http.StatusNotFound)
			return
		}
		h(w, r)
	}))

	ts.srv.Start()
}

// Stop stops the http server
func (ts *TestHTTPServer) Stop() {
	ts.srv.Close()
}
