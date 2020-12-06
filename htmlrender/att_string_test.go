package htmlrender

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomgray/egg"
)

func TestAsStringTakeLineNoTruncate(t *testing.T) {
	as := makeAttString("hello", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)

	line, remainder := takeLine(as, 10, 100)
	assert.Equal(t, as, line)
	assert.Equal(t, attString{}, remainder)
}

func TestAsStringTakeLine(t *testing.T) {
	as := makeAttString("hello is it me you're looking for?", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)

	line, remainder := takeLine(as, 10, 100)
	pre := makeAttString("hello is", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)
	aft := makeAttString("it me you're looking for?", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)
	assert.Equal(t, pre, line)
	assert.Equal(t, aft, remainder)
}

func TestAsStringTakeLineMultipleSegments(t *testing.T) {
	as1 := makeAttString("hello!", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)
	as2 := makeAttString(" is it me you're looking for?", egg.ColorBlack, egg.ColorCyan, egg.AttrNormal)

	as := joinAttStrings(as1, as2)
	line, remainder := takeLine(as, 10, 100)

	pre1 := makeAttString("hello!", egg.ColorBlack, egg.ColorBlue, egg.AttrNormal)
	pre2 := makeAttString(" is", egg.ColorBlack, egg.ColorCyan, egg.AttrNormal)
	pre := joinAttStrings(pre1, pre2)

	aft := makeAttString("it me you're looking for?", egg.ColorBlack, egg.ColorCyan, egg.AttrNormal)
	assert.Equal(t, pre, line)
	assert.Equal(t, aft, remainder)
}
