package controller

import (
	"log"
	"strings"

	"github.com/thomgray/codebook/model"
)

func queryDocument(doc *model.Document, query string) *model.Document {
	var traverse func(d *model.Document, q string) *model.Document

	traverse = func(d *model.Document, q string) *model.Document {
		for _, subdoc := range d.SubDocuments {
			st := subdoc.Heading.Context[model.ContextSearchTerm]
			log.Printf("Traversing heading %s\n", st)
			if strings.HasPrefix(q, st) {
				remaining := strings.Trim(strings.TrimPrefix(q, st), " ")
				if remaining != "" {
					res := traverse(subdoc, remaining)
					if res != nil {
						return res
					}
				} else {
					return subdoc
				}
			}
		}
		return nil
	}
	return traverse(doc, query)
}
