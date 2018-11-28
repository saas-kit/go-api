package main

import (
	"log"
	"saas-kit-api/app/address"
	"saas-kit-api/pkg/config"
	"saas-kit-api/pkg/database"
	"saas-kit-api/pkg/server"
)

func main() {
	// Init config
	cnf := config.New()

	// Create db connection
	db, err := database.New(cnf)
	if err != nil {
		log.Fatal(err)
	}
	defer checkErr(db.Close())

	// Init new HTTP server
	httpServer := server.New(cnf)
	router := httpServer.Router()

	// Set up address microservice
	checkErr(address.SetUp(cnf, router, db))

	// Run web server via HTTP
	log.Fatal(httpServer.ListenAndServe())
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
