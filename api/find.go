package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler responds with the IP address of the request
func Handler(w http.ResponseWriter, r *http.Request) {
	// if only one expected
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Specify name query string", http.StatusInternalServerError)
		return
	}

	if isValidAddress(name) {
		w.Write([]byte(name))
		return
	}

	name = strings.TrimSuffix(name, ".find")

	url := fmt.Sprintf(`https://prod-test-net-dashboard-api.azurewebsites.net/api/company/04bd44ea-0ff1-44be-a5a0-e502802c56d8/search?eventType=A.6f265aa45d8b4875.FIND.Register&name="%s"`, name)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var body []Graffle
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "s-maxage=2, stale-while-revalidate")
	output := ""
	if len(body) != 0 {

		now := time.Now().Unix()
		if now <= int64(body[0].BlockEventData.ValidUntil) {
			output = body[0].BlockEventData.Owner
			w.Write([]byte(output))
			return
		}
	}
	http.Error(w, fmt.Sprintf("Cannot find %s", name), http.StatusNotFound)
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

func isValidAddress(h string) bool {
	trimmed := strings.TrimPrefix(h, "0x")
	if len(trimmed)%2 == 1 {
		trimmed = "0" + trimmed
	}
	_, err := hex.DecodeString(trimmed)
	return err == nil
}
