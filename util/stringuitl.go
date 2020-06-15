package util

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func ReadLines(bytes []byte) []string {
	lines := make([]string, 0)

	var f []rune = make([]rune, 0)
	for len(bytes) > 0 {
		char, l := utf8.DecodeRune(bytes)
		bytes = bytes[l:]

		if char == '\n' {
			f, lines = _flushPathBuffer(f, lines)
		} else {
			f = append(f, char)
		}
	}
	f, lines = _flushPathBuffer(f, lines)
	return lines
}

func _flushPathBuffer(buffer []rune, array []string) ([]rune, []string) {
	if len(buffer) > 0 {
		return buffer[0:0], append(array, string(buffer))
	}
	return buffer, array
}

func StringSliceFilter(in []string, pred func(string) bool) []string {
	out := []string{}
	for _, s := range in {
		if pred(s) {
			out = append(out, s)
		}
	}
	return out
}

func StringSplitFlat(s string) []string {
	return StringSliceFlatten(strings.Split(s, " "))
}

func StringSliceFlatten(in []string) []string {
	return StringSliceFilter(in, func(s string) bool {
		return s != ""
	})
}

// func IsCaseInsensitiveStringSubslice(slice, subslice []string) (is bool, remainder []string) {
// 	if len(subslice) > len(slice) {
// 		return false, slice
// 	}

// 	x, sub := LongestCommonSubSlice(slice, subslice)
// 	if len(sub) == len(subslice) {
// 		return true, slice[x:]
// 	} else {
// 		return false, slice
// 	}
// }

func IsCaseInsensitiveStringSubslice(slice, subslice []string, proper bool) (is bool, remainder []string) {
	if len(subslice) > len(slice) {
		return false, slice
	} else if proper && len(subslice) == len(slice) {
		return false, slice
	}

	x, sub := LongestCommonSubSlice(slice, subslice)
	if len(sub) == len(subslice) {
		return true, slice[x:]
	} else {
		return false, slice
	}
}

func LongestCommonSubSlice(s1, s2 []string) (int, []string) {
	i := minInt(len(s1), len(s2))
	x := 0
	for ; x < i; x++ {
		str1 := s1[x]
		str2 := s2[x]

		if !strings.EqualFold(str1, str2) {
			break
		}
	}
	return x, s1[:x]
}

func minInt(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

func StringHasSubstringFlex(str, substr string) (bool, string) {
	return false, ""
}

var __spaceRegex, _ = regexp.Compile("\\s+")

func SanitiseString(str string) string {
	return __spaceRegex.ReplaceAllLiteralString(str, " ")
}

func LongestCommonPrefix(str1, str2 string) string {
	str1B := []byte(str1)
	str2B := []byte(str2)
	res := make([]byte, 0)

	for true {
		r1, w1 := utf8.DecodeRune(str1B)
		r2, w2 := utf8.DecodeRune(str2B)

		if r1 == utf8.RuneError || r2 == utf8.RuneError {
			break
		}
		if r1 == r2 {
			res = append(res, str1B[:w1]...)
		} else {
			break
		}
		str1B = str1B[w1:]
		str2B = str2B[w2:]
	}

	return string(res)
}

func StringSliceContains(slice []string, s string) bool {
	for _, str := range slice {
		if str == s {
			return true
		}
	}
	return false
}
