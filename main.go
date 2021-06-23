// CLI client for Shhh

// Shhh app repo: https://github.com/smallwat3r/shhh
// Shhh CLI repo: https://github.com/smallwat3r/shhh-cli

// Author: Matthieu Petiteau <mpetiteau.pro@gmail.com>

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

	"github.com/hashicorp/go-retryablehttp"
	. "github.com/logrusorgru/aurora"
)

// Parameters help information
const (
	// create mode
	helpSecret            = "Secret message to encrypt."
	helpEncryptPassphrase = "Passphrase to encrypt secret."
	helpExpire            = "(opt) How long to keep the secret alive: 10m, 30m, 1h, 3h, 6h, 1d, 2d, 3d, 5d or 7d (default: 3d)."
	helpTries             = "(opt) Max nb of tries to open the secret: 3, 5 or 7 (default: 5)."
	helpServer            = "(opt) Shhh target server (ex: https://<server>.com)."
	helpHaveibeenpwned    = "(opt) Check passphrase against the haveibeenpwned API."

	// read mode
	helpLink              = "URL link to access secret."
	helpDecryptPassphrase = "Passphrase to decrypt secret."
)

// Program version
var shhhVersion = "1.3.0"

func main() {

	// create mode
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Usage = func() {
		h := usageCreate()
		fmt.Println(h)
		os.Exit(0)
	}

	var secret string
	createCmd.StringVar(&secret, "m", "", helpSecret)
	createCmd.StringVar(&secret, "message", "", helpSecret)

	var encryptPassphrase string
	createCmd.StringVar(&encryptPassphrase, "p", "", helpEncryptPassphrase)
	createCmd.StringVar(&encryptPassphrase, "passphrase", "", helpEncryptPassphrase)

	var expire string
	createCmd.StringVar(&expire, "e", "3d", helpExpire)
	createCmd.StringVar(&expire, "expire", "3d", helpExpire)

	var tries int
	createCmd.IntVar(&tries, "t", 3, helpTries)
	createCmd.IntVar(&tries, "tries", 3, helpTries)

	var haveibeenpwned bool
	createCmd.BoolVar(&haveibeenpwned, "s", false, helpHaveibeenpwned)
	createCmd.BoolVar(&haveibeenpwned, "secure", false, helpHaveibeenpwned)

	var server string
	createCmd.StringVar(&server, "h", "", helpServer)
	createCmd.StringVar(&server, "host", "", helpServer)

	// read mode
	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readCmd.Usage = func() {
		h := usageRead()
		fmt.Println(h)
		os.Exit(0)
	}

	var link string
	readCmd.StringVar(&link, "l", "", helpLink)
	readCmd.StringVar(&link, "link", "", helpLink)

	var decryptPassphrase string
	readCmd.StringVar(&decryptPassphrase, "p", "", helpDecryptPassphrase)
	readCmd.StringVar(&decryptPassphrase, "passphrase", "", helpDecryptPassphrase)

	// Case when no mode or parameters provided
	if len(os.Args) == 1 {
		usage()
		version()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-h", "--help":
		usage()
		version()
		os.Exit(0)
	case "-v", "--version":
		version()
		os.Exit(0)
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
			os.Exit(1)
		}
		if encryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the passphrase using -p / --passphrase option.\n-p  %s\n",
				helpEncryptPassphrase,
			)
			os.Exit(1)
		}
		createSecret(secret, encryptPassphrase, expire, tries, haveibeenpwned, server)
	}

	if readCmd.Parsed() {
		if link == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the link using -l / --link option.\n-l  %s\n",
				helpLink,
			)
			os.Exit(1)
		}
		if decryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Please supply the passphrase using -p / --passphrase option.\n-p  %s\n",
				helpDecryptPassphrase,
			)
			os.Exit(1)
		}
		readSecret(link, decryptPassphrase)
	}
}

// Build a simple separator
func separator(length int) string {
	sep := ""
	for i := 0; i < length; i++ {
		sep += "-"
	}
	return sep
}

// Check URL format
func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Get targeted server
func getTargetServer(server string) string {
	target := os.Getenv("SHHH_SERVER") // if Shhh server set up in env
	if !(server == "") {
		target = server
	}
	// Default Shhh server target if none specified nor in env or params
	if target == "" {
		return "https://shhh-encrypt.herokuapp.com/api/secret"
	}
	// Check url is valid and add API endpoint
	if !isUrl(target) {
		fmt.Fprintf(os.Stderr, "Shhh server target URL invalid: %s\n", target)
		os.Exit(1)
	}
	return target + "/api/secret"
}

// Create a secret
func createSecret(secret string, passphrase string, expire string, tries int, haveibeenpwned bool, server string) {
	target := getTargetServer(server) // Get target Shhh host

	// Request
	payload := map[string]interface{}{
		"secret":         secret,
		"passphrase":     passphrase,
		"expire":         expire,
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
		fmt.Println(Red(result["details"]))
		os.Exit(1)
	case "created":
		fmt.Println(Green(separator(79)))
		fmt.Println(Green("Secret link         :"), Bold(Green(result["link"])))
		fmt.Println(Green("One time passphrase :"), Bold(Green(passphrase)))
		fmt.Println(Green("Expires on          :"), Bold(Green(result["expires_on"])))
		fmt.Println(Green(separator(79)))
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}
}

// Read a secret
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

	apiUrl := host + "/api/secret"
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
		os.Exit(1)
	case "success":
		fmt.Println(Green(separator(79)))
		fmt.Println(Bold(Green(result["msg"])))
		fmt.Println(Green(separator(79)))
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Cannot parse response from server.\n")
		os.Exit(1)
	}
}

// Print program version
func version() {
	fmt.Printf("shhh-cli version %s\n", shhhVersion)
}

// Print create mode help
func usageCreate() string {
	h := "Usage of create:"
	h += "\n  -h, --help                 Show create help message and exit."
	h += "\n  -m, --message    <string>  " + helpSecret
	h += "\n  -p, --passphrase <string>  " + helpEncryptPassphrase
	h += "\n  -e, --expire     <int>     " + helpExpire
	h += "\n  -t, --tries      <int>     " + helpTries
	h += "\n  -h, --host       <string>  " + helpServer
	h += "\n  -s, --secure               " + helpHaveibeenpwned
	h += "\n\n  example: shhh create -m [secret] -p [passphrase] -e 30m -t 3 -s\n"

	return h
}

// Print read mode help
func usageRead() string {
	h := "Usage of read:"
	h += "\n  -h, --help                 Show read help message and exit."
	h += "\n  -l, --link       <string>  " + helpLink
	h += "\n  -p, --passphrase <string>  " + helpDecryptPassphrase
	h += "\n\n  example: shhh read -l [link] -p [passphrase]\n"

	return h
}

// Print program help
func usage() {
	h := "Create or read secrets from a Shhh server."
	h += "\n\nFind more information at https://github.com/smallwat3r/shhh-cli/blob/master/README.md"

	h += "\n\nUsage:"
	h += "\n  shhh [mode] [options]"

	h += "\n\nOptions:"
	h += "\n  -h, --help     Show help message and exit."
	h += "\n  -v, --version  Show program version and exit."

	h += "\n\nModes:"
	h += "\n  create         Creates a secret message."
	h += "\n  read           Read a secret message.\n\n"

	h += usageCreate() + "\n"
	h += usageRead()

	fmt.Println(h)
}
