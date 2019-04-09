package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProxyInfo struct {
	Route        string `json:"route"`
	RemoteServer string `json:"remoteServer"`
}

type Config struct {
	ReverseProxies []ProxyInfo `json:"reverseProxies"`
	Port           int         `json:"port"`
	StaticFilesDir string      `json:"staticFilesDir"`
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
