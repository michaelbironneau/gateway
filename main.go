package main

import (
	"github.com/michaelbironneau/gateway/lib"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: gateway path-to-config.json")
	}
	var path = os.Args[1]
	c, err := lib.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", lib.New(c))
	log.Fatal(http.ListenAndServe(":"+c.Port, nil))
}
