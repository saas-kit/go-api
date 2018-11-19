package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// This option can increase application performance via decrease number of the GC starts
	// See full documentation to detail https://golang.org/pkg/runtime/debug/#SetGCPercent
	// debug.SetGCPercent(10000)

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
