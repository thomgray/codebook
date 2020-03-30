package controller

import (
	"github.com/thomgray/codebook/model"
)

func TopLevelDocument(doc *model.Document) *model.Document {
	d := doc
	for d.Super != nil {
		d = d.Super
	}
	return d
}
