package util

import "unicode/utf8"

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

func StringSliceFlatten(in []string) []string {
	return StringSliceFilter(in, func(s string) bool {
		return s != ""
	})
}
