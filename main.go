/*
  Shhh CLI client

  Shhh application repository:   https://github.com/smallwat3r/shhh
  Shhh CLI repository:           https://github.com/smallwat3r/shhh-cli

  Author: Matthieu Petiteau <mpetiteau.pro@gmail.com>

  MIT License

  Copyright (c) 2020 Matthieu Petiteau

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in all
  copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  SOFTWARE.
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	helpSecret            = "Secret message to encrypt."
	helpEncryptPassphrase = "Passphrase to encrypt secret."
	helpDays              = "Optional, number of days to keep the secret alive."
	helpLink              = "URL link to access secret."
	helpDecryptPassphrase = "Passphrase to decrypt secret."
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	secret := createCmd.String("m", "", helpSecret)
	encryptPassphrase := createCmd.String("p", "", helpEncryptPassphrase)
	days := createCmd.Int("d", 3, helpDays)

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	link := readCmd.String("l", "", helpLink)
	decryptPassphrase := readCmd.String("p", "", helpDecryptPassphrase)

	if len(os.Args) == 1 {
		usage()
		return
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
	case "read":
		readCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if createCmd.Parsed() {
		if *secret == "" {
			fmt.Println("Please supply a secret message using -m option.")
			fmt.Println(helpSecret)
			return
		}
		if *encryptPassphrase == "" {
			fmt.Println("Please supply the passphrase using -p option.")
			fmt.Println(helpEncryptPassphrase)
			return
		}
		createSecret(*secret, *encryptPassphrase, *days)
	}

	if readCmd.Parsed() {
		if *link == "" {
			fmt.Println("Please supply the link using -l option.")
			fmt.Println(helpLink)
			return
		}
		if *decryptPassphrase == "" {
			fmt.Println("Please supply the passphrase using -p option.")
			fmt.Println(helpDecryptPassphrase)
			return
		}
		readSecret(*link, *decryptPassphrase)
	}
}

func createSecret(secret string, passphrase string, days int) {

	// Check env var for custom shhh server.
	domain := os.Getenv("SHHH_SERVER")
	if domain == "" {
		// If no custom server of shhh, defaults to standard.
		domain = "https://shhh-encrypt.com/api/c"
	} else {
		domain += "/api/c"
	}

	message := map[string]interface{}{
		"secret":     secret,
		"passphrase": passphrase,
		"days":       days,
	}

	bytesRepr, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(domain, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		log.Fatalln(err)
	}

	// Get response from server.
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result := response["response"].(map[string]interface{})

	if result["status"] == "error" {
		fmt.Println("Error:", result["details"])
		return
	}

	if result["status"] == "created" {
		fmt.Println("Secret link         :", result["link"])
		fmt.Println("One time passphrase :", passphrase)
		fmt.Println("Expires on          :", result["expires_on"])
		return
	}
}

func readSecret(link string, passphrase string) {
	l, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	host := l.Scheme + "://" + l.Host

	// Get slug in URL path.
	p := l.Path
	slug := path.Base(p)

	// Build API URL with args provided.
	api_url := host + "/api/r"

	u, err := url.Parse(api_url)
	if err != nil {
		log.Fatalln(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatalln(err)
	}

	// Add args in querystring.
	q.Add("slug", slug)
	q.Add("passphrase", passphrase)

	u.RawQuery = q.Encode()

	read_url := u.String()

	resp, err := http.Get(read_url)
	if err != nil {
		log.Fatalln(err)
	}

	// Get response from server.
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result := response["response"].(map[string]interface{})

	if result["status"] == "error" {
		fmt.Println("Error:", result["msg"])
		return
	}

	if result["status"] == "success" {
		fmt.Println(result["msg"])
		return
	}
}

func usage() {
	fmt.Println("Usage: shhh-cli <command> [<args>]")
	fmt.Println("\nCommands available:")
	fmt.Println("  create    Creates a secret message.")
	fmt.Println("  read      Read a secret message")
	fmt.Println("\nUsage of create:")
	fmt.Println("  -m string")
	fmt.Println("        " + helpSecret)
	fmt.Println("  -p string")
	fmt.Println("        " + helpEncryptPassphrase)
	fmt.Println("  -d int")
	fmt.Println("        " + helpDays + "(default 3)")
	fmt.Println("\nUsage of read:")
	fmt.Println("  -l string")
	fmt.Println("        " + helpLink)
	fmt.Println("  -p string")
	fmt.Println("        " + helpDecryptPassphrase)
	fmt.Println("\nExamples: ")
	fmt.Println("    shhh-cli create -m 'this is a secret msg.' -p SuperPassphrase123 -d 2")
	fmt.Println("    shhh-cli read -l https://shhh-encrypt.com/api/r/jKD8Uy0A9_51c8asqAYL -p SuperPassphrase123")
}
