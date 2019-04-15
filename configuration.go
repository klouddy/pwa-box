package main

import (
	"encoding/json"
	"github.com/inconshreveable/log15"
	"io/ioutil"
)

type StaticApp struct {
	Route     string `json:"route"`
	Directory string `json:"directory"`
}

type ProxyInfo struct {
	Route        string `json:"route"`
	RemoteServer string `json:"remoteServer"`
}

type MetricsConfig struct {
	Endpoint string `json:"endpoint"`
}

type Config struct {
	ReverseProxies []ProxyInfo   `json:"reverseProxies"`
	Port           int           `json:"port"`
	StaticApps     []StaticApp   `json:"staticApps"`
	Metrics        MetricsConfig `json:"metrics"`
	LogLevel       string        `json:"logLevel"`
}

var acceptableLogLevels = []string{"DEBUG", "INFO", "WARN", "ERROR"}

func (c *Config) Validate() {
	c.ValidateLevelWithDefault()
	//default metrics endpoint
	// port number exists

}

func (c *Config) ValidateLevelWithDefault() {
	var isValid bool = false
	// validate log level.
	for _, curLvl := range acceptableLogLevels {
		if curLvl == c.LogLevel {
			isValid = true
		}
	}
	if !isValid {
		c.LogLevel = "WARN"
	}
}

func LoadConfiguration(file string) Config {

	var config Config
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log15.Error("Could not read config file", err)
	}
	config = LoadConfigByString(b)
	setLogLevel(&config)
	return config

}

func LoadConfigByString(configTxt []byte) Config {
	var config Config
	err := json.Unmarshal(configTxt, &config)
	if err != nil {
		log15.Error("Error marshalling config file", err)
	}
	return config
}

/**
Default level will be warn if not configured correctly.
*/
func setLogLevel(c *Config) {
	switch c.LogLevel {
	case "WARN":
		h := log15.CallerStackHandler("%+v", log15.StdoutHandler)
		h = log15.LvlFilterHandler(log15.LvlWarn, h)
		log15.Root().SetHandler(h)
	case "ERROR":
		h := log15.CallerStackHandler("%+v", log15.StdoutHandler)
		h = log15.LvlFilterHandler(log15.LvlError, h)
		log15.Root().SetHandler(h)
	case "INFO":
		h := log15.CallerStackHandler("%+v", log15.StdoutHandler)
		h = log15.LvlFilterHandler(log15.LvlInfo, h)
		log15.Root().SetHandler(h)
	case "DEBUG":
		h := log15.CallerStackHandler("%+v", log15.StdoutHandler)
		h = log15.LvlFilterHandler(log15.LvlDebug, h)
		log15.Root().SetHandler(h)
	default:
		h := log15.CallerStackHandler("%+v", log15.StdoutHandler)
		h = log15.LvlFilterHandler(log15.LvlWarn, h)
		log15.Root().SetHandler(h)
		log15.Warn("Log level was set incorrectly in config file.  Default WARN level was set.")
		c.LogLevel = "WARN"
	}
}
