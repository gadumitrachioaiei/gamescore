package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gadumitrachioaiei/gamescore/service"
)

var address = flag.String("address", "", "Address for the api")

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatal("Missing address parameter, see help")
	}
	mux := http.NewServeMux()
	mux.Handle("/scores/", service.New())
	s := http.Server{
		Addr:              *address,
		Handler:           mux,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      time.Second,
		IdleTimeout:       time.Second,
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start service: %v", err)
	}
}
