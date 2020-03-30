package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSubslice(t *testing.T) {
	var slice1 []string
	var slice2 []string

	slice1 = []string{"one", "two", "three"}
	slice2 = []string{"one", "two"}

	yes, remainder := IsCaseInsensitiveStringSubslice(slice1, slice2, false)
	assert.True(t, yes)
	assert.Equal(t, []string{"three"}, remainder)

	yes2, remainder2 := IsCaseInsensitiveStringSubslice(slice1[:2], slice2, false)
	assert.True(t, yes2)
	assert.Equal(t, []string{}, remainder2)

	no, remainder3 := IsCaseInsensitiveStringSubslice(slice2, slice1, false)
	assert.False(t, no)
	assert.Equal(t, slice2, remainder3)
}

func TestStringSubsliceCaseInsensitice(t *testing.T) {
	var slice1 []string
	var slice2 []string

	slice1 = []string{"one", "two", "three"}
	slice2 = []string{"One", "Two"}

	yes1, remainder1 := IsCaseInsensitiveStringSubslice(slice1, slice2, false)

	assert.True(t, yes1)
	assert.Equal(t, []string{"three"}, remainder1)
}

func TestStringSliceFlatter(t *testing.T) {
	res := StringSplitFlat("Hello       there my m an")
	assert.Equal(t, []string{"Hello", "there", "my", "m", "an"}, res)
}

func TestSanitizeString(t *testing.T) {
	res1 := SanitiseString("hello   there  this is a    string")
	assert.Equal(t, res1, "hello there this is a string")

	res2 := SanitiseString("  hello there ")
	assert.Equal(t, " hello there ", res2)
}

func TestLongestCommonPrefix(t *testing.T) {
	res1 := LongestCommonPrefix("hello", "he's a big boy")
	assert.Equal(t, "he", res1)

	res2 := LongestCommonPrefix("not", "here")
	assert.Equal(t, "", res2)
}
