package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var proxyList []*Prox

func main() {

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reverse_proxy_requests_total",
			Help: "A counter for configured reverse proxies.",
		},
		[]string{"target"},
	)

	proxyMetrics := ProxyMetrics{Counter: counter}

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(proxyMetrics.Counter)

	//config file Location
	configFileLocation := flag.String("config", "config.json", "Location of config file.")
	flag.Parse()

	// get config setup stuff.
	var config = LoadConfiguration(*configFileLocation)
	var port = fmt.Sprintf(":%d", config.Port)

	r := mux.NewRouter()
	//perform setup on mux
	setupReverseProxies(config, r, proxyMetrics)
	setupStaticApps(config, r)

	r.Handle("/metrics", promhttp.Handler())
	// start server.
	http.ListenAndServe(port, r)
}

func setupMetrics() {

}

func setupStaticApps(config Config, r *mux.Router) {
	for _, appConfig := range config.StaticApps {
		fs := http.FileServer(http.Dir(appConfig.Directory))
		r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
		fmt.Printf("Setup static app for path %s. Serving content from: %s.\n", appConfig.Route, appConfig.Directory)
	}
}

func setupReverseProxies(config Config, r *mux.Router, metrics ProxyMetrics) {
	for indx, proxyConfig := range config.ReverseProxies {
		proxyList = append(proxyList, NewProxy(proxyConfig.Route, proxyConfig.RemoteServer, metrics))
		handler := http.HandlerFunc(proxyList[indx].handle)
		r.PathPrefix(proxyConfig.Route).HandlerFunc(handler)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}
}
