package htmlrender

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/thomgray/egg"
)

type attStringSegment struct {
	str  string
	fg   egg.Color
	bg   egg.Color
	atts egg.Attribute
}

type attString = []attStringSegment

func makeAttString(str string, fg, bg egg.Color, atts egg.Attribute) attString {
	seg := attStringSegment{
		str, fg, bg, atts,
	}
	return []attStringSegment{seg}
}

func joinAttStrings(strings ...attString) attString {
	str := make([]attStringSegment, 0)
	for _, s := range strings {
		str = append(str, s...)
	}
	return str
}

func attStringWidth(as attString) int {
	res := 0
	for _, seg := range as {
		res += runewidth.StringWidth(seg.str)
	}
	return res
}

func takeLine(as attString, length, maxLength int) (attString, attString) {
	pre := make([]attStringSegment, 0)
	aft := as
	lengthRemaining := length
	for len(aft) > 0 {
		slice := aft[0]
		sliceW := runewidth.StringWidth(slice.str)
		if sliceW < lengthRemaining {
			pre = append(pre, slice)
			lengthRemaining -= sliceW
		} else {
			// need to slice up the slice
			strPre, strRemainder := takeString(slice.str, lengthRemaining)
			asPre := attStringSegment{strPre, slice.fg, slice.bg, slice.atts}
			asRem := attStringSegment{strRemainder, slice.fg, slice.bg, slice.atts}
			pre = append(pre, asPre)
			aft = append([]attStringSegment{asRem}, aft[1:]...)
			break
		}
		aft = aft[1:]
	}
	return pre, aft
}

func takeString(str string, length int) (string, string) {
	if runewidth.StringWidth(str) <= length {
		return str, ""
	}
	pieces := strings.Split(str, " ")
	lengthRemaining := length
	foreSlice := make([]string, 0)
	aftSlice := pieces
	for _, piece := range pieces {
		pieceW := runewidth.StringWidth(piece)
		if pieceW <= lengthRemaining {
			lengthRemaining -= pieceW + 1 // for whitespace
			foreSlice = append(foreSlice, piece)
			aftSlice = aftSlice[1:]
		}
	}

	return strings.Join(foreSlice, " "), strings.Join(aftSlice, " ")
}
