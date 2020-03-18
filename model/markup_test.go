package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

func TestParse(t *testing.T) {
	node, _ := html.Parse(strings.NewReader("<body><h1>Hello</h1></body>"))
	doc := DocumentFromNode(node)

	assert.Equal(t, "", doc.Node.Data)
}
