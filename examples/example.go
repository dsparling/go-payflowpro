// Copyright 2015 Doug Sparling. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/dsparling/go-payflowpro"
)

func main() {

	data := map[string]string{
		"USER":    "MyUserId",
		"VENDOR":  "MyVendorId",
		"PARTNER": "MyPartnerId",
		"PWD":     "MyPassword",

		"TRXTYPE":       "R", // Required - transaction type (Recurring)
		"ACTION":        "I", // Required - recurring action (Inquiry)
		"TENDER":        "C", // Required - method of payment
		"ORIGPROFILEID": "12DigitProfileId",
	}

	res := payflowpro.Pfpro(data)
	if res["RESULT"] == "0" {
		fmt.Println("Success")
	} else {
		fmt.Println("Failure")
	}
}
