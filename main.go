package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var proxyList []*Prox

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
	for indx, proxyConfig := range config.ReverseProxies {
		proxyList = append(proxyList, NewProxy(proxyConfig.Route, proxyConfig.RemoteServer))
		r.PathPrefix(proxyConfig.Route).HandlerFunc(proxyList[indx].handle)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}
}
