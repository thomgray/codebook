package view

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/thomgray/codebook/model"
	"github.com/thomgray/egg"
)

type OutputView struct {
	*egg.View
	file *model.File
	text *[]model.AttributedString
}

func MakeOutputView() *OutputView {
	vw := egg.MakeView()

	ov := OutputView{
		View: vw,
	}
	ov.OnDraw(ov.draw)

	return &ov
}

func (ov *OutputView) SetFile(f *model.File) {
	ov.file = f
	ov.text = nil
}

func (ov *OutputView) SetText(s *[]model.AttributedString) {
	ov.file = nil
	ov.text = s
	if ov.text != nil {
		ov.SetHeight(len(*s))
	}
}

func (ov *OutputView) draw(c egg.Canvas) {
	if ov.file != nil {
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
	if ov.file.Document != nil {
		h := _drawDocument(ov.file.Document, c)
		if h != ov.GetBounds().Height {
			ov.SetHeight(h)
			// need to redraw if changed size after drawing
			egg.GetApplication().ReDraw()
		}
	} else {
		lines := strings.Split(string(ov.file.Content), "\n")
		for i, s := range lines {
			c.DrawString2(s, 0, i)
		}
	}
}

func _drawDocument(doc *model.Document, c egg.Canvas) int {
	y := 0
	x := 0
	for _, e := range doc.Elements {
		_, y = printGeneric(e, c, x, y)
		y++
	}
	// _, y = printGeneric(doc.Elements, c, x, y)
	return y
}

func printGeneric(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	switch e.Type {
	case model.ElementTypeHeading:
		x, y = printHeading(e, c, 0, y)
	case model.ElementTypeString:
		_, y = printParagraph(e, c, 0, y)
	case model.ElementTypeCode:
		_, y = printCode(e, c, 0, y)
	case model.ElementTypeQuote:
		_, y = printQuote(e, c, 0, y)
	case model.ElementTypeUnorderedList:
		_, y = printList(e, c, 0, y)
	default:
		s := fmt.Sprintf("%s - %s", e.Tag, e.Content)
		c.DrawString2(s, 0, y)
		y++
	}
	return x, y
}

func printHeading(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	headingN := 0
	switch e.Tag {
	case "h1":
		headingN = 5
	case "h2":
		headingN = 4
	case "h3":
		headingN = 3
	case "h4":
		headingN = 2
	case "h5":
		headingN = 1
	}
	prebit := strings.Repeat("│", headingN+1)
	thisX := x
	c.DrawString(prebit, thisX, y, egg.ColorBlue, c.Background, c.Attribute)
	thisX += runewidth.StringWidth(prebit) + 1
	for _, seg := range e.Content {
		c.DrawString(seg.Raw, thisX, y, egg.ColorRed, c.Background, egg.AttrBold)
		thisX += runewidth.StringWidth(seg.Raw)
	}
	y++
	postBit := "└" + strings.Repeat("┴", headingN)
	remainder := c.Width - runewidth.StringWidth(postBit)
	allPostBitStr := postBit + strings.Repeat("─", remainder)
	c.DrawString(allPostBitStr, 0, y, egg.ColorBlue, c.Background, c.Attribute)
	return x, y + 1
}

func printParagraph(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	for _, l := range e.Content {
		s := fmt.Sprintf("%s", l.Raw)
		atts := c.Attribute
		bg := c.Background
		fg := c.Foreground

		if l.Attribution&model.AttributionBold != 0 {
			atts = atts | egg.AttrBold
		}
		if l.Attribution&model.AttributionCode != 0 {
			bg = egg.ColorBlack
			fg = egg.ColorWhite
		}
		if l.Attribution&model.AttributionEmphasis != 0 {
			atts = atts | egg.AttrUnderline
		}
		if l.Attribution&model.AttributeAnchor != 0 {
			href, ok := l.Context["href"]
			if ok {
				c.DrawString("[", x, y, egg.ColorMagenta, bg, atts)
				x++
				c.DrawString(s, x, y, egg.ColorCyan, bg, atts)
				x += runewidth.StringWidth(s)
				c.DrawString(" @ ", x, y, egg.ColorMagenta, bg, atts)
				x += 2
				c.DrawString(href, x, y, egg.ColorBlue, bg, atts)
				x += runewidth.StringWidth(href)
				c.DrawString("]", x, y, egg.ColorMagenta, bg, atts)
				x++
				continue
			}
		}
		c.DrawString(s, x, y, fg, bg, atts)
		x += runewidth.StringWidth(s)
	}
	y++
	return x, y
}

func printCode(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	// code should only have 1 plain content
	content := e.Content[0]
	lines := strings.Split(content.Raw, "\n")

	prepad := strings.Repeat("\000", c.Width)
	c.DrawString(prepad, x, y, egg.ColorWhite, egg.ColorBlack, c.Attribute)
	y++
	for _, l := range lines {
		s := fmt.Sprintf(" %s", l)
		remainingL := c.Width - runewidth.StringWidth(s)
		padding := strings.Repeat("\000", remainingL)
		c.DrawString(s+padding, x, y, egg.ColorWhite, egg.ColorBlack, egg.AttrNormal)
		y++
	}
	return x, y
}

func printQuote(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	// code should only have 1 plain content
	content := e.Content[0]
	lines := strings.Split(content.Raw, "\n")

	prepad := strings.Repeat("\000", c.Width)
	c.DrawString(prepad, x, y, egg.ColorBlack, egg.ColorWhite, c.Attribute)
	y++
	for _, l := range lines {
		s := fmt.Sprintf(" %s", l)
		remainingL := c.Width - runewidth.StringWidth(s)
		padding := strings.Repeat("\000", remainingL)
		c.DrawString(s+padding, x, y, egg.ColorBlack, egg.ColorWhite, egg.AttrNormal)
		y++
	}
	c.DrawString(prepad, x, y, egg.ColorBlack, egg.ColorWhite, c.Attribute)
	y++
	return x, y
}

func printList(e *model.Element, c egg.Canvas, x, y int) (nextX, nextY int) {
	for _, subE := range e.SubElements {
		for _, subSubE := range subE.SubElements {
			_, y = printGeneric(subSubE, c, x+3, y)
		}
	}
	return x, y
}

func printListItem(e *model.Element, c egg.Canvas, x, y int, index int) (nextX, nextY int) {
	return x, y
}
