package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	api "github.com/ioggstream/simple/api"
)

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	// Create an instance of our handler which satisfies the generated interface
	petStore := api.CreateApplication()

	// We now register our petStore above as the handler for the interface
	h := api.HandlerCustom(petStore)

	s := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
