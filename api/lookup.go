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
	if name != "" {
		lookup(name, w)
		return
	}
	address := r.URL.Query().Get("address")
	if address != "" {
		reverse(address, w)
		return
	}
	http.Error(w, "Specify name or address query string", http.StatusInternalServerError)

}

func reverse(address string, w http.ResponseWriter) {
	of := overflow.NewOverflowMainnet().Start()

	value, err := of.Script(`
import FIND from 0x097bafa4e0b48eef
import Profile from 0x097bafa4e0b48eef

//Check the status of a fin user
pub fun main(address: Address) : String?{

	let account=getAccount(address)
	let leaseCap = account.getCapability<&FIND.LeaseCollection{FIND.LeaseCollectionPublic}>(FIND.LeasePublicPath)

	if !leaseCap.check() {
		return nil
	}

	let profile= Profile.find(address).asProfile()
	let leases = leaseCap.borrow()!.getLeaseInformation() 
	var time : UFix64?= nil
	var name :String?= nil
	for lease in leases {

		//filter out all leases that are FREE or LOCKED since they are not actice
		if lease.status != "TAKEN" {
			continue
		}

		//if we have not set a 
		if profile.findName == "" {
			if time == nil || lease.validUntil < time! {
				time=lease.validUntil
				name=lease.name
			}
		}

		if profile.findName == lease.name {
			return lease.name
		}
	}
	return name
}
`).Args(of.Arguments().String(address)).RunReturns()

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find %s error:%v", address, err), http.StatusNotFound)
		return
	}
	w.Write([]byte(value.String()))

}
func lookup(name string, w http.ResponseWriter) {
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
`).Args(of.Arguments().String(name)).RunReturns()

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
