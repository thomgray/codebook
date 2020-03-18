package model

import "github.com/thomgray/egg"

type AttributedStringSegment struct {
	Text       string
	Foreground egg.Color
	Background egg.Color
	Attributes egg.Attribute
}

type AttributedString = []AttributedStringSegment

func MakeAttributedStringSegment(text string, fg, bg egg.Color, atts egg.Attribute) AttributedStringSegment {
	return AttributedStringSegment{
		text, fg, bg, atts,
	}
}

func MakeAttributedString(segs ...AttributedStringSegment) AttributedString {
	return AttributedString(segs)
}

func MakeASFromPlainString(s string) AttributedString {
	return AttributedString{AttributedStringSegment{s, egg.ColorDefault, egg.ColorDefault, egg.AttrNormal}}
}
