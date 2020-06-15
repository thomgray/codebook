package controller

import (
	"fmt"
	"strings"

	"github.com/thomgray/codebook/model"
	"github.com/thomgray/codebook/util"
	"github.com/thomgray/egg"
)

type command struct {
	aliases     []string
	desctiption string
	action      func(*MainController, []string) bool
	keyhandler  func(ke *egg.KeyEvent)
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

			mc.View.OutputView.CustomDraw(func(c egg.Canvas) {
				for i, spath := range sp {
					c.DrawString2(spath, 0, i)
				}
			})

			return true
		},
	},
	&command{
		aliases:     []string{"l"},
		desctiption: "List top level documents",
		action: func(mc *MainController, args []string) bool {
			allFilesPaths := mc.FileManager.FindSupportedFilePaths()

			mc.View.OutputView.CustomDraw(func(c egg.Canvas) {
				for i, f := range allFilesPaths {
					c.DrawString2(f.QueryPath(), 0, i)
				}
			})

			curBounds := mc.View.OutputView.GetBounds()
			if curBounds.Height < len(allFilesPaths) {
				curBounds.Height = len(allFilesPaths)
				mc.View.OutputView.SetBounds(curBounds)
				mc.View.ScrollView.ReDraw()
			}
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
	&command{
		aliases:     []string{"cd"},
		desctiption: "Change document root",
		action: func(mc *MainController, args []string) bool {
			positional, flags, _ := parseOptions(args)
			if len(positional) == 0 {
				return false
			}

			path := positional[0]
			if info, exists := util.PathExists(path); exists && info.IsDir() {
				mc.Config.SetCurrentDocRoot(path)

				if util.StringSliceContains(flags, "default") {
					// make this the default root
					mc.Config.SetDefaultDocRoot(path)
				}
				return true
			}
			return false
		},
	},
	&command{
		aliases:     []string{"pwd"},
		desctiption: "Output current document root",
		action: func(mc *MainController, args []string) bool {
			root := mc.Config.DocumentRoot()
			str := "?"
			if root != nil {
				str = *root
			}

			mc.View.OutputView.CustomDraw(func(c egg.Canvas) {
				c.DrawString2(str, 0, 0)
			})

			return true
		},
	},
}

func parseOptions(args []string) ([]string, []string, map[string]string) {
	positional := make([]string, 0)
	flags := make([]string, 0)
	options := make(map[string]string)

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			flags = append(flags, strings.TrimPrefix(arg, "--"))
		} else if arg == "-d" {
			flags = append(flags, "default")
		} else if strings.HasPrefix(arg, "-") {
			// noop - ignore unknown shorthand
		} else {
			// positional
			positional = append(positional, arg)
		}
	}

	return positional, flags, options
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
		mc.View.OutputView.CustomDraw(func(c egg.Canvas) {
			y := 0
			for _, cmd := range commands {
				aliases := strings.Join(cmd.aliases, " | ")
				plain := fmt.Sprintf("- %s : %s", aliases, cmd.desctiption)

				c.DrawString2(plain, 0, y)
				y++
			}
		})
		return true
	}
}
