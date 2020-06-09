package util

import (
	"strings"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

func MarkdownToNode(data []byte) (*html.Node, error) {
	md := blackfriday.Run(data)
	node, err := html.Parse(strings.NewReader(string(md)))

	return HTMLBody(node), err
}

func HtmlToNode(data []byte) (*html.Node, error) {
	node, err := html.Parse(strings.NewReader(string(data)))
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
