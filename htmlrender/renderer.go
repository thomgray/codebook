package htmlrender

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
	"golang.org/x/net/html"
)

var blockTags []string = []string{
	"address", "article", "aside", "blockquote", "canvas", "dd", "div", "dl", "dt", "fieldset",
	"figcaption", "figure", "footer", "form", "h1", "h2", "h3", "h4", "h5", "h6", "header", "hr", "li", "main", "nav",
	"noscript", "ol", "p", "pre", "section", "table", "tfoot", "ul", "video",
}

var inlineTags []string = []string{
	"a", "abbr", "acronym", "b", "bdo", "big", "br", "button", "cite", "code", "dfn", "del", "em", "i", "img", "input", "kbd", "label",
	"map", "object", "output", "q", "samp", "select", "small", "span", "strong", "sub", "sup", "textarea", "time", "tt", "var",
}

var otherVisibleTags []string = []string{
	"html", "body", "thead", "tbody", "tr", "th", "td",
}

const strikethroughCombining rune = '̶'

/*
These were dropped from the list of tags

block:
inline: "script"

*/

type Box struct {
	leftMargin  int
	rightMargin int
	topMargin   int
}

type RenderingContext struct {
	Canvas egg.Canvas
	Box
	cursorX          int
	cursorY          int
	endsInWhitespace bool
	didEndBlock      bool
	preformatted     bool
	strikethrough    bool
	listTier         int
	listItemIndex    int
	listType         string
}

func (rc RenderingContext) applyPost(prc PostRenderingContext) RenderingContext {
	// should be pointer maybe?
	rc.cursorX = prc.cursorX
	rc.cursorY = prc.cursorY
	rc.endsInWhitespace = prc.endsInWhitespace
	rc.didEndBlock = prc.didEndBlock
	rc.listItemIndex = prc.listItemIndex
	return rc
}

func (rc RenderingContext) setLeftMargin(x int) RenderingContext {
	rc.leftMargin = x
	// if rc.cursorX < x {
	rc.cursorX = x
	// }
	return rc
}

func (rc RenderingContext) applyBlock(tag string) RenderingContext {
	if !rc.didEndBlock {
		rc.cursorY += 2
		rc.didEndBlock = true
	}
	rc.cursorX = rc.leftMargin
	return rc
}

func (rc RenderingContext) copy() RenderingContext {
	return rc
}

type PostRenderingContext struct {
	cursorX          int
	cursorY          int
	endsInWhitespace bool
	didEndBlock      bool
	listItemIndex    int
}

func (prc PostRenderingContext) applyBlock(rc RenderingContext) PostRenderingContext {
	if !prc.didEndBlock {
		prc.cursorY += 2
	}
	prc.didEndBlock = true

	prc.cursorX = rc.leftMargin
	return prc
}

func (prc PostRenderingContext) noOp(rc RenderingContext) PostRenderingContext {
	prc.cursorX = rc.cursorX
	prc.cursorY = rc.cursorY
	prc.endsInWhitespace = rc.endsInWhitespace
	prc.didEndBlock = rc.didEndBlock
	prc.listItemIndex = rc.listItemIndex
	return prc
}

func RenderHtml(node *html.Node, c egg.Canvas) int {
	rc := RenderingContext{
		Canvas: c,
		Box: Box{
			leftMargin:  0,
			rightMargin: c.Width,
			topMargin:   0,
		},
		cursorX:     0,
		cursorY:     0,
		didEndBlock: true, // initially true to prompt
	}
	pc := renderRecursive(node, rc)
	return pc.cursorY
}

func renderRecursive(n *html.Node, c RenderingContext) PostRenderingContext {
	switch n.Type {
	case html.ElementNode:
		return renderElement(n, c)
	case html.TextNode:
		return renderText(n, c)
	default:
		return renderChildren(n, c, c)
	}
}

func renderElement(n *html.Node, rc RenderingContext) PostRenderingContext {
	tagName := n.Data

	c := rc
	if elementIsBlock(tagName) {
		c = c.applyBlock(tagName)
	} else if !elementIsInline(tagName) && !elementIsOtherVisisble(tagName) {
		// not a visible tag type, so skip
		return PostRenderingContext{}.noOp(rc)
	}

	switch tagName {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return renderHeading(n, c)
	case "hr":
		return renderHr(n, c)
	// check the tag for some simple rendering rules
	case "code":
		c.Canvas.Foreground = egg.ColorWhite
		c.Canvas.Background = egg.ColorBlack
	case "pre":
		c.preformatted = true
	case "em":
		c.Canvas.Attribute |= egg.AttrUnderline
	case "strong":
		c.Canvas.Attribute |= egg.AttrBold
	case "ul", "ol":
		return renderList(n, c)
	case "dl":
		c = c.setLeftMargin(c.leftMargin + 2)
	case "dt":
		c.Canvas.Attribute |= egg.AttrBold
		c.Canvas.Foreground = egg.ColorGreen
	case "dd":
		c = c.setLeftMargin(c.leftMargin + 2)
	case "li":
		var liStr string
		switch c.listType {
		case "ol":
			liStr = strconv.Itoa(c.listItemIndex+1) + "."
		default:
			liStr = " •"
		}
		c.Canvas.DrawString(liStr, c.leftMargin, c.cursorY, egg.ColorMagenta, c.Canvas.Background, c.Canvas.Attribute)
		c = c.setLeftMargin(c.leftMargin + 3)
		c.listItemIndex++
	case "del":
		c.strikethrough = true
	case "a":
		return renderAnchor(n, c)
	}

	return renderChildren(n, c, rc)
}

// delegate priming the render context
func renderHeading(n *html.Node, rc RenderingContext) PostRenderingContext {
	thisRc := rc
	hval := 0
	switch n.Data {
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

	padW := 7 - hval
	pre := strings.Repeat("│", padW)
	underPre := "└" + strings.Repeat("┴", padW-1)

	rc = rc.setLeftMargin(rc.leftMargin + runewidth.StringWidth(pre) + 1)

	rc.Canvas.Foreground = egg.ColorRed
	rc.Canvas.Attribute |= egg.AttrBold
	prc := renderChildren(n, rc, thisRc)

	if prc.didEndBlock && prc.cursorY > 0 {
		// if child rendering applied a newline gap, back up one
		prc.cursorY--
	}

	y := prc.cursorY
	yBegin := thisRc.cursorY

	for ; yBegin < y; yBegin++ {
		rc.Canvas.DrawString(pre, thisRc.leftMargin, yBegin, egg.ColorBlue, thisRc.Canvas.Background, thisRc.Canvas.Attribute)
	}
	rc.Canvas.DrawString(underPre, thisRc.leftMargin, yBegin, egg.ColorBlue, thisRc.Canvas.Background, thisRc.Canvas.Attribute)
	underline := strings.Repeat("─", rc.Canvas.Width-thisRc.leftMargin-padW-1)
	rc.Canvas.DrawString(underline, thisRc.leftMargin+padW, yBegin, egg.ColorBlue, thisRc.Canvas.Background, thisRc.Canvas.Attribute)

	prc.cursorY += 2
	prc.cursorX = thisRc.leftMargin
	prc.didEndBlock = true
	return prc
}

func renderHr(n *html.Node, rc RenderingContext) PostRenderingContext {
	prc := PostRenderingContext{}.noOp(rc)
	line := strings.Repeat("─", rc.Canvas.Width)
	rc.Canvas.DrawString(line, 0, rc.cursorY, egg.ColorMagenta, rc.Canvas.Background, rc.Canvas.Attribute)
	prc.didEndBlock = false
	prc = prc.applyBlock(rc)

	return prc
}

func renderList(n *html.Node, rc RenderingContext) PostRenderingContext {
	listStart := 0

	if attr, err := getAttribute(n, "start"); err == nil {
		if val, err2 := strconv.Atoi(attr); err2 == nil {
			listStart = val
		}
	}

	c := rc.copy()
	c.listTier++
	c.listItemIndex = listStart
	c.listType = n.Data
	c = c.setLeftMargin(c.leftMargin + 2)

	prc := renderChildren(n, c, rc)
	// ensure that li index is whatever is was before this list was rendered
	// so that nested lists don't mess up indexing
	prc.listItemIndex = rc.listItemIndex
	return prc
}

func renderAnchor(n *html.Node, c RenderingContext) PostRenderingContext {
	if href, err := getAttribute(n, "href"); err == nil {
		nodeText, nodeTextErr := nodeText(n)
		log.Printf("href= %s", nodeText)
		if nodeTextErr == nil && nodeText == href {
			// href==text so it is a simple one
			c.Canvas.Foreground = egg.ColorBlue
			return renderChildren(n, c, c)
		}

		thisC := c.copy()

		prc := renderChildren(n, c, thisC)
		prc.cursorX++
		maxW := thisC.rightMargin - thisC.leftMargin
		remainingW := thisC.rightMargin - prc.cursorX
		// now we should draw the href, making sure to wrap lines if needed
		hrefWithBracket := fmt.Sprintf("(%s)", href)

		if remainingW < 1 {
			prc.cursorX = thisC.leftMargin
			prc.cursorY++
		}
		// draw the @...
		c.Canvas.DrawString("@", prc.cursorX, prc.cursorY, egg.ColorMagenta, c.Canvas.Background, c.Canvas.Attribute)
		prc.cursorX++

		toDraw := hrefWithBracket
		keepDrawing := true

		openingBracketX := prc.cursorX
		openingBracketY := prc.cursorY
		for keepDrawing {
			slice, remainder, done := sliceForLine(toDraw, thisC.rightMargin-prc.cursorX, maxW)
			log.Println("just keep drawing", slice)
			toDraw = remainder

			thisC.Canvas.DrawString(slice, prc.cursorX, prc.cursorY, egg.ColorBlue, c.Canvas.Background, c.Canvas.Attribute)
			if done {
				prc.cursorX += runewidth.StringWidth(slice)
			} else {
				prc.cursorX = thisC.leftMargin
				prc.cursorY++
			}
			keepDrawing = !done
		}

		// just need to tweak the bracket colour by re-drawing them...
		thisC.Canvas.DrawRune('(', openingBracketX, openingBracketY, egg.ColorMagenta, thisC.Canvas.Background, thisC.copy().Canvas.Attribute)
		thisC.Canvas.DrawRune(')', prc.cursorX-1, prc.cursorY, egg.ColorMagenta, thisC.Canvas.Background, thisC.copy().Canvas.Attribute)
		return prc
	}

	// anchor without an href? so render the content as normal then
	return renderChildren(n, c, c)
}

func renderText(n *html.Node, c RenderingContext) PostRenderingContext {
	if c.preformatted {
		return renerTextPreformatted(n, c)
	}
	normalS := normaliseText(n.Data)
	startsWithWs := strings.HasPrefix(normalS, " ")
	endsWithWs := strings.HasSuffix(normalS, " ")
	if c.endsInWhitespace && startsWithWs {
		normalS = strings.TrimLeft(normalS, " ")
	} else if c.cursorX == c.leftMargin && startsWithWs {
		normalS = strings.TrimLeft(normalS, " ")
	}
	// do this transformation afterwards
	if c.strikethrough {
		normalS = strikethroughString(normalS)
	}
	strLen := runewidth.StringWidth(normalS)
	prc := PostRenderingContext{}.noOp(c)

	if strLen == 0 {
		return prc
	}

	lineL := c.rightMargin - c.cursorX
	boxW := c.rightMargin - c.leftMargin
	if lineL < 0 {
		c.cursorY++
		c.cursorX = c.leftMargin
		lineL = boxW
	}

	keepWriting := true
	for keepWriting {
		slice, remainder, finised := sliceForLine(normalS, lineL, boxW)
		normalS = remainder
		c.Canvas.DrawString2(slice, c.cursorX, c.cursorY)
		if !finised {
			// new line
			c.cursorX = c.leftMargin
			c.cursorY++
			lineL = boxW
		} else {
			c.cursorX = c.cursorX + runewidth.StringWidth(slice)
		}
		keepWriting = !finised
	}

	prc.endsInWhitespace = endsWithWs
	prc.didEndBlock = false
	prc.cursorX = c.cursorX
	prc.cursorY = c.cursorY
	return prc
}

func renerTextPreformatted(n *html.Node, c RenderingContext) PostRenderingContext {
	s := killTabsDead(n.Data)
	boxW := c.Canvas.Width - c.leftMargin - 1
	lines := strings.Split(s, "\n")
	blankR := regexp.MustCompile("^\\s*$")
	firstLineIsBlank := blankR.Match([]byte(lines[0]))
	lastLineIsBlank := blankR.Match([]byte(lines[len(lines)-1]))

	if !firstLineIsBlank {
		pad := strings.Repeat("\000", boxW)
		c.Canvas.DrawString2(pad, c.leftMargin, c.cursorY)
		c.cursorY++
	}

	for _, l := range lines {
		padL := boxW - runewidth.StringWidth(l)
		if padL < 0 {
			padL = 0
		}
		pad := strings.Repeat("\000", padL)
		c.Canvas.DrawString2(l+pad, c.leftMargin, c.cursorY)
		c.cursorY++
	}

	if !lastLineIsBlank {
		pad := strings.Repeat("\000", boxW)
		c.Canvas.DrawString2(pad, c.leftMargin, c.cursorY)
		c.cursorY++
	}

	prc := PostRenderingContext{}.noOp(c)
	prc.didEndBlock = true
	prc.cursorY++
	return prc
}

func strikethroughString(s string) string {
	//todo - interleave with strikethrough
	var out string = ""
	for _, c := range s {
		out = fmt.Sprintf("%s%c%c", out, c, strikethroughCombining)
	}
	return out
}

func renderChildren(n *html.Node, c RenderingContext, thisC RenderingContext) PostRenderingContext {
	prc := PostRenderingContext{}.noOp(c)
	for nc := n.FirstChild; nc != nil; nc = nc.NextSibling {
		prc = renderRecursive(nc, c)
		c = c.applyPost(prc)
	}
	if elementIsBlock(n.Data) {
		prc = prc.applyBlock(thisC)
	}
	return prc
}

func elementIsBlock(tagName string) bool {
	for _, tag := range blockTags {
		if tag == tagName {
			return true
		}
	}
	return false
}

func elementIsInline(tagName string) bool {
	for _, tag := range inlineTags {
		if tag == tagName {
			return true
		}
	}
	return false
}

func elementIsOtherVisisble(tagName string) bool {
	for _, tag := range otherVisibleTags {
		if tag == tagName {
			return true
		}
	}
	return false
}

func getAttribute(node *html.Node, key string) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, nil
		}
	}
	return "", errors.New("attribute not found")
}

func nodeText(node *html.Node) (string, error) {
	if node.FirstChild == node.LastChild &&
		node.FirstChild != nil &&
		node.FirstChild.Type == html.TextNode {
		return node.FirstChild.Data, nil
	}
	return "", errors.New("node contains non-text elements")
}

func renderStringAttributed(str attString, rc RenderingContext) PostRenderingContext {
	prc := PostRenderingContext{}.noOp(rc)

	return prc
}
