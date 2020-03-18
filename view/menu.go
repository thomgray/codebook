package view

import (
	"github.com/thomgray/egg"
	"github.com/thomgray/egg/eggc"
)

type ModalMenu struct {
	*eggc.BorderView
	label *egg.View
}

func MakeModalMenu() *ModalMenu {
	mm := ModalMenu{}
	mm.BorderView = eggc.MakeBorderView()
	// mm.View = egg.MakeView()
	mm.SetVisible(false)
	mm.label = egg.MakeView()
	mm.AddSubView(mm.label)

	w := egg.WindowWidth()

	mm.SetBounds(egg.MakeBounds(3, 3, w-6, 10))
	mm.label.SetBounds(egg.MakeBounds(1, 1, w-8, 8))
	mm.label.OnDraw(func(c egg.Canvas) {
		c.DrawString2("Quit", 1, 0)
		c.DrawString("q", 7, 0, c.Foreground, c.Background, c.Attribute|egg.AttrUnderline)
		c.DrawString2("Search Paths", 1, 1)
	})

	return &mm
}

func (mm *ModalMenu) GetView() *egg.View {
	return mm.View
}
