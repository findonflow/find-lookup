package handler

import (
	"fmt"
	"net/http"

	"github.com/bjartek/overflow/overflow"
)

// Handler responds with the IP address of the request
func Handler(w http.ResponseWriter, r *http.Request) {
	// if only one expected
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Specify address query string", http.StatusInternalServerError)
		return
	}

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
