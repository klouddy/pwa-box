package main

import (
	"flag"
)

var proxyService *ProxService

func main() {

	//config file Location
	configFileLocation := flag.String("config", "config.json", "Location of config file.")
	flag.Parse()

	// get config setup stuff.
	var config = LoadConfiguration(*configFileLocation)
	BootstrapBox(&config)

}
