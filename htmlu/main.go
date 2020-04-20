package htmlu

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"strings"

// 	"github.com/logrusorgru/aurora"
// 	"golang.org/x/net/html"
// )

// var knownTags []string = []string{
// 	"body", "p",
// 	"h1", "h2", "h3", "h4", "h5", "h6",
// }

// type NodeType int8

// type RederingContext uint8

// const (
// 	RenderNormal RederingContext = 1 << iota
// 	RenderPreformatted
// )

// type Format uint16

// const (
// 	FormatPlain Format = 1 << iota
// 	FormatStromg
// 	FormatUnderline
// 	FormatCode
// 	FormatItalic
// 	FormatStrikethough
// )

// const (
// 	NodeDiv NodeType = iota
// )

// type Node struct {
// 	Type NodeType
// }

// func main() {
// 	file := os.Args[1]
// 	fp, err1 := filepath.Abs("./html/" + file + ".html")
// 	if err1 != nil {
// 		panic(err1)
// 	}

// 	bytes, err := ioutil.ReadFile(fp)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// asHtml := blackfriday.Run(bytes)
// 	node, _ := html.Parse(strings.NewReader(string(bytes)))
// 	if node == nil {
// 		return
// 	}
// 	// fmt.Println(string(asHtml))
// 	body := getHtmlBody(node)
// 	if body == nil {
// 		return
// 	}

// 	fixWhitespace(body, true)
// 	// printDocTree(body)
// 	print(body, FormatPlain)
// }

// func parseHtml(node *html.Node) {

// }

// var blankRegex *regexp.Regexp = regexp.MustCompile(`\s+`)

// func printDocTree(node *html.Node) {
// 	var traverse func(n *html.Node, tier int)
// 	traverse = func(n *html.Node, tier int) {
// 		switch n.Type {
// 		case html.ElementNode, html.DocumentNode:
// 			tag := n.Data
// 			fmt.Printf("%s %s\n", strings.Repeat("- ", tier), tag)
// 			for c := n.FirstChild; c != nil; c = c.NextSibling {
// 				traverse(c, tier+1)
// 			}
// 		case html.TextNode:
// 			fmt.Printf("%s TXT: |%s|\n", strings.Repeat("- ", tier), n.Data)
// 		}
// 	}

// 	traverse(node, 0)
// }

// func getHtmlBody(node *html.Node) *html.Node {
// 	var traverse func(n *html.Node) *html.Node

// 	traverse = func(n *html.Node) *html.Node {
// 		switch n.Type {
// 		case html.ElementNode, html.DocumentNode:
// 			if n.Data == "body" {
// 				return n
// 			}

// 			for c := n.FirstChild; c != nil; c = c.NextSibling {
// 				bodyInHere := traverse(c)
// 				if bodyInHere != nil {
// 					return bodyInHere
// 				}
// 			}
// 		}
// 		return nil
// 	}

// 	return traverse(node)
// }

// func print(n *html.Node, format Format) (bool, bool) {
// 	var printedSomething bool
// 	var wroteNewLine bool
// 	var shouldNewLineForBlock bool = true

// 	switch n.Type {
// 	case html.ElementNode:
// 		if n.Data == "strong" {
// 			format = format | FormatStromg
// 		} else if n.Data == "code" || n.Data == "pre" {
// 			format = format | FormatCode
// 		} else if n.Data == "em" {
// 			format = format | FormatItalic
// 		} else if n.Data == "del" {
// 			format = format | FormatStrikethough
// 		}

// 		if n.Data == "hr" {
// 			fmt.Print("___________________________")
// 			printedSomething = true
// 			shouldNewLineForBlock = true
// 		} else if n.Data == "li" {
// 			fmt.Print("* ")
// 			printedSomething = true
// 			shouldNewLineForBlock = true
// 		}

// 		var childrenPrintedAnything bool = false
// 		// var childrenEndedOnnewLine bool = false
// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			// printed somethign means just that
// 			// ended block means just did a newline (so don't need to repeat it)
// 			// printing something would invalidate an ended block
// 			printed, newlined := print(c, format)

// 			childrenPrintedAnything = childrenPrintedAnything || printed // if a child printed anything, then yes
// 			if newlined {
// 				// if a new line just happened, we should note this
// 				// which can be invalidated by a subsequent print
// 				shouldNewLineForBlock = false
// 			} else if printed {
// 				// this should update if something was just printed
// 				shouldNewLineForBlock = true
// 			}
// 		}

// 		if isBlockNode(n) {
// 			// a block node should have written a new line (unless it hasn't written anything?)
// 			wroteNewLine = true
// 			if shouldNewLineForBlock {
// 				fmt.Println("\n")
// 				wroteNewLine = true
// 			}
// 			printedSomething = false
// 		}
// 		// if isBlockNode(n) && printedSomething && !childEndedBlock {
// 		// 	// fmt.Printf("\n>>PRINTED |%s - %s|<<\n", n.Data, n.FirstChild.Data)
// 		// 	fmt.Print("\n")
// 		// 	printedSomething = false
// 		// }
// 	case html.TextNode:
// 		data := n.Data
// 		if strings.Trim(data, " ") != "" {
// 			color := aurora.Color(0)
// 			if format&FormatStromg != 0 {
// 				color = color | aurora.BoldFm
// 			}
// 			if format&FormatItalic != 0 {
// 				color = color | aurora.ItalicFm
// 			}
// 			if format&FormatCode != 0 {
// 				color = color | aurora.BlackBg | aurora.WhiteFg
// 			}
// 			if format&FormatStrikethough != 0 {
// 				color = color | aurora.StrikeThroughFm
// 			}
// 			fmt.Print(aurora.Colorize(data, color))
// 			printedSomething = true
// 		}
// 		wroteNewLine = false
// 	}
// 	return printedSomething, wroteNewLine
// }

// func isBlockNode(n *html.Node) bool {
// 	switch n.Data {
// 	case "address", "article", "aside", "blockquote", "canvas", "dd", "div", "dl", "dt", "fieldset", "figcaption",
// 		"figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr",
// 		"li", "main", "nav", "noscript", "ol", "p", "pre", "section", "table", "tfoot", "ul", "video":
// 		return true
// 	default:
// 		return false
// 	}
// }

// func isEmpty(n *html.Node) {
// 	switch n.Type {
// 	case html.TextNode:

// 	}
// }

// func normaliseWhitespace(n *html.Node) {
// 	switch n.Type {
// 	case html.ElementNode:
// 		trimLeadingWS(n)
// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			normaliseWhitespace(n)
// 		}
// 	}
// }

// func trimLeadingWS(n *html.Node) bool {
// 	p := n
// 	c := n.FirstChild
// 	for c != nil {
// 		fmt.Println("c =", c.Data)
// 		if c.Type == html.TextNode {
// 			c.Data = blankRegex.ReplaceAllLiteralString(c.Data, " ")
// 			if c.Data == "" || c.Data == " " {
// 				sib := c.NextSibling
// 				p.RemoveChild(c)
// 				c = sib
// 			} else {
// 				c.Data = strings.TrimLeft(c.Data, " ")
// 				return true
// 			}
// 		} else {
// 			// not a text node, need to dig deeper into children
// 			for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
// 				stop := trimLeadingWS(cc)
// 				if stop {
// 					return stop
// 				}
// 			}
// 			c = c.NextSibling
// 		}
// 		//
// 	}

// 	return false
// }

// func fixWhitespace(n *html.Node, trimLeading bool) bool {
// 	switch n.Type {
// 	case html.TextNode:
// 		// if this node is a text node, then just normalise it
// 		// but if there is ws carried over, then trim leading
// 		// if it ends in ws, then carry that forward
// 		n.Data = blankRegex.ReplaceAllLiteralString(n.Data, " ")
// 		if trimLeading {
// 			// trim leading ws if required
// 			n.Data = strings.TrimLeft(n.Data, " ")
// 		}
// 		// if there is trailing ws, or carry forward current trimLeading if blank
// 		if n.Data == "" {
// 			return trimLeading
// 		} else {
// 			return strings.HasSuffix(n.Data, " ")
// 		}
// 	default:
// 		// this is not a text node, so need to establish if it is a block
// 		// in that case, be sure to remove leading ws at beginning
// 		if n.Data == "pre" {
// 			// preformatted, so don't fix anything,
// 			// but return needs to be true as it is block type
// 			return true
// 		}

// 		hasT := trimLeading
// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			hasT = fixWhitespace(c, hasT)
// 		}

// 		return hasT || isBlockNode(n)
// 	}
// }

// func isEmptyTextNode(n *html.Node) bool {
// 	return n.Type == html.TextNode && (n.Data == " " || n.Data == "")
// }

// func isEmptyNode() {

// }

// func textContent(n *html.Node) string {
// 	switch n.Type {
// 	case html.TextNode:
// 		return n.Data
// 	default:
// 		s := ""
// 		for c := n.FirstChild; c != nil; c = c.NextSibling {
// 			s += textContent(c)
// 		}
// 		return s
// 	}
// }
