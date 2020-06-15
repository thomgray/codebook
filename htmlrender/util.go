package htmlrender

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

func normaliseText(txt string) string {
	regex := regexp.MustCompile("\\s+")
	s := regex.ReplaceAllString(txt, " ")
	return s
}

func killTabsDead(txt string) string {
	return strings.ReplaceAll(txt, "\t", "  ")
}

func sliceForLine(str string, lineRemainder, maxwidth int) (string, string, bool) {
	if len(str) <= lineRemainder {
		return str, "", true
	}

	slice := str[:lineRemainder]
	// peek := str[lineRemainder]

	remainder := str[lineRemainder:]
	remainder = strings.TrimPrefix(remainder, " ") // remove leading whitespace

	// if strings.HasSuffix(slice, " ") {
	// 	// were were fortunate that this line fits to a word break exactly
	// 	return slice[:len(slice)-1], remainder, false
	// } else if peek == ' ' {
	// 	// similar situation, only the break in on the next line
	// 	return slice, remainder, false
	// } else if isWordBreakAt(slice[len(slice)-1], peek) {
	// 	// there is a legit work break at this line anyway, so
	// 	return slice, remainder, false
	// }
	// should check if it ends in other word-breaking characters (e.g. '-')
	// otherwise we assume we have chopped a word, and so must wrap
	longestContinuous, remainders := findLongestNonBreakingSegment(str, lineRemainder)

	if runewidth.StringWidth(longestContinuous) > lineRemainder {
		// still can't fit on the line
		if lineRemainder < maxwidth {
			// might be able to fit the word on the next line, so wrap the whole thing on the next line
			return "", str, false
		}
		// otherwise, we have to just chop the word
		return slice, remainder, false
	}
	// otherwise, we can break
	return longestContinuous, remainders, false
}

func findLongestNonBreakingSegment(s string, rng int) (string, string) {
	if runewidth.StringWidth(s) <= rng {
		return s, ""
	}

	splitWords := strings.Split(s, " ") // split into words
	working := splitWords[0]
	splitWords = splitWords[1:]
	for i, str := range splitWords {
		withNext := fmt.Sprintf("%s %s", working, str)
		if len(withNext) > rng {
			return working, strings.Join(splitWords[i:], " ")
		}
		working = withNext
	}

	return working, ""
}

func isWordBreakAt(char1, char2 rune) bool {
	return false
}
