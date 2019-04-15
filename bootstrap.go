package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func BootstrapBox(c *Config) {
	var port = fmt.Sprintf(":%d", c.Port)

	r := mux.NewRouter()
	//perform setup on mux
	setupReverseProxies(c, r)
	setupStaticApps(c, r)

	fmt.Println("Setting metrics endpoint at ", c.Metrics.Endpoint)
	r.Handle(c.Metrics.Endpoint, promhttp.Handler())

	// start server.
	http.ListenAndServe(port, r)
}

func setupStaticApps(config *Config, r *mux.Router) {
	for _, appConfig := range config.StaticApps {
		fs := http.FileServer(http.Dir(appConfig.Directory))
		r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
		fmt.Printf("Setup static app for path %s. Serving content from: %s.\n", appConfig.Route, appConfig.Directory)
	}
}

func setupReverseProxies(config *Config, r *mux.Router) {

	proxyService = NewPoxService()

	for indx, proxyConfig := range config.ReverseProxies {
		proxyService.AddNewProxy(proxyConfig.Route, proxyConfig.RemoteServer)
		handler := http.HandlerFunc(proxyService.proxyList[indx].handle)
		r.PathPrefix(proxyConfig.Route).Handler(handler)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}
}
