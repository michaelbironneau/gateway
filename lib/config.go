package lib

import (
	"encoding/json"
	"io/ioutil"
)

// Config is a struct that holds the configuration of the gateway.
type Config struct {

	// Version is the version that should be prepended to the URL, e.g. /v1/users (in fact this could be a more general prefix)
	Version string `json:"version"`

	// Port that the gateway should listen on.
	Port string `json:"port"`

	//  Mapping from path prefix to backend service, e.g. /users -> users-api.mycompany.net. The full path will be
	//  passed through to the backend service, i.e. in the above example users-api.mycompany.net will see /users as
	//  the first part of the URL.
	Rules map[string]string `json:"rules"`

	//  The response to return when an unmapped route is requested - an arbitrary JSON object that will be marshalled.
	NotFoundResponse interface{} `json:"not_found_error"`
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
		return nil, nil
	}
	return &config, nil
}
