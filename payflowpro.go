package payflowpro

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var debug = false
var pfprohost = ""
var version = "0.9.0"
var agent = "MailerMailer PFPro/" + version
var numretries = 3
var timeout = 30

func init() {
	Pftestmode(false) // set "live" mode as default.
	Pfdebug(false)
}

/*
 * Set test mode on or off.  Test mode means it uses the testing server
 * rather than the live one.  Default mode is live. (testmode=false)
 */
func Pftestmode(testmode bool) bool {
	if testmode {
		pfprohost = "https://pilot-payflowpro.paypal.com"
	} else {
		pfprohost = "https://payflowpro.paypal.com"
	}
	return true
}

/*
 * Set debug mode on or off.  Turns on some warn statements to track progress
 * of the request.  Default mode is off (C<$mode> == 0).
 *
 * Returns current setting.
 */
func Pfdebug(mode bool) bool {
	if mode {
		os.Setenv("HTTPS_DEBUG", "1") // assumes Crypt::SSLeay as the SSL engine
	} else {
		os.Setenv("HTTPS_DEBUG", "0") // assumes Crypt::SSLeay as the SSL engine
	}
	debug = mode
	return debug
}

func Pfpro(data map[string]string) map[string]string {

	// for the case of a referenced credit,
	// the INVNUM is not required to be set
	// so use the ORIGID instead.
	// If that's not set, just use a fixed string
	// to avoid undef warnings.
	var id string
	if value, ok := data["INVNUM"]; ok {
		id = value
	} else if value, ok := data["ORIGID"]; ok {
		id = value
	} else {
		id = "NOID"
	}

	requestId := strconv.FormatInt(time.Now().Unix(), 10) + data["TRXTYPE"] + id
	if len(requestId) >= 32 {
		requestId = requestId[:32]
	}

	if value, ok := data["TIMEOUT"]; ok {
		timeout, _ = strconv.Atoi(value)
	}

	if debug {
		fmt.Printf("%#v", data)
	}
	content := ""
	for key, value := range data {
		content += key + "[" + strconv.Itoa(len(value)) + "]=" + value + "&"
	}
	content += "VERBOSITY[6]=MEDIUM"

	client := &http.Client{}

	var query = []byte(content)
	req, err := http.NewRequest("POST", pfprohost, bytes.NewBuffer(query))
	if err != nil {
		fmt.Println("Request Error")
	}
	req.Header.Set("Content-Type", "text/namevalue")
	req.Header.Set("Content-Length", strconv.Itoa(len(content)))
	req.Header.Set("User-Agent", agent)
	req.Header.Set("Connection", "close")
	req.Header.Set("Host", pfprohost)

	// X-Headers
	req.Header.Add("X-VPS-REQUEST-ID", requestId)
	req.Header.Add("X-VPS-CLIENT-TIMEOUT", strconv.Itoa(timeout))
	req.Header.Add("X-VPS-VIT-INTEGRATION-PRODUCT", agent)
	req.Header.Add("X-VPS-VIT-INTEGRATION-VERSION", version)
	req.Header.Add("X-VPS-VIT-OS-NAME", runtime.GOOS)
	//req.Header.Add("X-VPS-VIT-OS-VERSION", $Config::Config{osvers})
	req.Header.Add("X-VPS-VIT-RUNTIME-VERSION", runtime.Version())

	if debug {
		fmt.Printf("HTTP Request:\n\n%#v", req)
		debugHttp(httputil.DumpRequestOut(req, true))
	}

	retval := make(map[string]string) // hash of values to return
	maxtries := numretries
	var response *http.Response

	//  Keep trying the request until we succeed, or fail NUMRETRIES times.
	// Since the REQUEST_ID is the same, we don't ever process
	// the request more than once, but we deal with timout cases:
	// If the request worked and we failed to get the response, we just
	// get the original response back; if it failed to reach PayPal, we
	// just retry it.  NOTE: This does not retry on payflow errors, just
	// when the HTTP protocol has failures/errors such as timeout.
	isSuccess := false
	for value := 0; ; {
		if debug {
			fmt.Printf("Running request, %d left\n", maxtries)
		}

		// delay for a bit between failures
		time.Sleep(time.Duration((numretries-maxtries)*30) * time.Second)

		response, err = client.Do(req)
		defer response.Body.Close()
		if err != nil {
			panic(err)
		}
		value = response.StatusCode
		if value == 200 || maxtries == 0 {
			isSuccess = true
			break
		}
		maxtries--
	}

	if isSuccess {
		// parse the return value into the hash and send it back.
		if debug {
			debugHttp(httputil.DumpResponse(response, true))
		}

		htmlData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		responseSlice := strings.Split(string(htmlData), "&")
		for _, part := range responseSlice {
			pair := strings.Split(part, "=")
			retval[pair[0]] = pair[1]
		}

	} else {
		// some error. fake up the old API's error code so existing code continues
		// to work.  this should just cause a retry on the application.
		if debug {
			fmt.Printf("HTTP communication error: %s\n", response.Status)
		}
		retval["RESULT"] = "-1"
		retval["RESPMSG"] = "Failed to connect to host"
	}
	retval["X-VPS-REQUEST-ID"] = requestId // useful for debugging
	return retval
}

func debugHttp(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
