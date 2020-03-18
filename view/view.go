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
	w, h := egg.WindowSize()
	mv.Resize(w, h)

	mv.ScrollView.AddSubView(mv.OutputView.View)
	app.AddViewController(mv.ScrollView)
	return &mv
}

func (mv *MainView) Resize(w, h int) {
	mv.ScrollView.SetBounds(egg.MakeBounds(0, 2, w, h-2))
	mv.OutputView.SetBounds(egg.MakeBounds(0, 0, w, h+10))
}

func (mv *MainView) SetActiveFile(file *model.File) {
	mv.activeFile = file

	if file != nil {
		mv.OutputView.SetFile(file)
	}
}
