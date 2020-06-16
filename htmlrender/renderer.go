package htmlrender

import (
	"fmt"
	"regexp"
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
	cursorX             int
	cursorY             int
	endsInWhitespace    bool
	shouldStartNewBlock bool
	didEndBlock         bool
	preformatted        bool
	strikethrough       bool
}

func (rc RenderingContext) applyPost(prc PostRenderingContext) RenderingContext {
	// should be pointer maybe?
	rc.cursorX = prc.cursorX
	rc.cursorY = prc.cursorY
	rc.endsInWhitespace = prc.endsInWhitespace
	rc.shouldStartNewBlock = prc.shouldStartNewBlock
	rc.didEndBlock = prc.didEndBlock
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
	// if rc.shouldStartNewBlock {
	// 	// otherwise there is already a block, so don't double the blocks!
	// 	// rc.Canvas.DrawString2(fmt.Sprintf("{%s}", tag), rc.cursorX, rc.cursorY)
	// 	rc.cursorY += 2
	// 	rc.shouldStartNewBlock = false
	// }
	// either way, a new block should reset the x cursor to the left margin
	rc.cursorX = rc.leftMargin
	return rc
}

type PostRenderingContext struct {
	cursorX             int
	cursorY             int
	endsInWhitespace    bool
	shouldStartNewBlock bool
	didEndBlock         bool
}

func (prc PostRenderingContext) applyBlock(rc RenderingContext, tag string) PostRenderingContext {
	if !prc.didEndBlock {
		prc.cursorY += 2
		// no need to apply a new block
	}
	prc.didEndBlock = true
	// if prc.shouldStartNewBlock {
	// 	prc.cursorY += 2
	// }
	prc.cursorX = rc.leftMargin
	prc.shouldStartNewBlock = false
	return prc
}

func (prc PostRenderingContext) noOp(rc RenderingContext) PostRenderingContext {
	prc.cursorX = rc.cursorX
	prc.cursorY = rc.cursorY
	prc.endsInWhitespace = rc.endsInWhitespace
	prc.shouldStartNewBlock = rc.shouldStartNewBlock
	prc.didEndBlock = rc.didEndBlock
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
		// not a visible tag type, so skip!
		return PostRenderingContext{}.noOp(rc)
	}

	switch tagName {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return renderHeading(n, c)
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
	case "ul", "ol", "dl":
		c = c.setLeftMargin(c.leftMargin + 2)
	case "dt":
		c.Canvas.Attribute |= egg.AttrBold
		c.Canvas.Foreground = egg.ColorGreen
	case "dd":
		c = c.setLeftMargin(c.leftMargin + 2)
	case "li":
		// this should be for either type of list!
		c.Canvas.DrawString("• ", c.leftMargin, c.cursorY, egg.ColorMagenta, c.Canvas.Background, c.Canvas.Attribute)
		c = c.setLeftMargin(c.leftMargin + 2)
	case "del":
		c.strikethrough = true
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
	prc.shouldStartNewBlock = true
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
		prc = prc.applyBlock(thisC, n.Data)
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
