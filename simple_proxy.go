package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type ProxService struct {
	proxyList []Prox
	Metrics   *ProxyMetrics
}

type Prox struct {
	hostBasePath *string
	target       *url.URL
	proxy        *httputil.ReverseProxy
	service      *ProxService
}

/**
Metrics that the handle function will instrument
*/
type ProxyMetrics struct {
	Counter *prometheus.CounterVec
}

func NewPoxService() *ProxService {
	ps := ProxService{proxyList: make([]Prox, 0)}
	ps.generateMetrics()
	return &ps
}

func (ps *ProxService) AddNewProxy(basePath string, target string) {
	origin, _ := url.Parse(target)

	curProxy := &Prox{hostBasePath: &basePath,
		target: origin,
		proxy:  httputil.NewSingleHostReverseProxy(origin)}

	curProxy.proxy.Director = func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = origin.Scheme
		r.URL.Host = origin.Host
		r.URL.Path = strings.TrimPrefix(r.URL.Path, basePath)
	}
	curProxy.service = ps
	ps.proxyList = append(ps.proxyList, *curProxy)

}

/**
Wrapper for reverse proxy handlerFunc.
*/
func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	mrw := NewMetricsResponseWriter(w)
	p.proxy.ServeHTTP(mrw, r)
	statusCode := mrw.statusCode
	p.service.Metrics.Counter.With(prometheus.Labels{"target": p.target.Host, "status": strconv.Itoa(statusCode)}).Inc()
}

func (p *ProxService) generateMetrics() {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reverse_proxy_requests_total",
			Help: "A counter for configured reverse proxies.",
		},
		[]string{"target", "status"},
	)

	p.Metrics = &ProxyMetrics{Counter: counter}
	prometheus.MustRegister(p.Metrics.Counter)
}

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{w, http.StatusOK}
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}
