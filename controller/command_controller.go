package controller

import (
	"fmt"
	"strings"

	"github.com/thomgray/codebook/model"
)

type command struct {
	aliases     []string
	desctiption string
	action      func(*MainController, []string) bool
}

var commands = []*command{
	&command{
		aliases:     []string{"h", "help"},
		desctiption: "Print this help message",
	},
	&command{
		aliases:     []string{"q", "quit"},
		desctiption: "Exit application",
		action: func(mc *MainController, args []string) bool {
			app.Stop()
			return true
		},
	},
	&command{
		aliases:     []string{"ls", "list"},
		desctiption: "List configured search paths",
		action: func(mc *MainController, args []string) bool {
			sp := mc.Config.SearchPaths
			as := make([]model.AttributedString, len(sp))
			for i, spath := range sp {
				as[i] = model.MakeASFromPlainString(spath)
			}
			mc.View.SetSpecialOutput(&as)
			return true
		},
	},
	&command{
		aliases:     []string{"l"},
		desctiption: "List top level documents",
		action: func(mc *MainController, args []string) bool {
			as := make([]model.AttributedString, 0)
			for _, f := range mc.FileManager.Files {
				if f.Document != nil {
					as = append(as, model.MakeASFromPlainString(f.Document.SearchTerm))
				}
			}
			mc.View.SetSpecialOutput(&as)
			return true
		},
	},
	&command{
		aliases:     []string{"sp-add", "+"},
		desctiption: "Add a search path",
		action: func(mc *MainController, args []string) bool {
			if len(args) == 0 {
				return false
			}
			sp := args[0]
			mc.Config.AddSearchPath(sp)
			// mc.Config.ReloadNotes()
			mc.reloadFiles()
			return true
		},
	},
	&command{
		aliases:     []string{"sp-remove", "-"},
		desctiption: "Remove a search path",
		action: func(mc *MainController, args []string) bool {
			if len(args) == 0 {
				return false
			}
			sp := args[0]
			mc.Config.RemoveSearchPath(sp)
			// mc.Config.ReloadNotes()
			mc.reloadFiles()
			return true
		},
	},
	&command{
		aliases:     []string{"reload"},
		desctiption: "Reload notes",
		action: func(mc *MainController, args []string) bool {
			mc.Config.Init() // to reload config
			mc.reloadFiles() // to reload files
			return true
		},
	},
}

func (mc *MainController) handleCommand(str string) {
	trimmed := strings.Trim(str, " ")
	split := strings.Split(trimmed, " ")
	if len(split) == 0 {
		return
	}
	cmdIn := split[0]
	args := make([]string, 0)
	split = split[1:]
	for _, s := range split {
		if s != "" {
			args = append(args, s)
		}
	}
	hit := false

here:
	for _, cmd := range commands {
		for _, alias := range cmd.aliases {
			if alias == cmdIn {
				hit = cmd.action(mc, args)
				if hit {
					break here
				}
			}
		}
	}

	if hit {
		mc.InputView.SetTextContentString("")
		mc.InputView.SetCursorX(0)
		app.ReDraw()
	}
}

var __helpTxt *[]model.AttributedString = nil

func GetHelp() *[]model.AttributedString {
	if __helpTxt == nil {
		initHelp()
	}
	return __helpTxt
}

func initHelp() {
	txt := make([]model.AttributedString, 0)

	for _, cmd := range commands {
		aliases := strings.Join(cmd.aliases, " | ")
		plain := fmt.Sprintf("- %s : %s", aliases, cmd.desctiption)
		txt = append(txt, model.MakeASFromPlainString(plain))
	}

	__helpTxt = &txt
}

func bootstrapCommands() {
	// needs to be done this way for circularity reasons :(
	commands[0].action = func(mc *MainController, args []string) bool {
		txt := GetHelp()
		mc.View.SetSpecialOutput(txt)
		return true
	}
}
