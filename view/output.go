package view

import (
	"log"

	"github.com/thomgray/codebook/htmlrender"
	"github.com/thomgray/codebook/model"
	"github.com/thomgray/egg"
)

type OutputView struct {
	*egg.View
	doc        *model.Document
	file       *model.File
	customDraw func(egg.Canvas)
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

func (ov *OutputView) CustomDraw(f func(egg.Canvas, egg.State)) {
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

func (ov *OutputView) draw(c egg.Canvas, s egg.State) {
	if ov.customDraw != nil {
		ov.customDraw(c)
	} else if ov.file != nil {
		ov.drawFile(c, s)
	}
}

func (ov *OutputView) drawFile(c egg.Canvas, _ egg.State) {
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
