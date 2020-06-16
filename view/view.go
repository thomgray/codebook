package view

import (
	"github.com/thomgray/codebook/model"
	"github.com/thomgray/egg"
	"github.com/thomgray/egg/eggc"
)

// MainView ...
type MainView struct {
	OutputView *OutputView
	ScrollView *eggc.ScrollView
	activeFile *model.File
}

var app *egg.Application

// MakeMainView ...
func MakeMainView(application *egg.Application) *MainView {
	app = application
	mv := MainView{
		OutputView: MakeOutputView(),
		ScrollView: eggc.MakeScrollView(),
	}
	mv.fitToWindow()

	mv.ScrollView.AddSubView(mv.OutputView.View)
	app.AddViewController(mv.ScrollView)
	app.OnResizeEvent(func(re *egg.ResizeEvent) {
		mv.resize(re.Width, re.Height)
		app.ReDraw()
	})
	return &mv
}

func (mv *MainView) resize(w, h int) {
	mv.ScrollView.SetBounds(egg.MakeBounds(0, 2, w, h-2))
	mv.OutputView.SetBounds(egg.MakeBounds(0, 0, w, h-2))
}

func (mv *MainView) fitToWindow() {
	w, h := egg.WindowSize()
	mv.ScrollView.SetBounds(egg.MakeBounds(0, 2, w, h-2))
	mv.refit()
}

func (mv *MainView) refit() {
	bs := mv.ScrollView.GetBounds()
	outputY := 0
	mv.OutputView.SetBounds(egg.MakeBounds(0, outputY, bs.Width-1, bs.Height-outputY))
}

// func (mv *MainView) SetActiveDocument(doc *model.Document) {
// 	mv.activeDocument = doc

// 	mv.DocumentView.SetVisible(doc != nil)
// 	mv.DocumentView.SetDocument(doc)
// 	mv.OutputView.SetDocument(doc)

// 	if mv.DocumentView.IsVisible() {
// 		docBnds := mv.DocumentView.GetBounds()
// 		outBnds := mv.OutputView.GetBounds()
// 		outBnds.Y = docBnds.Height + 1
// 		mv.OutputView.SetBounds(outBnds)
// 	}
// 	mv.refit()
// }

func (mv *MainView) SetActiveFile(file *model.File) {
	mv.activeFile = file
	mv.OutputView.SetFile(file)
	mv.refit()
}

func (mv *MainView) HandleKeyEvent(e *egg.KeyEvent) {
	// switch e.Key {
	// case egg.KeyArrowUp:
	// 	if mv.OutputView.GetBounds().Y == 0 && !mv.DocumentView.IsVisible() {
	// 		mv.DocumentView.SetVisible(true)
	// 		mv.refit()
	// 		return
	// 	}
	// case egg.KeyArrowDown:
	// 	if mv.DocumentView.IsVisible() {
	// 		mv.DocumentView.SetVisible(false)
	// 		mv.refit()
	// 		return
	// 	}
	// }
	mv.ScrollView.ReceiveKeyEvent(e)
}
