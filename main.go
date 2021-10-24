package main

import (
	"log"
	"net/http"

	"github.com/findonflow/find-lookup/api"
)

func main() {
	http.HandleFunc("/", api.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
