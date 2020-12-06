package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

func TestParse(t *testing.T) {
	node, _ := html.Parse(strings.NewReader("<body><h1>Hello</h1></body>"))
	doc := DocumentFromNode(node, "filename")

	assert.Equal(t, "", doc.Node.Data)
}

func TestZip(t *testing.T) {
	node, _ := html.Parse(strings.NewReader("<body><h1>Hello</h1><p>content</p><h2>heading 2</h2><p>h2 content</p></body>"))

	body := node.FirstChild

	var traverse func(*html.Node)

	fmt.Println(body.Data)

	traverse = func(n *html.Node) {
		if n != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c)
			}
		}
	}

	traverse(body)

	doc := DocumentFromNode(node, "file")
	assert.Equal(t, 4, len(doc.Elements))
}
