package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var proxyList []*Prox

func main() {
	r := mux.NewRouter()
	var config = LoadConfiguration("config.json")
	var port = fmt.Sprintf(":%d", config.Port)
	for indx, proxyConfig := range config.ReverseProxies {
		proxyList = append(proxyList, NewProxy(proxyConfig.Route, proxyConfig.RemoteServer))
		r.PathPrefix(proxyConfig.Route).HandlerFunc(proxyList[indx].handle)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}

	for _, appConfig := range config.StaticApps {
		fs := http.FileServer(http.Dir(appConfig.Directory))
		r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
		fmt.Printf("Setup static app for path %s. Serving content from: %s.\n", appConfig.Route, appConfig.Directory)
	}

	http.ListenAndServe(port, r)
}
