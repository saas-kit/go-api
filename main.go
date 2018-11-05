package main

import (
	"fmt"
	"go-saas-kit/delivery/http"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	apiPort    string = "API_PORT"
	apiVersion string = "API_VERSION"
)

func init() {
	// Check whether .env config is loaded or not
	if os.Getenv("APP_NAME") != "" {
		log.Println("Environment configuration file is already loaded")
		return
	}
	// Loading environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Run web server via HTTP
	httpServer := http.NewServer(fmt.Sprintf(":%s", os.Getenv(apiPort)), fmt.Sprintf("v%s", os.Getenv(apiVersion)))
	log.Fatal(httpServer.ListenAndServe())
}
