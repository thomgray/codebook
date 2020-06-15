package model

import "os"

type AutocompleteResult struct {
	Str   string
	IsDir bool
}

func (ar *AutocompleteResult) CompletionStr() string {
	if ar.IsDir {
		return ar.Str + string(os.PathSeparator)
	}
	return ar.Str
}
