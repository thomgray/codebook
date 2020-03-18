package controller

import (
	"github.com/thomgray/codebook/constants"
	"github.com/thomgray/codebook/model"

	"github.com/thomgray/codebook/config"
	"github.com/thomgray/codebook/view"
	"github.com/thomgray/egg"
)

// MainController ...
type MainController struct {
	View        *view.MainView
	InputView   *view.InputView
	ModalMenu   *view.ModalMenu
	Config      *config.Config
	FileManager *model.FileManager
}

// Mode ...
type Mode uint8

// Mode ...
const (
	ModeInput Mode = iota
	ModeMenu
)

var app *egg.Application
var mode Mode = ModeInput
var inputMode constants.InputMode = constants.InputModeTraverse

// InitMainController ...
func InitMainController(config *config.Config) *MainController {
	app = egg.InitOrPanic()

	mc := MainController{
		View:        view.MakeMainView(app),
		InputView:   view.MakeInputView(app),
		ModalMenu:   view.MakeModalMenu(),
		Config:      config,
		FileManager: model.MakeFileManager(),
	}

	app.AddViewController(mc.ModalMenu)
	app.OnKeyEvent(mc.keyEventDelegate)

	mc.init()

	return &mc
}

func (mc *MainController) init() {
	mc.reloadFiles()
	bootstrapCommands()
}

func (mc *MainController) reloadFiles() {
	mc.FileManager.LoadFiles(mc.Config.NotePaths)
}

func (mc *MainController) keyEventDelegate(e *egg.KeyEvent) {
	switch e.Key {
	case egg.KeyEsc:
		m := mc.toggleMode()
		mc.ModalMenu.SetVisible(m == ModeMenu)
		app.ReDraw()
		return
	}

	if mode == ModeInput {
		mc.handleEventInputMode(e)
	} else if mode == ModeMenu {
		mc.handleEventMenuMode(e)
	}
}

func (mc *MainController) handleEventInputMode(e *egg.KeyEvent) {
	// shoud be based on cursor position, but egg doesn't expose that
	if mc.InputView.GetCursorX() == 0 {
		switch e.Char {
		case '?':
			e.StopPropagation = true
			mc.setInputMode(constants.InputModeSearch)
			app.ReDraw()
			return
		case '>':
			e.StopPropagation = true
			mc.setInputMode(constants.InputModeTraverse)
			app.ReDraw()
			return
		case ':':
			e.StopPropagation = true
			mc.setInputMode(constants.InputModeCommand)
			app.ReDraw()
			return
		}
	}

	switch e.Key {
	case egg.KeyEnter:
		mc.handleEnter(e)
	case egg.KeyTab:
		e.StopPropagation = true
		mc.handleAutocomplete(mc.InputView.GetTextContentString())
	case egg.KeyArrowUp, egg.KeyArrowDown:
		e.StopPropagation = true
		mc.View.ScrollView.ReceiveKeyEvent(e)
		app.ReDraw()
	}
}

func (mc *MainController) handleEventMenuMode(e *egg.KeyEvent) {
	switch e.Char {
	case 'x':
		app.Stop()
	}
}

func (mc *MainController) toggleMode() Mode {
	if mode == ModeInput {
		app.SetFocusedView(nil)
		mode = ModeMenu
	} else {
		mc.InputView.GainFocus()
		mode = ModeInput
	}
	return mode
}

func (mc *MainController) setInputMode(m constants.InputMode) {
	if inputMode == m {
		return
	}
	// old := mode
	inputMode = m
	mc.InputView.SetMode(inputMode)
}

func (mc *MainController) handleEnter(e *egg.KeyEvent) {
	e.StopPropagation = true
	txt := mc.InputView.GetTextContentString()
	switch inputMode {
	case constants.InputModeTraverse:
		mc.handleTraverse(txt)
	case constants.InputModeSearch:
		mc.handleSearch(txt)
	case constants.InputModeCommand:
		mc.handleCommand(txt)
	}
}

func (mc *MainController) handleSearch(str string) {
	var f *model.File = nil
	for _, file := range mc.FileManager.Files {
		if file.Name == str {
			f = file
			break
		}
	}

	if f != nil {
		mc.View.SetActiveFile(f)
		mc.InputView.SetTextContentString("")
		app.ReDraw()
	}
}

func (mc *MainController) handleTraverse(str string) {
	var f *model.File = nil
	for _, file := range mc.FileManager.Files {
		if file.Name == str {
			f = file
			break
		}
	}

	if f != nil {
		mc.View.SetActiveFile(f)
		mc.InputView.SetTextContentString("")
		app.ReDraw()
	}
}

func (mc *MainController) handleAutocomplete(str string) {

}

// Start ...
func (mc *MainController) Start() {
	defer app.Start()
}
