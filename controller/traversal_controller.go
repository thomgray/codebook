package controller

import (
	"strings"

	"github.com/thomgray/codebook/model"
)

func (mc *MainController) handleTraverse(str string) {
	allPaths := mc.FileManager.FindSupportedFilePaths()
	complete := func(file *model.File) {
		mc.SetActiveFile(file)
		app.ReDraw()
	}

	// can do something with matching dirs perhaps?
	// dirPaths := mc.FileManager.FindPossibleBasePaths()
	// log.Println(">>>>>> dirs = ", dirPaths)

	for _, p := range allPaths {
		qp := p.QueryPath()

		if strings.EqualFold(qp, str) {
			// is exact match
			f := model.LoadCodeFile(p.Full)
			if f != nil {
				complete(f)
			}
		}
	}
}
