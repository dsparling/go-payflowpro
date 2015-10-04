# go-payflowpro
PayflowPro - Library for accessing PayPal's Payflow Pro HTTP interface

The go-payflowpro package is a simple port based on the
[PayflowPro](http://search.cpan.org/~vkhera/PayflowPro/) Perl module
by [Vivek Khera](http://search.cpan.org/~vkhera/).

## Caveats/TODO

I haven't quite yet taken care of the differences beween the original untyped language (Perl) and the typed languge Go. So for example, at the moment all results values are strings, meaning you must check for:

	if res["RESULT"] == "0" {

instead of:

	if res["RESULT"] == 0 {
	
which matches the Perl code.

## Installation

Simply install the package to your [$GOPATH](http://code.google.com/p/go-wiki/wiki/GOPATH "GOPATH") with the [go tool](http://golang.org/cmd/go/ "go command") from shell:
```bash
$ go get github.com/dsparling/go-payflowpro
```
Make sure [Git is installed](http://git-scm.com/downloads) on your machine and in your system's `PATH`.

*`go get` installs the latest tagged release*

## Examples

[Example.go](https://github.com/dsparling/go-payflowpro/blob/master/examples/example.go) is a sort of hello world for go-payflowpro and should get you started for the barebones necessities of using the package.

	cd examples
	go run example.go

## Recurring Inquiry
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
