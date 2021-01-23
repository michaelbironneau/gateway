package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func bytesToResponse(statusCode int, contentType string, buffer []byte) *http.Response {
	r := ioutil.NopCloser(bytes.NewReader(buffer))
	h := http.Header{}
	h.Set("Content-type", contentType)
	return &http.Response{
		Status:           "",
		StatusCode:       statusCode,
		Header: h,
		Body:             r,
		ContentLength:    int64(len(buffer)),
	}
}

//  writeJSON marshals a JSON object to bytes and writes as body, ignoring any errors (will just be an empty response)
func (c *Config) writeJSON(req *http.Request, w http.ResponseWriter, status int, obj interface{})  {
	resp, _ := json.Marshal(obj)
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	w.Write(resp)
	c.Interceptor(req, bytesToResponse(status, "application/json", resp))
}

// New creates a new gateway.
func New(c *Config) http.HandlerFunc {
	c.setDefaults()
	return func(w http.ResponseWriter, req *http.Request) {

		//  1. Apply Filter
		allow, filterStatus, filterBody := c.Filter(req)
		if !allow {
			c.writeJSON(req, w, filterStatus, filterBody)
			return
		}

		// 2. Find backend, or return "not found" if not found
		b, url, ok := backend(c, req)
		if !ok {
			c.writeJSON(req, w, http.StatusNotFound, c.NotFoundResponse)
			return
		}

		// 3. Reverse proxy request
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = c.Scheme
				r.URL.Host = b
				r.URL.Path = url
				r.Host = b
			},
			ModifyResponse: func(resp *http.Response) error {
				c.Interceptor(req, resp)
				return nil
			},
		}).ServeHTTP(w, req)
	}
}
