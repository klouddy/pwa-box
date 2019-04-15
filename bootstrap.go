package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/inconshreveable/log15"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

var loggers map[string]log15.Logger = make(map[string]log15.Logger)

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

type AppLogRequest struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func setupStaticApps(config *Config, r *mux.Router) {
	for _, appConfig := range config.StaticApps {
		loggers[appConfig.Name] = log15.New("staticAppName", appConfig.Name)
		fs := http.FileServer(http.Dir(appConfig.Directory))
		r.PathPrefix(appConfig.Route).Handler(http.StripPrefix(appConfig.Route, fs))
		r.HandleFunc("/loggers"+appConfig.LoggerPath, func(writer http.ResponseWriter, r *http.Request) {
			var logReq AppLogRequest
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &logReq)
			performAppLog(&logReq, loggers[appConfig.Name])
			writer.WriteHeader(http.StatusCreated)
		})
		log15.Info(fmt.Sprintf("Setup static app for path %s. Serving content from: %s.  Logger path /loggers/%s", appConfig.Route, appConfig.Directory, appConfig.LoggerPath))
	}
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

func setupReverseProxies(config *Config, r *mux.Router) {

	proxyService = NewPoxService()

	for indx, proxyConfig := range config.ReverseProxies {
		proxyService.AddNewProxy(proxyConfig.Route, proxyConfig.RemoteServer)
		handler := http.HandlerFunc(proxyService.proxyList[indx].handle)
		r.PathPrefix(proxyConfig.Route).Handler(handler)
		fmt.Printf("Setup revers proxy at %s, redirecting to %s \n", proxyConfig.Route, proxyConfig.RemoteServer)
	}
}
