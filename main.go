package main

import (
	"net/http"

	handler "github.com/findonflow/find-lookup/api"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	http.ListenAndServe(":8080", nil)
}
