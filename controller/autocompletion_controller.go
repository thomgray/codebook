package controller

import (
	"strings"

	"github.com/thomgray/codebook/model"
)

// todo
// should be case-insensitive
func traverseDocumentForAutocompletes(doc *model.Document, query string) [][]string {
	res := make([][]string, 0)
	thisDocSearchTerm := doc.Heading.Context[model.ContextSearchTerm]
	if strings.HasPrefix(thisDocSearchTerm, query) {
		// the query is contained in the documet, so it is a candidate auto-complete
		res = append(res, []string{thisDocSearchTerm})
	} else if strings.HasPrefix(query, thisDocSearchTerm) {
		// this document is included in the search term, so it or a child of it may match
		qTrimmed := strings.TrimPrefix(query, thisDocSearchTerm)
		qTrimmed = strings.TrimLeft(qTrimmed, " ")
		for _, sd := range doc.SubDocuments {
			// need to trim the working query,
			subDocResults := traverseDocumentForAutocompletes(sd, qTrimmed)
			for _, sdr := range subDocResults {
				// prepend "this search term" to any results

				prefixed := append([]string{thisDocSearchTerm}, sdr...)
				res = append(res, prefixed)
			}
		}
	} else {
		// the is no match, so can bottle out
	}

	return res
}
