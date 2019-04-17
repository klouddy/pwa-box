package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/inconshreveable/log15"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var appLogger log15.Logger

func BootstrapBox(c *Config) ProxService {
	var port = fmt.Sprintf(":%d", c.Port)

	r := mux.NewRouter()
	//perform setup on mux
	proxyService := setupReverseProxies(c, r)
	//setupReversePorxiesNew(c, r)
	setupStaticApps(c, r)

	fmt.Println("Setting metrics endpoint at ", c.Metrics.Endpoint)
	r.Handle(c.Metrics.Endpoint, promhttp.Handler())

	// start server.
	http.ListenAndServe(port, r)
	return proxyService
}

type AppLogRequest struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func setupStaticApps(config *Config, r *mux.Router) {
	appConfig := config.StaticApps
	appLogger = log15.New()

	r.HandleFunc("/loggers"+appConfig.LoggerPath, func(writer http.ResponseWriter, r *http.Request) {
		var logReq AppLogRequest
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &logReq)
		performAppLog(&logReq, appLogger)
		writer.WriteHeader(http.StatusCreated)
	})

	fs := http.FileServer(http.Dir(appConfig.Directory))
	//r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
	r.PathPrefix(appConfig.Route).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newPath := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, appConfig.Route), "/")
		if !strings.Contains(newPath, ".") {
			http.ServeFile(w, r, appConfig.Directory+"/index.html")
		} else {
			fmt.Println("equal")
			http.StripPrefix(appConfig.Route, fs).ServeHTTP(w, r)
		}
	})
	log15.Info(fmt.Sprintf("Setup static app for path %s. Serving content from: %s.  Logger path /loggers/%s", appConfig.Route, appConfig.Directory, appConfig.LoggerPath))

}

func performAppLog(request *AppLogRequest, logger log15.Logger) {
	switch request.Level {
	case "INFO":
		logger.Info(request.Message)
	case "DEBUG":
		logger.Debug(request.Message)
	case "WARN":
		logger.Warn(request.Message)
	case "ERROR":
		logger.Error(request.Message)
	default:
		logger.Info(request.Message)
	}
}

func setupReversePorxiesNew(config *Config, r *mux.Router) {

	for _, config := range config.ReverseProxies {
		r.PathPrefix(config.Route).HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
			u, _ := url.Parse(config.RemoteServer)
			revProxy := httputil.NewSingleHostReverseProxy(u)
			revProxy.ServeHTTP(writer, r)
		})
	}
}

func setupReverseProxies(config *Config, r *mux.Router) ProxService {

	proxyService := NewPoxService()

	for indx, proxyConfig := range config.ReverseProxies {
		proxyService.AddNewProxy(proxyConfig.Route, proxyConfig.RemoteServer)
		handler := http.HandlerFunc(proxyService.proxyList[indx].handle)
		r.PathPrefix(proxyConfig.Route).Handler(handler)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}

	return proxyService
}
