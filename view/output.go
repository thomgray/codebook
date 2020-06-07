package view

import (
	"log"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/net/html"

	"github.com/thomgray/codebook/htmlu"
	"github.com/thomgray/codebook/model"
	"github.com/thomgray/egg"
)

type lineTracker struct {
	doubleLineLock bool
}

func (lt *lineTracker) drew() {
	lt.doubleLineLock = false
}

func (lt *lineTracker) broke() {
	lt.doubleLineLock = true
}

const listIndent int = 2

type OutputView struct {
	*egg.View
	doc  *model.Document
	text *[]model.AttributedString
}
type contextListType uint8

const (
	contextListUl contextListType = iota
	contextListOl
)

// context passed from parent to children
type renderingContext struct {
	c                   egg.Canvas
	fg                  egg.Color
	bg                  egg.Color
	atts                egg.Attribute
	listType            contextListType
	listItemCardinality int
	preformatted        bool
	leftXMargin         int
	rightXMargin        int
	lineTracker         *lineTracker
}

// context passed back from child to parent
type renderingBackContext struct {
	x int
	y int
}

func MakeOutputView() *OutputView {
	vw := egg.MakeView()

	ov := OutputView{
		View: vw,
	}
	ov.OnDraw(ov.draw)

	return &ov
}

func (ov *OutputView) SetDocument(f *model.Document) {
	ov.doc = f
	ov.text = nil
}

func (ov *OutputView) SetText(s *[]model.AttributedString) {
	ov.doc = nil
	ov.text = s
	if ov.text != nil {
		bnds := ov.GetBounds()
		bnds.Height = len(*s)
		ov.SetBounds(bnds)
	}
}

func (ov *OutputView) draw(c egg.Canvas) {
	if ov.doc != nil {
		ov.drawFile(c)
	} else if ov.text != nil {
		ov.drawText(c)
	}
}

func (ov *OutputView) drawText(c egg.Canvas) {
	for i, l := range *ov.text {
		x := 0
		for _, seg := range l {
			c.DrawString(seg.Text, x, i, seg.Foreground, seg.Background, seg.Attributes)
			x += runewidth.StringWidth(seg.Text)
		}
	}
}

func (ov *OutputView) drawFile(c egg.Canvas) {
	if ov.doc != nil {
		h := _drawDocument(ov.doc, c)
		bnds := ov.GetBounds()
		if h != bnds.Height {
			bnds.Height = h
			ov.SetBounds(bnds)
			// need to redraw if changed size after drawing
			egg.GetApplication().ReDraw()
		}
	}
}

func _drawDocument(doc *model.Document, c egg.Canvas) int {
	y := 0
	x := 0

	rc := renderingContext{
		c:            c,
		fg:           c.Foreground,
		bg:           c.Background,
		atts:         c.Attribute,
		leftXMargin:  0,
		rightXMargin: c.Width,
		lineTracker:  &lineTracker{},
	}

	_, y = _drawNode(doc.Node, c, x, y, rc)
	return y
}

func _drawNode(n *html.Node, c egg.Canvas, x, y int, rc renderingContext) (nextX, nextY int) {
	switch n.Type {
	case html.ElementNode:
		switch n.Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return _drawHeader(n, c, x, y, rc)
		case "a":
			return _drawAnchor(n, c, x, y, rc)
		default:
			return _drawElement(n, c, x, y, rc)
		}
	case html.TextNode:
		return _drawText(n, c, x, y, rc)
	default:
		return x, y
	}
}

func _drawHeader(n *html.Node, c egg.Canvas, x, y int, rc renderingContext) (nextX, nextY int) {
	hval := htmlu.HVal(n)
	rc.fg = egg.ColorMagenta

	prebit := strings.Repeat("│", hval+1)
	thisX := x
	c.DrawString(prebit, thisX, y, egg.ColorBlue, c.Background, c.Attribute)
	x += runewidth.StringWidth(prebit) + 1

	postBit := "└" + strings.Repeat("┴", hval)
	remainder := c.Width - runewidth.StringWidth(postBit)
	allPostBitStr := postBit + strings.Repeat("─", remainder)

	c.DrawString(allPostBitStr, thisX, y+1, egg.ColorBlue, c.Background, c.Attribute)

	// todo - add left margin and prepare for multiline headings!
	x, y = _drawElementContent(n, c, x, y, rc, true)
	return x, y + 1
}

func _drawAnchor(n *html.Node, c egg.Canvas, x, y int, rc renderingContext) (nextX, nextY int) {
	href := ""
	for _, att := range n.Attr {
		if att.Key == "href" {
			href = att.Val
			break
		}
	}
	if href != "" {
		// c.DrawString(href, x, y, egg.ColorMagenta, rc.bg, rc.atts)
		rc.fg = egg.ColorBlue
		c.DrawString("[", x, y, egg.ColorMagenta, rc.bg, rc.atts)
		x++
		x, y = _drawElementContent(n, c, x, y, rc, false)
		c.DrawString(" @", x, y, egg.ColorMagenta, rc.bg, rc.atts)
		x += 2
		c.DrawString(href, x, y, egg.ColorBlue, rc.bg, rc.atts)
		x += runewidth.StringWidth(href)
		c.DrawString("]", x, y, egg.ColorMagenta, rc.bg, rc.atts)
		x++

		return x, y
	} else {
		return _drawElementContent(n, c, x, y, rc, false)
	}
}

func _drawElement(n *html.Node, c egg.Canvas, x, y int, rc renderingContext) (nextX, nextY int) {
	log.Printf("Drawing element, xmargin = %d x = %d data = %s", rc.leftXMargin, x, n.Data)
	isBlock := htmlu.IsBlockNode(n)
	if isBlock {
		x, y = _breakConditionally(x, y, rc)
	}

	switch n.Data {
	case "em":
		rc.atts = rc.atts | egg.AttrUnderline
	case "strong":
		rc.atts = rc.atts | egg.AttrBold
	case "code":
		rc.fg = egg.ColorWhite
		rc.bg = egg.ColorBlack
	case "hr":
		line := strings.Repeat("─", c.Width)
		rc.lineTracker.drew()
		c.DrawString(line, 0, y, egg.ColorYellow, c.Background, c.Attribute)
	case "pre":
		rc.preformatted = true
	case "ul":
		rc.listItemCardinality = 0
		// rc.leftXMargin += listIndent
		rc.listType = contextListUl
	case "ol":
		rc.listItemCardinality = 0
		// rc.leftXMargin += listIndent
		rc.listType = contextListOl
	case "li":
		rc.listItemCardinality++
		switch rc.listType {
		case contextListOl:
			n := strconv.Itoa(rc.listItemCardinality)
			indent := runewidth.StringWidth(n) + 2
			c.DrawString(n+".", x, y, egg.ColorCyan, rc.bg, rc.atts)
			x += indent
			rc.leftXMargin += indent
		default: //ul is default
			c.DrawString("* ", x, y, egg.ColorCyan, rc.bg, rc.atts)
			x += 2
			rc.leftXMargin += 2
			rc.lineTracker.drew()
		}
	}

	return _drawElementContent(n, c, x, y, rc, isBlock)
}

func _breakConditionally(x, y int, rc renderingContext) (int, int) {
	if rc.lineTracker.doubleLineLock {
		return x, y
	}
	log.Printf("> end of a block - should return to x margin %d", rc.leftXMargin)
	rc.lineTracker.broke()
	return rc.leftXMargin, y + 1
}

func _drawElementContent(n *html.Node, c egg.Canvas, x, y int, rc renderingContext, isBlock bool) (nextX, nextY int) {
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		x, y = _drawNode(ch, c, x, y, rc)
	}
	if isBlock && !rc.lineTracker.doubleLineLock {
		return _breakConditionally(x, y, rc)
	}
	return x, y
}

func _drawText(n *html.Node, c egg.Canvas, x, y int, rc renderingContext) (nextX, nextY int) {
	if rc.preformatted {
		fullWidth := rc.rightXMargin - rc.leftXMargin
		xx := x
		// if preformatted, just render as is!
		lines := strings.Split(n.Data, "\n")
		for _, l := range lines {
			llen := runewidth.StringWidth(l)
			pad := ""
			if llen < fullWidth {
				pad = strings.Repeat(" ", fullWidth-llen)
			}
			c.DrawString(l+pad, xx, y, rc.fg, rc.bg, rc.atts)
			y++
			xx = x
		}
		rc.lineTracker.drew()
		return xx, y - 1 //todo I think this will mess up the spacing - y might need to be decremented, and not sure about x
	}

	if strings.TrimSpace(n.Data) == "" {
		return x, y
	}

	words := htmlu.DataToWords(n)
	lastWordI := len(words) - 1
	hasTrailing := strings.HasSuffix(n.Data, " ")

	for i, w := range words {
		wlen := runewidth.StringWidth(w)
		if x+wlen > rc.rightXMargin {
			// must wrap word
			x = rc.leftXMargin
			y++
		}

		if i == lastWordI && !hasTrailing {
			c.DrawString(w, x, y, rc.fg, rc.bg, rc.atts)
			x += runewidth.StringWidth(w)
		} else {
			c.DrawString(w+" ", x, y, rc.fg, rc.bg, rc.atts)
			x += runewidth.StringWidth(w) + 1
		}
	}
	rc.lineTracker.drew()
	return x, y
}
