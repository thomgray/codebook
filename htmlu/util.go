package htmlu

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var blankRegex *regexp.Regexp = regexp.MustCompile(`\s+`)

func IsBlockNode(n *html.Node) bool {
	switch n.Data {
	case "address", "article", "aside", "blockquote", "canvas", "dd", "div", "dl", "dt", "fieldset", "figcaption",
		"figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr",
		"li", "main", "nav", "noscript", "ol", "p", "pre", "section", "table", "tfoot", "ul", "video":
		return true
	default:
		return false
	}
}

func FixWhitespace(n *html.Node, trimLeading bool) bool {
	switch n.Type {
	case html.TextNode:
		// if this node is a text node, then just normalise it
		// but if there is ws carried over, then trim leading
		// if it ends in ws, then carry that forward
		n.Data = blankRegex.ReplaceAllLiteralString(n.Data, " ")
		if trimLeading {
			// trim leading ws if required
			n.Data = strings.TrimLeft(n.Data, " ")
		}
		// if there is trailing ws, or carry forward current trimLeading if blank
		if n.Data == "" {
			return trimLeading
		} else {
			return strings.HasSuffix(n.Data, " ")
		}
	default:
		// this is not a text node, so need to establish if it is a block
		// in that case, be sure to remove leading ws at beginning
		if n.Data == "pre" {
			// preformatted, so don't fix anything,
			// but return needs to be true as it is block type
			return true
		}

		hasT := trimLeading
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			hasT = FixWhitespace(c, hasT)
		}

		return hasT || IsBlockNode(n)
	}
}

func DataToWords(n *html.Node) []string {
	return blankRegex.Split(n.Data, -1)
}

func HVal(node *html.Node) int {
	hval := 0

	switch node.Data {
	case "h1":
		hval = 1
	case "h2":
		hval = 2
	case "h3":
		hval = 3
	case "h4":
		hval = 4
	case "h5":
		hval = 5
	case "h6":
		hval = 6
	}

	return hval
}
