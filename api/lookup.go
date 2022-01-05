package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/bjartek/overflow/overflow"
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
	of := overflow.NewOverflowMainnet().Start()

	value, err := of.Script(`
import FIND from 0x097bafa4e0b48eef

pub fun main(name: String) : Address?  {
    return FIND.lookupAddress(name)
}
`).RunReturns()

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find %s error:%v", name, err), http.StatusNotFound)
		return
	}
	w.Write([]byte(value.String()))
}

func isValidAddress(h string) bool {
	trimmed := strings.TrimPrefix(h, "0x")
	if len(trimmed)%2 == 1 {
		trimmed = "0" + trimmed
	}
	_, err := hex.DecodeString(trimmed)
	return err == nil
}
