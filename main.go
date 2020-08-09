// CLI client for Shhh
//
// Shhh app repo: https://github.com/smallwat3r/shhh
// Shhh CLI repo: https://github.com/smallwat3r/shhh-cli
//
// Author: Matthieu Petiteau <mpetiteau.pro@gmail.com>
//
// MIT License
//
// Copyright (c) 2020 Matthieu Petiteau
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"

	"github.com/hashicorp/go-retryablehttp"
	. "github.com/logrusorgru/aurora"
)

// Params
const (
	// create
	helpSecret            = "Secret message to encrypt."
	helpEncryptPassphrase = "Passphrase to encrypt secret."
	helpDays              = "(optional) Number of days to keep the secret alive (default: 3)."
	helpTries             = "(optional) Number of tries to open secret before it gets deleted (default: 5)."
	helpServer            = "(optional) Shhh target server (ex: https://shhh-encrypt.herokuapp.com)."
	helpHaveibeenpwned    = "(optional) If flag set, check passphrase against the haveibeenpwned API."

	// read
	helpLink              = "URL link to access secret."
	helpDecryptPassphrase = "Passphrase to decrypt secret."
)

const separator = "-------------------------------------------------------------------------------"

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Usage = func() {
		h := usageCreate()
		fmt.Println(h)
	}

	var secret string
	createCmd.StringVar(&secret, "m", "", helpSecret)
	createCmd.StringVar(&secret, "message", "", helpSecret)

	var encryptPassphrase string
	createCmd.StringVar(&encryptPassphrase, "p", "", helpEncryptPassphrase)
	createCmd.StringVar(&encryptPassphrase, "passphrase", "", helpEncryptPassphrase)

	var days int
	createCmd.IntVar(&days, "d", 3, helpDays)
	createCmd.IntVar(&days, "days", 3, helpDays)

	var tries int
	createCmd.IntVar(&tries, "t", 3, helpTries)
	createCmd.IntVar(&tries, "tries", 3, helpTries)

	var haveibeenpwned bool
	createCmd.BoolVar(&haveibeenpwned, "s", false, helpHaveibeenpwned)
	createCmd.BoolVar(&haveibeenpwned, "secure", false, helpHaveibeenpwned)

	var server string
	createCmd.StringVar(&server, "h", "", helpServer)
	createCmd.StringVar(&server, "host", "", helpServer)

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readCmd.Usage = func() {
		h := usageRead()
		fmt.Println(h)
	}

	var link string
	readCmd.StringVar(&link, "l", "", helpLink)
	readCmd.StringVar(&link, "link", "", helpLink)

	var decryptPassphrase string
	readCmd.StringVar(&decryptPassphrase, "p", "", helpDecryptPassphrase)
	readCmd.StringVar(&decryptPassphrase, "passphrase", "", helpDecryptPassphrase)

	if len(os.Args) == 1 {
		usage()
		return
	}

	switch os.Args[1] {
	case "-h", "--help":
		usage()
		return
	case "create":
		createCmd.Parse(os.Args[2:])
	case "read":
		readCmd.Parse(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "%q is not valid command.\n", os.Args[1])
		os.Exit(127)
	}

	if createCmd.Parsed() {
		if secret == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply a secret message using -m / --message option.\n-m  %s\n",
				helpSecret,
			)
			return
		}
		if encryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the passphrase using -p / --passphrase option.\n-p  %s\n",
				helpEncryptPassphrase,
			)
			return
		}
		createSecret(secret, encryptPassphrase, days, tries, haveibeenpwned, server)
	}

	if readCmd.Parsed() {
		if link == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the link using -l / --link option.\n-l  %s\n",
				helpLink,
			)
			return
		}
		if decryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the passphrase using -p / --passphrase option.\n-p  %s\n",
				helpDecryptPassphrase,
			)
			return
		}
		readSecret(link, decryptPassphrase)
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getTargetServer(server string) string {
	target := os.Getenv("SHHH_SERVER") // if Shhh server set up in env
	if !(server == "") {
		target = server
	}
	// Default Shhh server target if none specified nor in env or params
	if target == "" {
		return "https://shhh-encrypt.herokuapp.com/api/c"
	}
	// Check url is valid and add API endpoint
	if !isUrl(target) {
		fmt.Fprintf(os.Stderr, "Shhh server target URL invalid: %s\n", target)
		os.Exit(1)
	}
	return target + "/api/c"
}

func createSecret(secret string, passphrase string, days int, tries int, haveibeenpwned bool, server string) {
	target := getTargetServer(server) // Get target Shhh host

	// Request
	payload := map[string]interface{}{
		"secret":         secret,
		"passphrase":     passphrase,
		"days":           days,
		"tries":          tries,
		"haveibeenpwned": haveibeenpwned,
	}
	bytesRepr, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 4
	retryClient.Logger = nil

	resp, err := retryClient.Post(target, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure the return code is expected
	expected := map[int]bool{200: true, 201: true, 422: true}
	if !expected[resp.StatusCode] {
		fmt.Fprintf(
			os.Stderr,
			"Failed to reach Shhh: returned %d on %s\n",
			resp.StatusCode,
			target,
		)
		os.Exit(1)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result, ok := response["response"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}

	switch result["status"] {
	case "error":
		details, ok := result["details"].(map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
			os.Exit(1)
		}
		errors, ok := details["json"].(map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
			os.Exit(1)
		}
		for _, v := range errors {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v)
				fmt.Println(Red(s.Index(0)))
			}
		}
	case "created":
		fmt.Println(Green(separator))
		fmt.Println(Green("Secret link         :"), Bold(Green(result["link"])))
		fmt.Println(Green("One time passphrase :"), Bold(Green(passphrase)))
		fmt.Println(Green("Expires on          :"), Bold(Green(result["expires_on"])))
		fmt.Println(Green(separator))
	default:
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}
}

func readSecret(link string, passphrase string) {
	// Check url is valid
	if !isUrl(link) {
		fmt.Fprintf(os.Stderr, "Shhh server link URL invalid: %s\n", link)
		os.Exit(1)
	}

	// Build API endpoint from link
	l, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	host := l.Scheme + "://" + l.Host

	p := l.Path
	slug := path.Base(p) // Get unique slug from link URL

	apiUrl := host + "/api/r"
	u, err := url.Parse(apiUrl)
	if err != nil {
		log.Fatalln(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatalln(err)
	}
	q.Add("slug", slug)
	q.Add("passphrase", passphrase)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 4
	retryClient.Logger = nil

	// Request
	u.RawQuery = q.Encode()
	readUrl := u.String()
	resp, err := retryClient.Get(readUrl)
	if err != nil {
		log.Fatalln(err)
	}

	// Make sure the return code is expected
	expected := map[int]bool{200: true, 401: true, 404: true, 422: true}
	if !expected[resp.StatusCode] {
		fmt.Fprintf(
			os.Stderr,
			"Failed to reach Shhh: returned %d on %s\n",
			resp.StatusCode,
			link,
		)
		os.Exit(1)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result, ok := response["response"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}

	switch result["status"] {
	case "error", "expired", "invalid":
		fmt.Println(Red(result["msg"]))
	case "success":
		fmt.Println(Green(separator))
		fmt.Println(Bold(Green(result["msg"])))
		fmt.Println(Green(separator))
	default:
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}
}

func usageCreate() string {
	h := "Usage of create:\n"
	h += "  -h, --help                 Show create help message and exit.\n"
	h += "  -m, --message    <string>  " + helpSecret + "\n"
	h += "  -p, --passphrase <string>  " + helpEncryptPassphrase + "\n"
	h += "  -d, --days       <int>     " + helpDays + "\n"
	h += "  -t, --tries      <int>     " + helpTries + "\n"
	h += "  -h, --host       <string>  " + helpServer + "\n"
	h += "  -s, --secure               " + helpHaveibeenpwned + "\n"
	h += "  example: shhh create --message 'a secret msg' --passphrase PdVUe3bdiu --days 2 --secure\n"

	return h
}

func usageRead() string {
	h := "Usage of read:\n"
	h += "  -h, --help                 Show read help message and exit.\n"
	h += "  -l, --link       <string>  " + helpLink + "\n"
	h += "  -p, --passphrase <string>  " + helpDecryptPassphrase + "\n"
	h += "  example: shhh read --link https://shhh-encrypt.herokuapp.com/r/jKD8Uy0A9_51c8asqAYL --passphrase PdVUe3bdiu\n"

	return h
}

func usage() {
	h := "Create or read secrets from a Shhh server.\n\n"
	h += "Find more information at https://github.com/smallwat3r/shhh-cli/blob/master/README.md\n\n"

	h += "Usage:\n"
	h += "  shhh [mode] [options]\n\n"

	h += "Options:\n"
	h += "  -h, --help   Show help message and exit.\n\n"

	h += "Modes:\n"
	h += "  create       Creates a secret message.\n"
	h += "  read         Read a secret message.\n\n"

	h += usageCreate() + "\n"
	h += usageRead()

	fmt.Println(h)
}
