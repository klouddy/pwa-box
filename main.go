package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var proxyService *ProxService

func main() {

	//config file Location
	configFileLocation := flag.String("config", "config.json", "Location of config file.")
	flag.Parse()

	// get config setup stuff.
	var config = LoadConfiguration(*configFileLocation)
	var port = fmt.Sprintf(":%d", config.Port)

	r := mux.NewRouter()
	//perform setup on mux
	setupReverseProxies(config, r)
	setupStaticApps(config, r)

	fmt.Println("Setting metrics endpoint at ", config.Metrics.Endpoint)
	r.Handle(config.Metrics.Endpoint, promhttp.Handler())

	// start server.
	http.ListenAndServe(port, r)
}

func setupStaticApps(config Config, r *mux.Router) {
	for _, appConfig := range config.StaticApps {
		fs := http.FileServer(http.Dir(appConfig.Directory))
		r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
		fmt.Printf("Setup static app for path %s. Serving content from: %s.\n", appConfig.Route, appConfig.Directory)
	}
}

func setupReverseProxies(config Config, r *mux.Router) {

	proxyService = NewPoxService()

	for indx, proxyConfig := range config.ReverseProxies {
		proxyService.AddNewProxy(proxyConfig.Route, proxyConfig.RemoteServer)
		handler := http.HandlerFunc(proxyService.proxyList[indx].handle)
		r.PathPrefix(proxyConfig.Route).Handler(handler)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}
}
