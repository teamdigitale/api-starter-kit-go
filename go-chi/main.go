package main

import (
	"flag"
	"log"
	"net/http"

	api "github.com/ioggstream/simple/api"
)

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "Address for test HTTP server")
	flag.Parse()

	// Create an instance of our handler which satisfies the generated interface
	petStore := api.CreateApplication()

	// We now register our petStore above as the handler for the interface
	h := api.HandlerCustom(petStore)

	s := &http.Server{
		Handler: h,
		Addr:    *addr,
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
