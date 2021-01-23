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

// New creates a new gateway. Config is mandatory - logger is optional and can be nil. It will be invoked once the response
// completes, either succeeding or returning an error. logger should not block.
func New(c *Config, logger func(*http.Request, *http.Response)) http.HandlerFunc {
	var l func(*http.Request, *http.Response)
	l = logger
	if l == nil {
		l = func(*http.Request, *http.Response){
			// no-op
			return
		}
	}
	return func(w http.ResponseWriter, req *http.Request) {
		b, url, ok := backend(c, req)
		if !ok {
			resp, _ := json.Marshal(c.NotFoundResponse)
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-type", "application/json")
			w.Write(resp)
			l(req, bytesToResponse(http.StatusNotFound, "application/json", resp))
			return
		}
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = c.Scheme
				r.URL.Host = b
				r.URL.Path = url
				r.Host = b
			},
			ModifyResponse: func(resp *http.Response) error {
				l(req, resp)
				return nil
			},
		}).ServeHTTP(w, req)
	}
}
