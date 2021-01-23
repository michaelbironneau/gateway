package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// Config is a struct that holds the configuration of the gateway.
type Config struct {
	// Port that the gateway should listen on.
	Port string `json:"port"`

	//  Versions maps from the version (eg. 'v1') to path prefix to backend service, e.g. /users -> users-api.mycompany.net. The full path will be
	//  passed through to the backend service, i.e. in the above example users-api.mycompany.net will see /users as
	//  the first part of the URL.
	Versions map[string]map[string]string `json:"versions"`

	//  The response to return when an unmapped route is requested - an arbitrary JSON object that will be marshalled.
	NotFoundResponse interface{} `json:"not_found_error"`

	//  Fallback URL if the version string is not used - this is commonly the documentation page for the APIs.
	//  The fallback URL will only be used if
	//  - The `version` configuration is not blank
	//  - The URL does not start with the version
	//  - The fallback rule is not blank
	FallbackRule string `json:"fallback_rule"`

	//  Scheme is the type of URL scheme to use for requests to the backend, such as "http" or "https".
	Scheme string `json:"scheme"`

	//  Interceptor is a post-response hook to log the request/response and/or modify the response.
	//  At this time, you can only modify responses sent from a backend (not the one you specify in not_found_error or
	//  any response which is created by a Filter you have configured if using this as a library).
	Interceptor func(r *http.Request, resp *http.Response)

	//  Filter is a pre-proxy hook to allow or deny a request - `true` means the request is allowed; `false` is denied
	//  By always returning true but blocking for some time, it can also be used to implement rate limiting.
	//  The second return argument can be used to return a given status code, and the third a given response (JSON).
	Filter func(r *http.Request) (bool, int, interface{})


}

// Load loads a configuration file and parses it into a Config struct.
func Load(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return load(b)
}

func load(b []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	for k := range config.Versions {
		config.Versions[strings.ToLower(k)] = config.Versions[k] //make version string case insensitive
	}
	return &config, nil
}

func (c *Config) setDefaults(){
	if c.Interceptor == nil {
		c.Interceptor = func(*http.Request, *http.Response){return}
	}
	if c.Filter == nil {
		c.Filter = func(*http.Request)(bool, int, interface{}){return true, 0, nil}
	}
}