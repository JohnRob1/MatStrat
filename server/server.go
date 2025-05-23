package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	fmt.Println("running server...")
	http.HandleFunc("/track-events", TrackEvents)
	http.HandleFunc("/usa-login", UsaLogin)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// NOTE: Move these functions to respective utility files once there is enough of them
func Assert(condition bool, message string) {
	if !condition {
		panic("assertion failed: " + message)
	}
}

func RespBodyToHtml(resp *http.Response) (*html.Node, error) {
	content, err := io.ReadAll(resp.Body)
	Assert(err == nil, "error reading events body")

	doc, err := html.Parse(strings.NewReader(string(content)))
	Assert(err == nil, "error parsing html")

	return doc, nil
}

func FindElementByID(n *html.Node, id string) (*html.Node, error) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == id {
				return n, nil
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, _ := FindElementByID(c, id); found != nil {
			return found, nil
		}
	}

	return nil, fmt.Errorf("Could not find element with id: %s", id)
}

func RenderInner(n *html.Node) string {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		RenderNode(&buf, c)
	}
	return buf.String()
}

func RenderNode(buf *bytes.Buffer, n *html.Node) {
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
			RenderNode(buf, c)
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
