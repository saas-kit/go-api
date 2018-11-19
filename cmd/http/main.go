package main

import (
	"log"
	"saas-kit-api/pkg/server"
)

func main() {
	// Run web server via HTTP
	httpServer := server.New(config{})
	log.Fatal(httpServer.ListenAndServe())
}
