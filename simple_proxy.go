package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	Counter        *prometheus.CounterVec
	RequestSummary *prometheus.SummaryVec
}

func NewPoxService() ProxService {
	ps := ProxService{proxyList: make([]Prox, 0)}
	ps.generateMetrics()
	return ps
}

func (ps *ProxService) AddNewProxy(basePath string, target string) {
	origin, _ := url.Parse(target)

	curProxy := &Prox{hostBasePath: &basePath,
		target: origin,
		proxy:  httputil.NewSingleHostReverseProxy(origin)}

	curProxy.proxy.Director = func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.Header.Add("Content-Type", "application/json")
		r.URL.Scheme = origin.Scheme
		r.URL.Host = origin.Host
		r.URL.Path = origin.Path + strings.TrimPrefix(r.URL.Path, *curProxy.hostBasePath)
		fmt.Println("url: ", r.URL.Host, ", ", r.URL.Path)
	}

	curProxy.proxy.ModifyResponse = func(r *http.Response) error {
		fmt.Println("response: ", r)
		return nil
	}
	curProxy.service = ps
	ps.proxyList = append(ps.proxyList, *curProxy)

}

/**
Wrapper for reverse proxy handlerFunc.
Collects metrics.
*/
func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	mrw := NewMetricsResponseWriter(w)
	now := time.Now()

	p.proxy.ServeHTTP(mrw, r)
	statusCode := mrw.statusCode
	p.service.Metrics.Counter.With(prometheus.Labels{"target": p.target.Host, "status": strconv.Itoa(statusCode), "path": r.URL.Path}).Inc()
	p.service.Metrics.RequestSummary.With(prometheus.Labels{"target": p.target.Host, "status": strconv.Itoa(statusCode), "path": r.URL.Path}).Observe(time.Since(now).Seconds())
}

/**
Setup for metrics of proxy service object
*/
func (p *ProxService) generateMetrics() {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pwa_box_reverse_proxy_requests_total",
			Help: "A counter for configured reverse proxies.",
		},
		[]string{"target", "status", "path"},
	)

	sum := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "pwa_box_reverse_proxy_requests_durations",
			Help: "Reverse proxy latencies in seconds"},
		[]string{"target", "status", "path"})

	p.Metrics = &ProxyMetrics{Counter: counter, RequestSummary: sum}
	prometheus.MustRegister(p.Metrics.Counter, p.Metrics.RequestSummary)
}

/**
gathers status code of response writer
*/
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

/**
Creates new Metrics Response writer
*/
func NewMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{w, http.StatusOK}
}

/**
Function for metricsResponseWriter.
*/
func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}
