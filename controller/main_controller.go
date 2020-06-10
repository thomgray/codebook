package controller

import (
	"log"
	"strings"

	"github.com/thomgray/codebook/constants"
	"github.com/thomgray/codebook/model"

	"github.com/thomgray/codebook/config"
	"github.com/thomgray/codebook/view"
	"github.com/thomgray/egg"
)

// MainController ...
type MainController struct {
	View           *view.MainView
	InputView      *view.InputView
	ModalMenu      *view.ModalMenu
	CompletionView *view.CompletionView
	Config         *config.Config
	FileManager    *model.FileManager
	activeDocument *model.Document
	activeFile     *model.File
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
		View:           view.MakeMainView(app),
		InputView:      view.MakeInputView(app),
		ModalMenu:      view.MakeModalMenu(),
		CompletionView: view.MakeCompletionView(),
		Config:         config,
		FileManager:    model.MakeFileManager(config),
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
	// mc.FileManager.LoadFiles(mc.Config.NotePaths)
}

func (mc *MainController) keyEventDelegate(e *egg.KeyEvent) {
	switch e.Key {
	case egg.KeyEsc:
		// m := mc.toggleMode()
		// mc.ModalMenu.SetVisible(m == ModeMenu)
		// app.ReDraw()
		return
	}

	if mode == ModeInput {
		mc.handleEventInputMode(e)
	} else if mode == ModeMenu {
		mc.handleEventMenuMode(e)
	}
}

func (mc *MainController) handleEventInputMode(e *egg.KeyEvent) {
	if mc.InputView.GetCursorX() == 0 {
		switch e.Char {
		case '?':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeSearch)
			app.ReDraw()
			return
		case '>':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeTraverse)
			app.ReDraw()
			return
		case ':':
			e.SetPropagate(false)
			mc.setInputMode(constants.InputModeCommand)
			app.ReDraw()
			return
		}
	}

	switch e.Key {
	case egg.KeyEnter:
		mc.handleEnter(e)
	case egg.KeyTab:
		e.SetPropagate(false)
		mc.handleAutocomplete(mc.InputView.GetTextContentString())
	case egg.KeyArrowUp, egg.KeyArrowDown:
		e.SetPropagate(false)
		mc.View.HandleKeyEvent(e)
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
	e.SetPropagate(false)
	mc.CompletionView.SetVisible(false)
	txt := mc.InputView.GetTextContentString()
	switch inputMode {
	case constants.InputModeTraverse:
		if !mc.handleSpecial(txt) {
			mc.handleTraverse(txt)
		}
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

	if f != nil && f.Document != nil {
		mc.SetActiveDocument(f.Document)
		mc.InputView.SetTextContentString("")
		app.ReDraw()
	}
}

func (mc *MainController) SetActiveDocument(d *model.Document) {
	mc.activeDocument = d
	mc.View.SetActiveDocument(d)
	mc.InputView.SetTextContentString("")
	mc.InputView.SetCursorX(0)
}

func (mc *MainController) SetActiveFile(f *model.File) {
	mc.activeFile = f
	mc.View.SetActiveFile(f)
	mc.InputView.SetTextContentString("")
	mc.InputView.SetCursorX(0)
}

func (mc *MainController) handleTraverse(strUntreated string) {
	complete := func(file *model.File) {
		mc.SetActiveFile(file)
		app.ReDraw()
	}

	strTrimmed := strings.TrimSpace(strings.ToLower(strUntreated))
	log.Println("traversing", strTrimmed)
	// pathFragments := strings.Split(strTrimmed, "/")

	files := mc.FileManager.TraversePath(strTrimmed)
	if len(files) == 1 {
		file := files[0]
		complete(file)
		log.Println("Active file set")
	}
}

func matchesLocation(f *model.File, path string) bool {
	for _, loc := range f.Locations {
		if strings.EqualFold(loc.RelativePathWithName, path) {
			// path matches fully
			return true
		}
	}
	return false
}

// func (mc *MainController)handleTraverse(strUntreated string) {
// 	complete := func(doc *model.Document) {
// 		mc.SetActiveDocument(doc)
// 		app.ReDraw()
// 	}

// 	strTrimmed := strings.TrimSpace(strings.ToLower(strUntreated))
// 	str, mode := getQueryAndMode(strTrimmed)
// 	log.Printf("Doing a traverse w/ mode %d and command %s", mode, str)

// 	var d *model.Document = nil

// 	if mc.activeDocument != nil {
// 		doc := mc.activeDocument

// 		if mode == TraversalModeDefault || mode == TraveralModeHere {
// 			d = queryDocument(doc, str, false)
// 			if d != nil {
// 				complete(d)
// 				return
// 			}
// 		} else if mode == TraveralModeRoot {
// 			top := TopLevelDocument(doc)
// 			d = queryDocument(top, str, false)
// 			if d != nil {
// 				complete(d)
// 				return
// 			}
// 		}
// 	}

// 	if mode == TraversalModeDefault || mode == TraveralModeExt {
// 		for _, file := range mc.FileManager.Files {
// 			if file.Document != nil {
// 				d = queryDocument(file.Document, str, true)
// 				if d != nil {
// 					complete(d)
// 					return
// 				}
// 			}
// 		}
// 	}
// }

func (mc *MainController) handleSpecial(str string) bool {
	overruled := false
	cleaned := strings.TrimSpace(str)
	if mc.activeDocument != nil {
		switch cleaned {
		case ".":
			mc.SetActiveDocument(mc.activeDocument)
			app.ReDraw()
			overruled = true
		case "*":
			// unset active document
			mc.SetActiveDocument(nil)
			app.ReDraw()
			overruled = true
		case "..":
			if mc.activeDocument != nil && mc.activeDocument.Super != nil {
				mc.SetActiveDocument(mc.activeDocument.Super)
				overruled = true
				app.ReDraw()
			}
		case "/":
			if mc.activeDocument != nil && mc.activeDocument.Super != nil {
				super := mc.activeDocument.Super
				for super.Super != nil {
					super = super.Super
				}
				mc.SetActiveDocument(super)
				overruled = true
				app.ReDraw()
			}
		}
	}
	return overruled
}

func (mc *MainController) handleAutocomplete(str string) {
	switch inputMode {
	case constants.InputModeTraverse:
		mc.handleAutocompleteNote(str)
	}
}

func (mc *MainController) Start() {
	defer app.Start()
}
