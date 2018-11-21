package main

import (
	"log"
	"saas-kit-api/api/v1"
	"saas-kit-api/pkg/server"
)

func main() {
	// Init new HTTP server
	httpServer := server.New(config{})
	// Set up API v1 services
	// TODO: missed params
	v1.SetUp(httpServer.Router())
	// Run web server via HTTP
	log.Fatal(httpServer.ListenAndServe())
}
