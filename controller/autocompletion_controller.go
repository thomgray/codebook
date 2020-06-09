package controller

import (
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/thomgray/codebook/util"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/codebook/model"
)

type autocompleteResult struct {
	toComplete       string
	entireString     string
	comparableString string
}

func (ar autocompleteResult) prefix(whole, comparable string) autocompleteResult {
	return autocompleteResult{
		toComplete:       whole + " " + ar.entireString,
		entireString:     whole + " " + ar.entireString,
		comparableString: comparable + " " + ar.comparableString,
	}
}

func (mc *MainController) handleAutocompleteNote(str string) {
	completeSuggestions := mc.FileManager.SuggestPaths(str)

	if len(completeSuggestions) == 1 {
		newQuery := completeSuggestions[0]
		mc.InputView.SetTextContentString(newQuery)
		mc.InputView.SetCursorX(runewidth.StringWidth(newQuery))
		app.ReDraw()
	}
	// querySanitied := strings.TrimLeft(util.SanitiseString(str), " ")
	// query, mode := getQueryAndMode(querySanitied)
	// queryLcase := strings.ToLower(query)

	// res := make([]autocompleteResult, 0)

	// if mc.activeDocument != nil {
	// 	if mode == TraversalModeDefault || mode == TraveralModeHere {
	// 		for _, c := range mc.activeDocument.SubDocuments {
	// 			res = append(res, findCompletionsForDocument(c, queryLcase)...)
	// 		}
	// 	} else if mode == TraveralModeRoot {
	// 		super := TopLevelDocument(mc.activeDocument)
	// 		for _, c := range super.SubDocuments {
	// 			res = append(res, findCompletionsForDocument(c, queryLcase)...)
	// 		}
	// 	}
	// }
	// if len(res) > 0 {
	// 	handleCompleteWithCompletions(mc, res, mode)
	// 	return
	// }

	// if mode == TraversalModeDefault || mode == TraveralModeExt {
	// 	for _, f := range mc.FileManager.Files {
	// 		if f.Document != nil {
	// 			res = append(res, findCompletionsForDocument(f.Document, queryLcase)...)
	// 		}
	// 	}
	// }
	// handleCompleteWithCompletions(mc, res, mode)
}

func findCompletionsForDocument(doc *model.Document, query string) []autocompleteResult {
	// query should already be sanitised
	res := make([]autocompleteResult, 0)
	docstr := doc.SearchTerm
	docstrLcase := strings.ToLower(docstr)

	if strings.HasPrefix(docstrLcase, query) {
		// the query is contained in the doc title, so this document is a candidate
		res = append(res, autocompleteResult{
			toComplete:       docstr[runewidth.StringWidth(query):],
			entireString:     docstr,
			comparableString: docstrLcase,
		})
	} else if strings.HasPrefix(query, docstrLcase) {
		// the query contains the document
		for _, c := range doc.SubDocuments {
			chopped := strings.TrimLeft(strings.TrimPrefix(query, docstrLcase), " ")
			subRes := findCompletionsForDocument(c, chopped)
			for _, subr := range subRes {
				res = append(res, subr.prefix(docstr, docstrLcase))
			}
		}
	}

	return res
}

func handleCompleteWithCompletions(mc *MainController, res []autocompleteResult, mode TraveralMode) {
	var modePrefix string = ""
	switch mode {
	case TraveralModeHere:
		modePrefix = ". "
	case TraveralModeRoot:
		modePrefix = "/ "
	case TraveralModeExt:
		modePrefix = "* "
	}
	if len(res) == 1 {
		totalStr := fmt.Sprintf("%s%s ", modePrefix, res[0].entireString)
		mc.InputView.SetTextContentString(totalStr)
		mc.InputView.SetCursorX(runewidth.StringWidth(totalStr))
		app.ReDraw()
	} else if len(res) > 1 {
		first := res[0]
		rest := res[1:]

		prefix := first.comparableString
		for _, other := range rest {
			prefix = util.LongestCommonPrefix(prefix, other.comparableString)
			if prefix == "" {
				break
			}
		}

		//todo could be better
		if prefix != "" {
			toCompl := first.entireString[:len(prefix)]
			toCompl = modePrefix + toCompl
			mc.InputView.SetTextContentString(toCompl)
			mc.InputView.SetCursorX(runewidth.StringWidth(toCompl))
			app.ReDraw()
		}
	}
}

func traverseDocumentForAutocompletes(doc *model.Document, query []string, trailingSpace bool) [][]string {
	res := make([][]string, 0)
	thisDocSearchTerm := doc.SearchTerm
	thisDocSearchTermSplit := util.StringSplitFlat(thisDocSearchTerm)

	log.Println("Asking for query candidates for", query, "given this doc has search term", thisDocSearchTermSplit)

	// string slice thing isn't working so well, consider using plain strings
	if yes, _ := util.IsCaseInsensitiveStringSubslice(thisDocSearchTermSplit, query, trailingSpace); yes && !trailingSpace {
		log.Println("Query contained in document")
		// the query is contained in the document term
		res = append(res, thisDocSearchTermSplit)
	} else if yes, remainder := util.IsCaseInsensitiveStringSubslice(query, thisDocSearchTermSplit, false); yes {
		log.Println("Document contained in query, looking for sub-doc")
		for _, sd := range doc.SubDocuments {
			// need to trim the working query,
			subDocResults := traverseDocumentForAutocompletes(sd, remainder, trailingSpace)
			for _, sdr := range subDocResults {
				// prepend "this search term" to any results

				prefixed := append(thisDocSearchTermSplit, sdr...)
				res = append(res, prefixed)
			}
		}
	} else {
		log.Println("Founds nothing")
	}

	log.Println("Autocomplete candidates = ")
	for _, r := range res {
		for _, rr := range r {
			log.Printf("'%s' ", rr)
		}
	}
	return res
}

// should not be sensitive to case or redundant whitespace
// assume that input is space-separated strings (i.e. no whitespace)
func longestCommonPath(in [][]string) ([]string, bool) {
	var endsAtWordBreak bool = false
	var out []string = make([]string, 0)
	if len(in) == 0 {
		return out, endsAtWordBreak
	}
	head := in[0]
	tail := in[1:]
	var minDepth int = len(head)
	for _, q := range tail {
		l := len(q)
		if l < minDepth {
			minDepth = l
		}
	}

here:
	for x := 0; x < minDepth; x++ {
		thisOne := head[x]
		for _, s := range tail {
			thatOne := s[x]
			thisOne = longestCommonSubstr(thisOne, thatOne)
			if thisOne == "" {
				endsAtWordBreak = true
				break here
			}
		}
		out = append(out, thisOne)
		if thisOne != head[x] {
			// trancated at this point, to end here
			break
		}

	}
	return out, endsAtWordBreak
}

func longestCommonSubstr(s1, s2 string) string {
	s1d := []byte(s1)
	s2d := []byte(s2)
	inCommon := make([]rune, 0)
	for {
		r1, s1 := utf8.DecodeRune(s1d)
		r2, s2 := utf8.DecodeRune(s2d)

		if r1 == utf8.RuneError || r2 == utf8.RuneError {
			break
		} else if r1 != r2 {
			// TODO check case insensitively
			// runes are different
			break
		}
		s1d = s1d[s1:]
		s2d = s2d[s2:]
		inCommon = append(inCommon, r1)
	}

	return string(inCommon)
}
