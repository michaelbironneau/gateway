package main

import (
	"github.com/michaelbironneau/gateway/lib"
	"log"
	"net/http"
	"os"
)

func main() {
	var configPath string
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

	http.HandleFunc("/", lib.New(c))
	log.Fatal(http.ListenAndServe(":"+c.Port, nil))
}
