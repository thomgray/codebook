package view

import (
	"log"

	"github.com/mattn/go-runewidth"

	"github.com/thomgray/codebook/htmlrender"
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
	file *model.File
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

func (ov *OutputView) SetFile(f *model.File) {
	ov.file = f
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
	if ov.file != nil {
		ov.drawFile(c)
	}
	if ov.doc != nil {
		// ov.drawFile(c)
	} else if ov.text != nil {
		ov.drawText(c)
	}
}

func (ov *OutputView) drawFile(c egg.Canvas) {
	f := ov.file
	if f == nil {
		log.Println("File is null, nothing to render")
		return
	}
	node := f.Body
	if node == nil {
		log.Println("Node is null, nothing to render")
		return
	}

	htmlrender.RenderHtml(node, c)
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
