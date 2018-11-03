package main

import "go-saas-kit/delivery/http"

func main() {
	cnf := GetConfig()
	http.NewServer(cnf)
}
