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
		fmt.Printf("Setup revers proxy at %s, redirecting to %s", proxyConfig.Route, proxyConfig.RemoteServer)
	}
	fs := http.FileServer(http.Dir(config.StaticFilesDir))

	r.PathPrefix("/").Handler(fs)
	http.ListenAndServe(port, r)
}
