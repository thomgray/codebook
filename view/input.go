package view

import (
	"github.com/thomgray/codebook/constants"
	"github.com/thomgray/egg"
	"github.com/thomgray/egg/eggc"
)

type InputView struct {
	*eggc.TextView
	label *egg.View
	mode  constants.InputMode
}

func MakeInputView(app *egg.Application) *InputView {
	tv := eggc.MakeTextView()
	label := egg.MakeView()
	iv := InputView{
		tv,
		label,
		constants.InputModeTraverse,
	}

	label.SetBounds(egg.MakeBounds(0, 0, 1, 1))
	app.AddView(label)
	label.OnDraw(func(c egg.Canvas, _ egg.State) {
		var char string
		switch iv.mode {
		case constants.InputModeTraverse:
			char = ">"
		case constants.InputModeSearch:
			char = "?"
		case constants.InputModeCommand:
			char = ":"
		}
		c.DrawString(char, 0, 0, egg.ColorCyan, c.Background, c.Attribute)
	})

	app.AddViewController(tv)
	app.SetFocusedView(tv.View)
	w := egg.WindowWidth()
	iv.SetBounds(egg.MakeBounds(2, 0, w-2, 1))

	return &iv
}

func (iv *InputView) SetMode(m constants.InputMode) {
	iv.mode = m
}

func (iv *InputView) GetMode() constants.InputMode {
	return iv.mode
}
