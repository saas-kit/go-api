package main

import (
	"net/http"
	"os"
)

func main() {
	_, err := http.Get(os.Getenv("HEALTHCHECK_URL"))
	if err != nil {
		os.Exit(1)
	}
}
