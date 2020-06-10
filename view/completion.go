package view

import (
	"strings"

	"github.com/thomgray/egg"
)

type CompletionView struct {
	*egg.View
	completions []string
}

func MakeCompletionView() *CompletionView{
	cv :=  CompletionView{
		View: egg.MakeView(),
		completions: make([]string, 0),
	}

	cv.SetVisible(false)
	cv.OnDraw(cv.draw)
	cv.Resize()	

	egg.GetApplication().AddViewController(cv)

	return &cv
}

func (cv CompletionView) GetView() *egg.View {
	return cv.View
}

func (cv *CompletionView) SetCompletions(compl []string) {
	cv.completions = compl
	cv.Resize()
}

func (cv *CompletionView) MaxHeight() int {
	return egg.WindowHeight() - 1
}

func (cv *CompletionView) Resize() {
	w := egg.WindowWidth()
	h := len(cv.completions)+1

	newBounds := egg.MakeBounds(0, 1, w, h)
	newBounds.Height = h
	newBounds.Width = w
	cv.SetBounds(newBounds)
}

func (cv *CompletionView) draw(c egg.Canvas) {
	h := cv.MaxHeight()
	if cv.completions != nil {
		for i, compl := range cv.completions {
			if i >= h-1 {
				break
			} 
			c.DrawString2(compl, 0, i)
		}
	}
	c.DrawString2("...", 0, h-1)
	y := len(cv.completions)

	c.DrawString(strings.Repeat("─", c.Width), 0, y,  egg.ColorBlue, c.Background, c.Attribute)

}
