package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	EventsUrl       = "https://www.trackwrestling.com/Login.jsp"
	SessionInfoLine = 479
)

type EventType int

const (
	PredefinedTournaments EventType = 1
	OpenTournaments       EventType = 2
	TeamTournaments       EventType = 3
	FreestyleTournaments  EventType = 4
	SeasonTournaments     EventType = 5
)

type Event struct {
	Id   string    `json:"id"`
	Name string    `json:"name"`
	Type EventType `json:"type"`
}

type Session struct {
	Id  string `json:"id"`
	Tim string `json:"tim"`
}

func main() {
	fmt.Println("running server...")
	http.HandleFunc("/events", events)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func assert(condition bool, message string) {
	if !condition {
		panic("assertion failed: " + message)
	}
}

func events(w http.ResponseWriter, r *http.Request) {
	tIdxParam := r.URL.Query().Get("tournamentIndex")
	if tIdxParam == "" {
		tIdxParam = "0"
	}

	tIdx, err := strconv.Atoi(tIdxParam)
	assert(err == nil, "error converting tournament index")

	resp, err := http.Get(fmt.Sprintf("%s?tournamentIndex=%d", EventsUrl, tIdx))
	assert(err == nil, "error getting events")

	defer resp.Body.Close()
	assert(resp.StatusCode == http.StatusOK, fmt.Sprintf("events responded with status code %d", resp.StatusCode))

	sessionData, err := getSessionData(resp)
	assert(err == nil, "error getting session data")

	sessionDataJson, err := json.Marshal(sessionData)
	assert(err == nil, "error marshaling session data")

	content, err := io.ReadAll(resp.Body)
	assert(err == nil, "error reading events body")

	doc, err := html.Parse(strings.NewReader(string(content)))
	assert(err == nil, "error parsing html")

	eventsList, err := getEventsList(doc, tIdx)
	assert(err == nil, "error getting events list")

	var eventsListJson []string
	for _, event := range eventsList {
		eventJson, err := json.Marshal(event)
		assert(err == nil, fmt.Sprintf("error marshaling json for event %s", event.Name))

		eventsListJson = append(eventsListJson, string(eventJson))
	}

	fmt.Fprintln(w, string(sessionDataJson))
	for _, event := range eventsListJson {
		fmt.Fprintln(w, event)
	}
}

func getSessionData(resp *http.Response) (Session, error) {
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; i <= SessionInfoLine; i++ {
		scanner.Scan()
	}

	re := regexp.MustCompile(`TIM=(\d+)&twSessionId=([^"]+)`)
	match := re.FindStringSubmatch(scanner.Text())

	return Session{
		Tim: match[1],
		Id:  match[2],
	}, nil
}

func getEventsList(n *html.Node, tIdx int) ([]Event, error) {
	var eventsList []Event
	currentNode := n

	for i := tIdx * 30; i < (tIdx*30)+30; i++ {
		currentNode, err := findElementByID(currentNode, fmt.Sprintf("anchor_%d", i))
		assert(err == nil, fmt.Sprintf("error finding anchor_%d", i))

		for _, attr := range currentNode.Attr {
			if attr.Key == "href" {
				re := regexp.MustCompile(`eventSelected\(([^,]+),([^,]+),([^,]+),([^,]+),([^)]+)\)`)
				match := re.FindStringSubmatch(attr.Val)

				eventType, err := strconv.Atoi(match[3])
				assert(err == nil, "error converting event type")

				eventsList = append(eventsList, Event{
					Id:   match[1],
					Name: match[2],
					Type: EventType(eventType),
				})
			}
		}
	}

	return eventsList, nil
}

func findElementByID(n *html.Node, id string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n, nil
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, _ := findElementByID(c, id); found != nil {
			return found, nil
		}
	}

	return nil, fmt.Errorf("Could not find element with id: %s", id)
}

func renderInner(n *html.Node) string {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(&buf, c)
	}
	return buf.String()
}

func renderNode(buf *bytes.Buffer, n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		buf.WriteString("<")
		buf.WriteString(n.Data)
		for _, attr := range n.Attr {
			buf.WriteString(" ")
			buf.WriteString(attr.Key)
			buf.WriteString("=\"")
			buf.WriteString(html.EscapeString(attr.Val))
			buf.WriteString("\"")
		}
		buf.WriteString(">")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderNode(buf, c)
		}
		buf.WriteString("</")
		buf.WriteString(n.Data)
		buf.WriteString(">")
	case html.TextNode:
		buf.WriteString(html.EscapeString(n.Data))
	case html.CommentNode:
		buf.WriteString("")
	}
}
