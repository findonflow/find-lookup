package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetRouter returns the router for the API
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler).Methods(http.MethodGet)
	return r
}

func respond(w http.ResponseWriter, r *http.Request, body []byte, err error) {

	w.Header().Set("Content-Type", "application/json")
	switch err {
	case nil:
		w.Write(body)
	default:
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

// Handler responds with the IP address of the request
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)

	reply := "foobar"
	body, err := json.Marshal(reply)
	respond(w, r, body, err)
}
