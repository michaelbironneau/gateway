package lib

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func backend(c *Config, r *http.Request) (string, string, bool) {
	var (
		pathToMatch string
	)

	if c.Version != "" {
		ps := strings.SplitN(r.URL.Path, "/", 3)
		if  len(ps) != 3 {
			return "", "", false //expect URL of form /{version}/
		}
		pathToMatch = "/" + ps[2]
	} else {
		pathToMatch = r.URL.Path
	}
	for k, v := range c.Rules {
		if strings.Index(pathToMatch, k) == 0 {
			return v, pathToMatch, true
		}
	}
	return "", "", false
}

// New creates a new gateway.
func New(c *Config) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		b, url, ok := backend(c, req)
		if !ok {
			resp, _ := json.Marshal(c.NotFoundResponse)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(resp)
			return
		}
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = "http"
				r.URL.Host = b
				r.URL.Path = url
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}
