package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Port           int    `json:"port"`
	StaticFilesDir string `json:"staticFilesDir"`
}

func main() {
	fmt.Println("Hello world.")
	config := LoadConfiguration("config.json")
	port := fmt.Sprintf(":%d", config.Port)
	fmt.Println("port", port)
	proxy := NewProxy("http://localhost:3000")
	http.HandleFunc("/proxyServer", ProxyServer)

	http.HandleFunc("/", proxy.handle)
	fs := http.FileServer(http.Dir(config.StaticFilesDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(port, nil)
}

func ProxyServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in proxy server")
	w.Write([]byte("Reverse proxy Server Running. "))
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config

}
