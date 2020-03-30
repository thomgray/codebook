package controller

import (
	"strings"

	"github.com/thomgray/codebook/model"
)

// TraveralMode ...
type TraveralMode = int8

// TraveralMode ...
const (
	TraversalModeDefault TraveralMode = iota
	TraveralModeHere
	TraveralModeRoot
	TraveralModeExt
)

func getQueryAndMode(cmd string) (string, TraveralMode) {
	mode := TraversalModeDefault
	str := cmd
	if strings.HasPrefix(cmd, ". ") {
		str = strings.TrimLeft(strings.TrimPrefix(str, "."), " ")
		mode = TraveralModeHere
	} else if strings.HasPrefix(cmd, "/ ") {
		str = strings.TrimLeft(strings.TrimPrefix(str, "/"), " ")
		mode = TraveralModeRoot
	} else if strings.HasPrefix(cmd, "* ") {
		str = strings.TrimLeft(strings.TrimPrefix(str, "*"), " ")
		mode = TraveralModeExt
	}

	return str, mode
}

func queryDocument(doc *model.Document, query string, includeThisOne bool) *model.Document {
	// query should already be sanitised
	var traverse func(d *model.Document, q string) *model.Document

	traverse = func(d *model.Document, q string) *model.Document {
		// search term should be trimmed already
		st := strings.ToLower(d.SearchTerm)
		if st == q {
			return d
		} else if strings.HasPrefix(q, st) {
			remainder := strings.TrimLeft(strings.TrimPrefix(q, st), " ")

			for _, subd := range d.SubDocuments {
				res := traverse(subd, remainder)
				if res != nil {
					return res
				}
			}
		}
		return nil
	}

	if includeThisOne {
		return traverse(doc, query)
	}
	for _, subd := range doc.SubDocuments {
		res := traverse(subd, query)
		if res != nil {
			return res
		}
	}
	return nil
}
