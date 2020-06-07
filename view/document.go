package view

import (
	"log"

	"github.com/thomgray/codebook/model"
	"github.com/thomgray/egg"
)

const _indent int = 4

type DocumentView struct {
	*egg.View
	doc *model.Document
}

func MakeDocumentView() *DocumentView {
	v := egg.MakeView()
	dv := DocumentView{
		View: v,
	}
	v.OnDraw(dv.draw)

	return &dv
}

func (dv *DocumentView) SetDocument(doc *model.Document) {
	dv.doc = doc
	dv.resizeForDocument()
}

func (dv *DocumentView) resizeForDocument() {
	if dv.doc == nil {
		bnds := dv.GetBounds()
		bnds.Height = 0
		dv.SetBounds(bnds)
	} else {
		childrenC := len(dv.doc.SubDocuments)
		parentC := 0
		d := dv.doc

		for d != nil {
			d = d.Super
			parentC++
		}

		log.Println("REsizieng. H = ", 1+childrenC+parentC)
		bnds := dv.GetBounds()
		bnds.Height = 1 + childrenC + parentC
		dv.SetBounds(bnds)
	}
}

func (dv *DocumentView) draw(c egg.Canvas) {
	if dv.doc == nil {
		return
	}

	var parentNames []string = make([]string, 0)
	for d := dv.doc.Super; d != nil; d = d.Super {
		name := d.SearchTerm
		parentNames = append([]string{name}, parentNames...)
	}

	y := 0
	x := 0

	pre := " / "
	for _, s := range parentNames {
		c.DrawString(pre, x, y, egg.ColorMagenta, c.Background, c.Attribute)
		c.DrawString2(s, x+3, y)
		pre = "└─╴"
		x += _indent
		y++
	}
	// ├└╴╵
	st := dv.doc.SearchTerm
	c.DrawString(pre, x, y, egg.ColorMagenta, c.Background, c.Attribute)
	c.DrawString(st, x+3, y, egg.ColorYellow, c.Background, egg.AttrBold)
	x += _indent
	y++

	pre = "├─╴"
	sdL := len(dv.doc.SubDocuments)
	for i, sv := range dv.doc.SubDocuments {
		if i == sdL-1 {
			pre = "└─╴"
		}

		childStr := sv.SearchTerm
		c.DrawString(pre, x, y, egg.ColorMagenta, c.Background, c.Attribute)
		c.DrawString2(childStr, x+3, y)
		y++
	}
}
