package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyStripPath(t *testing.T) {

	basePathToStrip := "/foo/bar"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEqual(t, basePathToStrip, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "done")
	}))
	defer ts.Close()

	myProxy := NewProxy(basePathToStrip, ts.URL)
	localTs := httptest.NewServer(http.HandlerFunc(myProxy.handle))

	res, err := http.Get(localTs.URL + basePathToStrip)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
