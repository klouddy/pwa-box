package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Prox struct {
	hostBasePath *string
	target       *url.URL
	proxy        *httputil.ReverseProxy
}

func NewProxy(basePath string, target string) *Prox {
	origin, _ := url.Parse(target)
	curProxy := &Prox{hostBasePath: &basePath, target: origin, proxy: httputil.NewSingleHostReverseProxy(origin)}
	curProxy.proxy.Director = func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = origin.Scheme
		r.URL.Host = origin.Host
		r.URL.Path = strings.TrimPrefix(r.URL.Path, basePath)
	}

	return curProxy
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle proxy request request", r.URL)
	// string leading

	p.proxy.ServeHTTP(w, r)
}
