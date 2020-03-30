package controller

import (
	"github.com/thomgray/codebook/model"
)

type query struct {
	q                string
	hasTrailingSpace bool
}

func documentMatchesQuery(doc *model.Document, q query) (yes bool, remainder string) {
	// docStr := doc.Heading.Context[model.ContextSearchTerm]
	return false, ""
}

func stringHasPrefixFlex(str string, prefix string) (yes bool, remainder string) {
	if len(prefix) > len(str) {
		return false, str
	}
	return false, ""
}
