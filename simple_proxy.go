package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Prox struct {
	hostBasePath *string
	target       *url.URL
	proxy        *httputil.ReverseProxy
}

func NewProxy(basePath string, target string) *Prox {
	url, _ := url.Parse(target)
	return &Prox{hostBasePath: &basePath, target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle proxy request request", r.URL)
	// string leading

	p.proxy.ServeHTTP(w, r)
}
