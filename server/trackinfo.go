package main

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
)

const (
	trackHomeUrl    = "https://www.trackwrestling.com/TWHome.jsp"
	sessionInfoLine = 92
)

type trackSession struct {
	Id  string `json:"id"`
	Tim string `json:"tim"`
}

func GetTrackSessionData() (trackSession, error) {
	resp, err := http.Get(trackHomeUrl)
	Assert(err == nil, "error getting track home")

	defer resp.Body.Close()
	Assert(resp.StatusCode == http.StatusOK, fmt.Sprintf("track home responded with status code %d", resp.StatusCode))

	scanner := bufio.NewScanner(resp.Body)

	for i := 0; i <= sessionInfoLine; i++ {
		scanner.Scan()
	}

	re := regexp.MustCompile(`TIM=(\d+)&twSessionId=([^"]+)`)
	match := re.FindStringSubmatch(scanner.Text())

	if len(match) < 3 {
		panic("could not find match for TIM and twSessionId in track home response\n" + "match: " + scanner.Text())
	}

	return trackSession{
		Tim: match[1],
		Id:  match[2],
	}, nil
}
