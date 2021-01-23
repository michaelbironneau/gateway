package main

import (
	"github.com/michaelbironneau/gateway/lib"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		configPath string
		port       string
	)
	if len(os.Args) != 2 {
		if s, ok := os.LookupEnv("GATEWAY_CONFIG_FILE"); ok {
			configPath = s
		} else {
			log.Fatal("Usage: gateway path-to-config.json")
		}

	} else {
		configPath = os.Args[1]
	}

	c, err := lib.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if c.Port == "" {
		p, ok := os.LookupEnv("HTTP_PLATFORM_PORT")
		if !ok {
			log.Fatal("Config file should specify port, or the HTTP_PLATFORM_PORT environment variable must be set.")
		}
		port = p
	} else {
		port = c.Port
	}

	http.HandleFunc("/", lib.New(c, nil))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
