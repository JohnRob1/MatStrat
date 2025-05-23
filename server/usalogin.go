package main

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
)

const (
	loginUrl   = "https://www.usabracketing.com/login"
	tokenLine  = 42
	tokenRegex = `value="([^"]+)"`
)

func UsaLogin(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(loginUrl)
	Assert(err == nil, "error getting usa-login")

	defer resp.Body.Close()

	token, err := getToken(resp)
	fmt.Fprintln(w, token)
}

func getToken(resp *http.Response) (string, error) {
	scanner := bufio.NewScanner(resp.Body)

	for range tokenLine {
		scanner.Scan()
	}

	re := regexp.MustCompile(tokenRegex)
	match := re.FindStringSubmatch(scanner.Text())

	return match[1], nil
}
