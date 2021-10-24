package main

import (
	"log"
	"net/http"

	handler "github.com/findonflow/find-lookup/api"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
