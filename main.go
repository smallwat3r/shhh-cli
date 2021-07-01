// CLI client for Shhh

// Shhh-cli source code: https://github.com/smallwat3r/shhh-cli
// Shhh source code: https://github.com/smallwat3r/shhh

// Author: Smallwat3r - Matthieu Petiteau <mpetiteau.pro@gmail.com>

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
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/fatih/color"
)

var (
	shhhVersion = "1.3.1"
)

// Create mode
const (
	helpSecret            = "Secret message to encrypt."
	helpEncryptPassphrase = "Passphrase to encrypt secret."
	helpExpire            = "(opt) How long to keep the secret alive: 10m, 30m, 1h, 3h, 6h, 1d, 2d, 3d, 5d or 7d (default: 3d)."
	helpTries             = "(opt) Max nb of tries to open the secret: 3, 5 or 7 (default: 5)."
	helpServer            = "(opt) Shhh target server (ex: https://<server>.com)."
	helpHaveibeenpwned    = "(opt) Check passphrase against the haveibeenpwned API."
)

// Read mode
const (
	helpLink              = "URL link to access secret."
	helpDecryptPassphrase = "Passphrase to decrypt secret."
)

func version() {
	fmt.Printf("shhh-cli version %s\n\n", shhhVersion)
}

func usageCreate() string {
	h := "Usage of create:"
	h += "\n  -h, --help                 Show create help message and exit."
	h += "\n  -m, --message    <string>  " + helpSecret
	h += "\n  -p, --passphrase <string>  " + helpEncryptPassphrase
	h += "\n  -e, --expire     <int>     " + helpExpire
	h += "\n  -t, --tries      <int>     " + helpTries
	h += "\n  -h, --host       <string>  " + helpServer
	h += "\n  -s, --secure               " + helpHaveibeenpwned
	h += "\n\n  example: shhh create -m [secret] -p [passphrase] -e 30m -t 3 -s\n\n"

	return h
}

func usageRead() string {
	h := "Usage of read:"
	h += "\n  -h, --help                 Show read help message and exit."
	h += "\n  -l, --link       <string>  " + helpLink
	h += "\n  -p, --passphrase <string>  " + helpDecryptPassphrase
	h += "\n\n  example: shhh read -l [link] -p [passphrase]\n\n"

	return h
}

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

func main() {
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

	if len(os.Args) == 1 {
		version()
		usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-h", "--help":
		version()
		usage()
		os.Exit(0)
	case "-v", "--version":
		version()
		os.Exit(0)
	case "create":
		createCmd.Parse(os.Args[2:])
	case "read":
		readCmd.Parse(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "%q is not valid command\n\n", os.Args[1])
		os.Exit(127)
	}

	if createCmd.Parsed() {
		if secret == "" {
			fmt.Fprintf(
				os.Stderr,
				"Supply a message to encrypt using --message\n\n",
			)
			os.Exit(1)
		}
		if encryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Supply a passphrase using --passphrase\n\n",
			)
			os.Exit(1)
		}
		createSecret(
			secret,
			encryptPassphrase,
			expire,
			tries,
			haveibeenpwned,
			server,
		)
	}

	if readCmd.Parsed() {
		if link == "" {
			fmt.Fprintf(os.Stderr, "Supply a link using --link\n\n")
			os.Exit(1)
		}
		if decryptPassphrase == "" {
			fmt.Fprintf(
				os.Stderr,
				"Supply a passphrase using --passphrase\n\n",
			)
			os.Exit(1)
		}
		readSecret(link, decryptPassphrase)
	}
}

func createSecret(
	secret string,
	passphrase string,
	expire string,
	tries int,
	haveibeenpwned bool,
	server string,
) {
	target := getTargetServer(server)

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

	resp, err := http.Post(target, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		log.Fatalln(err)
	}

	expected := map[int]bool{200: true, 201: true, 422: true}
	if !expected[resp.StatusCode] {
		fmt.Fprintf(
			os.Stderr,
			"Failed to reach server: returned %d on %s\n\n",
			resp.StatusCode,
			target,
		)
		os.Exit(1)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result, ok := response["response"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't read response from server\n\n")
		os.Exit(1)
	}

	switch result["status"] {
	case "error":
		color.Red("%s\n\n", result["details"].(string))
		os.Exit(1)
	case "created":
		color.Green(separator(79))
		color.Green("Secret link         : %s", result["link"].(string))
		color.Green("One time passphrase : %s", passphrase)
		color.Green("Expires on          : %s", result["expires_on"].(string))
		color.Green("%s\n\n", separator(79))
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Couldn't read response from server\n\n")
		os.Exit(1)
	}
}

func readSecret(link string, passphrase string) {
	if !isUrl(link) {
		fmt.Fprintf(os.Stderr, "Shhh server link URL invalid: %s\n\n", link)
		os.Exit(1)
	}

	l, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	host := l.Scheme + "://" + l.Host

	p := l.Path
	slug := path.Base(p)

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

	u.RawQuery = q.Encode()
	readUrl := u.String()
	resp, err := http.Get(readUrl)
	if err != nil {
		log.Fatalln(err)
	}

	expected := map[int]bool{200: true, 401: true, 404: true, 422: true}
	if !expected[resp.StatusCode] {
		fmt.Fprintf(
			os.Stderr,
			"Failed to reach server: returned %d on %s\n\n",
			resp.StatusCode,
			link,
		)
		os.Exit(1)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	result, ok := response["response"].(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Couldn't read response from server\n\n")
		os.Exit(1)
	}

	switch result["status"] {
	case "error", "expired", "invalid":
		color.Red("%s\n\n", result["msg"].(string))
		os.Exit(1)
	case "success":
		color.Green(separator(79))
		color.Green(result["msg"].(string))
		color.Green("%s\n\n", separator(79))
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Couldn't read response from server\n\n")
		os.Exit(1)
	}
}

func separator(length int) string {
	sep := ""
	for i := 0; i < length; i++ {
		sep += "-"
	}
	return sep
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Get endpoint for the targeted server
func getTargetServer(server string) string {
	target := os.Getenv("SHHH_SERVER")
	if !(server == "") {
		target = server
	}
	// Default Shhh server target if none specified nor in env or params
	if target == "" {
		return "https://shhh-encrypt.herokuapp.com/api/secret"
	}
	if !isUrl(target) {
		fmt.Fprintf(
			os.Stderr,
			"Shhh server target URL invalid: %s\n\n",
			target,
		)
		os.Exit(1)
	}
	return target + "/api/secret"
}
