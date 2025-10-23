// CLI client for Shhh
//
// Original: https://github.com/smallwat3r/shhh-cli
// Service:  https://github.com/smallwat3r/shhh
//
// MIT License © 2020 Matthieu Petiteau

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fatih/color"
)

var shhhCLIVersion = "1.3.3"

const (
	// Create mode
	helpSecret            = "Secret message to encrypt."
	helpEncryptPassphrase = "Passphrase to encrypt secret."
	helpExpire            = "(opt) How long to keep the secret alive: 10m, 30m, 1h, 3h, 6h, 1d, 2d, 3d, 5d or 7d (default: 3d)."
	helpTries             = "(opt) Max nb of tries to open the secret: 3, 5 or 7 (default: 5)."
	helpServer            = "(opt) Shhh target server (ex: https://<server>.com)."

	// Read mode
	helpLink              = "URL link to access secret."
	helpDecryptPassphrase = "Passphrase to decrypt secret."
)

func version() {
	fmt.Printf("shhh-cli version %s\n\n", shhhCLIVersion)
}

func usageCreate() string {
	var b strings.Builder
	b.WriteString("Usage of create:")
	b.WriteString("\n  -h, --help                 Show create help message and exit.")
	b.WriteString("\n  -m, --message    <string>  " + helpSecret)
	b.WriteString("\n  -p, --passphrase <string>  " + helpEncryptPassphrase)
	b.WriteString("\n  -e, --expire     <string>  " + helpExpire)
	b.WriteString("\n  -t, --tries      <int>     " + helpTries)
	b.WriteString("\n  -H, --host       <string>  " + helpServer)
	b.WriteString("\n\n  example: shhh create -m [secret] -p [passphrase] -e 30m -t 3\n\n")
	return b.String()
}

func usageRead() string {
	var b strings.Builder
	b.WriteString("Usage of read:")
	b.WriteString("\n  -h, --help                 Show read help message and exit.")
	b.WriteString("\n  -l, --link       <string>  " + helpLink)
	b.WriteString("\n  -p, --passphrase <string>  " + helpDecryptPassphrase)
	b.WriteString("\n\n  example: shhh read -l [link] -p [passphrase]\n\n")
	return b.String()
}

func usage() {
	var b strings.Builder
	b.WriteString("Create or read secrets from a Shhh server.")
	b.WriteString("\n\nFind more information at https://github.com/smallwat3r/shhh-cli/blob/master/README.md")
	b.WriteString("\n\nUsage:")
	b.WriteString("\n  shhh [mode] [options]")
	b.WriteString("\n\nOptions:")
	b.WriteString("\n  -h, --help     Show help message and exit.")
	b.WriteString("\n  -v, --version  Show program version and exit.")
	b.WriteString("\n\nModes:")
	b.WriteString("\n  create         Creates a secret message.")
	b.WriteString("\n  read           Read a secret message.\n\n")
	b.WriteString(usageCreate() + "\n")
	b.WriteString(usageRead())
	fmt.Println(b.String())
}

type createRequest struct {
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
	Expire     string `json:"expire"`
	Tries      int    `json:"tries"`
}

type createAPIResponse struct {
	Response struct {
		Status    string `json:"status"`
		Details   string `json:"details,omitempty"`
		Link      string `json:"link,omitempty"`
		ExpiresOn string `json:"expires_on,omitempty"`
	} `json:"response"`
}

type readAPIResponse struct {
	Response struct {
		Status  string `json:"status"`
		Msg     string `json:"msg,omitempty"`
		Details string `json:"details,omitempty"`
	} `json:"response"`
}

var httpClient = &http.Client{
	Timeout: 12 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   8 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   8 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Usage = func() { fmt.Println(usageCreate()); os.Exit(0) }

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
	createCmd.IntVar(&tries, "t", 5, helpTries)
	createCmd.IntVar(&tries, "tries", 5, helpTries)

	var server string
	createCmd.StringVar(&server, "H", "", helpServer)
	createCmd.StringVar(&server, "host", "", helpServer)

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	readCmd.Usage = func() { fmt.Println(usageRead()); os.Exit(0) }

	var link string
	readCmd.StringVar(&link, "l", "", helpLink)
	readCmd.StringVar(&link, "link", "", helpLink)

	var decryptPassphrase string
	readCmd.StringVar(&decryptPassphrase, "p", "", helpDecryptPassphrase)
	readCmd.StringVar(&decryptPassphrase, "passphrase", "", helpDecryptPassphrase)

	// no args: show version + usage
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
		_ = createCmd.Parse(os.Args[2:])
	case "read":
		_ = readCmd.Parse(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "%q is not a valid command\n\n", os.Args[1])
		os.Exit(127)
	}

	if createCmd.Parsed() {
		if err := runCreate(secret, encryptPassphrase, expire, tries, server); err != nil {
			fail(err)
		}
	}

	if readCmd.Parsed() {
		if err := runRead(link, decryptPassphrase); err != nil {
			fail(err)
		}
	}
}

func runCreate(secret, passphrase, expire string, tries int, server string) error {
	if secret == "" {
		return errors.New("supply a message using --message")
	}
	if passphrase == "" {
		return errors.New("supply a passphrase using --passphrase")
	}
	if err := validateExpire(expire); err != nil {
		return err
	}
	if err := validateTries(tries); err != nil {
		return err
	}

	target, err := getTargetServer(server)
	if err != nil {
		return err
	}

	reqBody := createRequest{
		Secret:     secret,
		Passphrase: passphrase,
		Expire:     expire,
		Tries:      tries,
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "shhh-cli/"+shhhCLIVersion)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusUnprocessableEntity:
		// proceed
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("failed to reach server: returned %d on %s\n%s", resp.StatusCode, target, strings.TrimSpace(string(body)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read server response: %w", err)
	}

	var apiResp createAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return fmt.Errorf("couldn't decode server response: %w\nBody: %s", err, strings.TrimSpace(string(bodyBytes)))
	}

	switch apiResp.Response.Status {
	case "error":
		msg := strings.TrimSpace(apiResp.Response.Details)
		if msg == "" {
			msg = "An unexpected error occurred while creating the secret."
		}
		fmt.Fprintf(os.Stderr, "%s\n\n", msg)
		return exitCode(1)
	case "created":
		printBanner()
		color.Green("Secret link         : %s", apiResp.Response.Link)
		color.Green("One time passphrase : %s", passphrase)
		color.Green("Expires on          : %s", apiResp.Response.ExpiresOn)
		printBanner()
		return nil
	default:
		return errors.New("couldn't read response from server")
	}
}

func runRead(link string, passphrase string) error {
	if link == "" {
		return errors.New("supply a link using --link")
	}
	if passphrase == "" {
		return errors.New("supply a passphrase using --passphrase")
	}
	if !isURL(link) {
		return fmt.Errorf("Shhh server link URL invalid: %s", link)
	}

	l, err := url.Parse(link)
	if err != nil {
		return err
	}
	host := l.Scheme + "://" + l.Host
	slug := path.Base(l.Path)

	u, err := url.Parse(host + "/api/secret")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("external_id", slug)
	q.Set("passphrase", passphrase)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "shhh-cli/"+shhhCLIVersion)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusUnauthorized, http.StatusNotFound, http.StatusUnprocessableEntity:
		// proceed
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("failed to reach server: returned %d on %s\n%s", resp.StatusCode, link, strings.TrimSpace(string(body)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read server response: %w", err)
	}

	var apiResp readAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return fmt.Errorf("couldn't decode server response: %w\nBody: %s", err, strings.TrimSpace(string(bodyBytes)))
	}

	status := strings.ToLower(strings.TrimSpace(apiResp.Response.Status))
	msg := strings.TrimSpace(apiResp.Response.Msg)

	switch status {
	case "success":
		printBanner()
		color.Green(msg)
		printBanner()
		return nil
	default:
		if msg == "" {
			msg = strings.TrimSpace(apiResp.Response.Details)
		}
		color.Red("%s\n", msg)
		return exitCode(1)
	}
}

func printBanner() {
	color.Green(strings.Repeat("-", 79))
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func validateExpire(exp string) error {
	allowed := map[string]struct{}{
		"10m": {}, "30m": {}, "1h": {}, "3h": {}, "6h": {},
		"1d": {}, "2d": {}, "3d": {}, "5d": {}, "7d": {},
	}
	if _, ok := allowed[exp]; !ok {
		return fmt.Errorf("invalid expire value %q; allowed: 10m, 30m, 1h, 3h, 6h, 1d, 2d, 3d, 5d, 7d", exp)
	}
	return nil
}

func validateTries(t int) error {
	switch t {
	case 3, 5, 7:
		return nil
	default:
		return fmt.Errorf("invalid tries value %d; allowed: 3, 5, 7", t)
	}
}

func getTargetServer(serverFlag string) (string, error) {
	target := strings.TrimSpace(os.Getenv("SHHH_SERVER"))
	if serverFlag != "" {
		target = strings.TrimSpace(serverFlag)
	}
	if target == "" {
		target = "https://shhh.smallwat3r.com"
	}
	if !isURL(target) {
		return "", fmt.Errorf("Shhh server target URL invalid: %s", target)
	}
	// avoid double-appending
	if strings.HasSuffix(target, "/api/secret") {
		return target, nil
	}
	return strings.TrimRight(target, "/") + "/api/secret", nil
}

type exitCode int

func (e exitCode) Error() string { return fmt.Sprintf("exit code %d", int(e)) }

func fail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	var ec exitCode
	if errors.As(err, &ec) {
		os.Exit(int(ec))
	}
	os.Exit(1)
}
