package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var configFiles = []string{
	`{"reverseProxies": [],"staticApps": [{"route": "/app","directory": "./static"}],"metrics": {"endpoint": "/metrics"},"port": 9898,"logLevel": "WARN"}`,
	`{"reverseProxies": [],"staticApps": [{"route": "/app","directory": "./static"}],"metrics": {"endpoint": "/metrics"},"port": 9898,"logLevel": "BOB"}`}

func TestLoadConfigByString(t *testing.T) {
	for _, config := range configFiles {
		c := LoadConfigByString([]byte(config))
		c.Validate()
		assert.Equal(t, "WARN", c.LogLevel)
		assert.Empty(t, c.ReverseProxies)
		assert.Equal(t, len(c.StaticApps), 1)
		assert.Equal(t, "/metrics", c.Metrics.Endpoint)
	}
}
