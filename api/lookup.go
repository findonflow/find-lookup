package handler

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bjartek/overflow/overflow"
)

// Handler responds with the IP address of the request
func Handler(w http.ResponseWriter, r *http.Request) {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	var files []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}
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
	of := overflow.NewOverflowMainnet().Config("../flow.json").Start()

	value, err := of.Script(`
import FIND from 0x09a86f2493ce2e9d

//Check the status of a fin user
pub fun main(name: String) : Address? 
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
