package htmlrender

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLongestNonBreakingSegment(t *testing.T) {
	str := "one two three four five"

	var res string
	var remainder string

	res, remainder = findLongestNonBreakingSegment(str, 7)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 8)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 9)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 10)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 11)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 12)
	assert.Equal(t, "one two", res)
	assert.Equal(t, "three four five", remainder)

	res, remainder = findLongestNonBreakingSegment(str, 13)
	assert.Equal(t, "one two three", res)
	assert.Equal(t, "four five", remainder)
}

func TestSliceForLineWhenWholeStringFits(t *testing.T) {
	str := "one two three four five"
	lineW := 100
	maxW := 100

	var slice string
	var remainder string
	var finished bool
	slice, remainder, finished = sliceForLine(str, lineW, maxW)

	assert.Equal(t, "one two three four five", slice)
	assert.Equal(t, "", remainder)
	assert.Equal(t, true, finished)
}

func TestSliceForLine(t *testing.T) {
	str := "one two three four five"
	maxW := 100

	var slice string
	var remainder string
	var finished bool

	for i := 0; i < 25; i++ {
		slice, remainder, finished = sliceForLine(str, i, maxW)

		var expectedLine string
		var expectedRemainder string
		var expectedConplete bool

		if i < 3 {
			expectedLine = ""
			expectedRemainder = str
			expectedConplete = false
		} else if i < 7 {
			expectedLine = "one"
			expectedRemainder = "two three four five"
			expectedConplete = false
		} else if i < 13 {
			expectedLine = "one two"
			expectedRemainder = "three four five"
			expectedConplete = false
		} else if i < 18 {
			expectedLine = "one two three"
			expectedRemainder = "four five"
			expectedConplete = false
		} else if i < 23 {
			expectedLine = "one two three four"
			expectedRemainder = "five"
			expectedConplete = false
		} else {
			expectedLine = "one two three four five"
			expectedRemainder = ""
			expectedConplete = true
		}

		assert.Equal(t, expectedLine, slice)
		assert.Equal(t, expectedRemainder, remainder)
		assert.Equal(t, expectedConplete, finished)
	}
}

func TestSliceForLineWhenNoBreak(t *testing.T) {
	var str = "onehellolongwordthatdoesntbreak is that a fact?"

	var slice string
	var remainder string
	var finished bool

	slice, remainder, finished = sliceForLine(str, 10, 11)

	assert.Equal(t, "", slice)
	assert.Equal(t, str, remainder)
	assert.Equal(t, false, finished)

	// does chop if the next line isn't going to be any wider
	slice, remainder, finished = sliceForLine(str, 10, 10)

	assert.Equal(t, "onehellolo", slice)
	assert.Equal(t, "ngwordthatdoesntbreak is that a fact?", remainder)
	assert.Equal(t, false, finished)
}
