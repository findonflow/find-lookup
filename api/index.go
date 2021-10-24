package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func respond(w http.ResponseWriter, r *http.Request, body []byte, err error) {

	w.Header().Set("Content-Type", "application/json")
	switch err {
	case nil:
		w.Write(body)
	default:
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler responds with the IP address of the request
func Handler(w http.ResponseWriter, r *http.Request) {
	reply := strings.TrimPrefix(r.URL.Path, "/")

	url := fmt.Sprintf(`https://prod-test-net-dashboard-api.azurewebsites.net/api/company/04bd44ea-0ff1-44be-a5a0-e502802c56d8/search?eventType=A.85f0d6217184009b.FIND.Register&name="%s"`, reply)

	resp, err := http.Get(url)
	if err != nil {
		respond(w, r, []byte{}, err)
	}
	defer resp.Body.Close()

	var body []Graffle
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		respond(w, r, []byte{}, err)
	}

	output := ""
	if len(body) != 0 {

		now := time.Now().Unix()
		if now <= int64(body[0].BlockEventData.ValidUntil) {
			output = body[0].BlockEventData.Owner
		}
	}

	result, err := json.Marshal(output)

	respond(w, r, result, err)
}

type Graffle struct {
	ID             string `json:"id"`
	BlockEventData struct {
		Name        string `json:"name"`
		Owner       string `json:"owner"`
		ValidUntil  int    `json:"validUntil"`
		LockedUntil int    `json:"lockedUntil"`
	} `json:"blockEventData"`
	EventDate         time.Time `json:"eventDate"`
	FlowEventID       string    `json:"flowEventId"`
	FlowTransactionID string    `json:"flowTransactionId"`
}
