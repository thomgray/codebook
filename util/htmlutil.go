package util

import (
	"log"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/russross/blackfriday.v2"
)

func MarkdownToNode(data []byte) (*html.Node, error) {
	md := blackfriday.Run(data)
	log.Print(string(md))
	node, err := html.Parse(strings.NewReader(string(md)))

	return HTMLBody(node), err
}

func HTMLBody(n *html.Node) *html.Node {
	var body *html.Node = nil
	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return body
}
