package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

const (
	eventsUrl       = "https://www.trackwrestling.com/Login.jsp"
	sessionInfoLine = 479
)

type trackEventType int

const (
	predefinedTournaments trackEventType = 1
	openTournaments       trackEventType = 2
	teamTournaments       trackEventType = 3
	freestyleTournaments  trackEventType = 4
	seasonTournaments     trackEventType = 5
)

type trackEvent struct {
	Id   string         `json:"id"`
	Name string         `json:"name"`
	Type trackEventType `json:"type"`
}

type session struct {
	Id  string `json:"id"`
	Tim string `json:"tim"`
}

func TrackEvents(w http.ResponseWriter, r *http.Request) {
	tIdxParam := r.URL.Query().Get("tournamentIndex")
	if tIdxParam == "" {
		tIdxParam = "0"
	}

	tIdx, err := strconv.Atoi(tIdxParam)
	Assert(err == nil, "error converting tournament index")

	resp, err := http.Get(fmt.Sprintf("%s?tournamentIndex=%d", eventsUrl, tIdx))
	Assert(err == nil, "error getting events")

	defer resp.Body.Close()
	Assert(resp.StatusCode == http.StatusOK, fmt.Sprintf("events responded with status code %d", resp.StatusCode))

	sessionData, err := getSessionData(resp)
	Assert(err == nil, "error getting session data")

	sessionDataJson, err := json.Marshal(sessionData)
	Assert(err == nil, "error marshaling session data")

	doc, err := RespBodyToHtml(resp)
	Assert(err == nil, "error parsing html")

	eventsList, err := getEventsList(doc, tIdx)
	Assert(err == nil, "error getting events list")

	var eventsListJson []string
	for _, event := range eventsList {
		eventJson, err := json.Marshal(event)
		Assert(err == nil, fmt.Sprintf("error marshaling json for event %s", event.Name))

		eventsListJson = append(eventsListJson, string(eventJson))
	}

	fmt.Fprintln(w, string(sessionDataJson))
	for _, event := range eventsListJson {
		fmt.Fprintln(w, event)
	}
}

func getSessionData(resp *http.Response) (session, error) {
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; i <= sessionInfoLine; i++ {
		scanner.Scan()
	}

	re := regexp.MustCompile(`TIM=(\d+)&twSessionId=([^"]+)`)
	match := re.FindStringSubmatch(scanner.Text())

	return session{
		Tim: match[1],
		Id:  match[2],
	}, nil
}

func getEventsList(n *html.Node, tIdx int) ([]trackEvent, error) {
	var eventsList []trackEvent
	currentNode := n

	for i := tIdx * 30; i < (tIdx*30)+30; i++ {
		currentNode, err := FindElementByID(currentNode, fmt.Sprintf("anchor_%d", i))
		Assert(err == nil, fmt.Sprintf("error finding anchor_%d", i))

		for _, attr := range currentNode.Attr {
			if attr.Key == "href" {
				re := regexp.MustCompile(`eventSelected\(([^,]+),([^,]+),([^,]+),([^,]+),([^)]+)\)`)
				match := re.FindStringSubmatch(attr.Val)

				eventType, err := strconv.Atoi(match[3])
				Assert(err == nil, "error converting event type")

				eventsList = append(eventsList, trackEvent{
					Id:   match[1],
					Name: match[2],
					Type: trackEventType(eventType),
				})
			}
		}
	}

	return eventsList, nil
}
