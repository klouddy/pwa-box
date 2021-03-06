package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var configFiles = []string{
	`{"reverseProxies": [],"staticApps": {"route": "/app","directory": "./static"},"metrics": {"endpoint": "/metrics"},"port": 9898,"logLevel": "WARN"}`,
	`{"reverseProxies": [],"staticApps": {"route": "/app","directory": "./static"},"metrics": {"endpoint": "/metrics"},"port": 9898,"logLevel": "BOB"}`}

func TestLoadConfigByString(t *testing.T) {
	for _, config := range configFiles {
		c := LoadConfigByString([]byte(config))
		c.Validate()
		assert.Equal(t, "WARN", c.LogLevel)
		assert.Empty(t, c.ReverseProxies)
		assert.Equal(t, "/metrics", c.Metrics.Endpoint)
	}
}

func TestConfig_ValidateLevelWithDefault(t *testing.T) {
	c := Config{LogLevel: "INFO"}
	c.ValidateLevelOrDefault()
	assert.Equal(t, "INFO", c.LogLevel)

	c.LogLevel = "BOB"
	c.ValidateLevelOrDefault()
	assert.Equal(t, "WARN", c.LogLevel)
}

func TestConfig_ValidateMetricsOrDefault(t *testing.T) {
	c := Config{}
	c.ValidateMetricsOrDefault()
	assert.Equal(t, "/metrics", c.Metrics.Endpoint)

	c = Config{Metrics: MetricsConfig{Endpoint: "/bob"}}
	c.ValidateMetricsOrDefault()
	assert.Equal(t, "/bob", c.Metrics.Endpoint)
}

func TestConfig_ValidatePortOrDefault(t *testing.T) {
	c := Config{}
	c.ValidatePortOrDefault()
	assert.Equal(t, 9898, c.Port)

	c = Config{Port: -10}
	c.ValidatePortOrDefault()
	assert.Equal(t, 9898, c.Port)

	c = Config{Port: 4200}
	c.ValidatePortOrDefault()
	assert.Equal(t, 4200, c.Port)
}
