package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
)

func backend(c *Config, r *http.Request) (string, string, bool) {
	ps := strings.SplitN(r.URL.Path, "/", 3)
	if len(ps) != 3 {
		return tryFallback(c, r)
	}
	rules, ok := c.Versions[strings.ToLower(ps[1])]
	if !ok {
		return tryFallback(c, r)
	}
	pathToMatch := "/" + ps[2]
	for k, v := range rules {
		if strings.Index(pathToMatch, k) == 0 {
			return v, pathToMatch, true
		}
	}
	return "", "", false
}

func tryFallback(c *Config, r *http.Request) (string, string, bool){
	if c.FallbackRule != "" {
		return c.FallbackRule, r.URL.Path, true
	}
	return "", "", false
}

// New creates a new gateway.
func New(c *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		b, url, ok := backend(c, req)
		if !ok {
			resp, _ := json.Marshal(c.NotFoundResponse)
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-type", "application/json")
			w.Write(resp)
			return
		}
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = c.Scheme
				r.URL.Host = b
				r.URL.Path = url
				r.Host = b
			},
		}).ServeHTTP(w, req)
	}
}
