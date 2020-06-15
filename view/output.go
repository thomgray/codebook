package view

import (
	"log"

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
	doc        *model.Document
	file       *model.File
	customDraw func(egg.Canvas)
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

func (ov *OutputView) UnbindDraw() {
	ov.View.OnDraw(ov.draw)
}

func (ov *OutputView) CustomDraw(f func(egg.Canvas)) {
	ov.View.OnDraw(f)
}

func (ov *OutputView) SetDocument(f *model.Document) {
	ov.doc = f
	ov.UnbindDraw()
}

func (ov *OutputView) SetFile(f *model.File) {
	ov.file = f
	bnds := ov.GetBounds()
	bnds.Origin.Y = 0
	ov.SetBounds(bnds)
	ov.UnbindDraw()
}

func (ov *OutputView) draw(c egg.Canvas) {
	if ov.customDraw != nil {
		ov.customDraw(c)
	} else if ov.file != nil {
		ov.drawFile(c)
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

	h := htmlrender.RenderHtml(node, c) + 1
	if ov.GetBounds().Height != h {
		newb := ov.GetBounds()
		newb.Height = h
		ov.SetBounds(newb)
		app.ReDraw()
	}
}
