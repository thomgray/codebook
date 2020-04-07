package model

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type (
	ElementType uint8
	Attribution uint8
)

const (
	ElementTypeString ElementType = iota
	ElementTypeHeading
	ElementTypeCode
	ElementTypeQuote
	ElementTypeListItem
	ElementTypeUnorderedList
	ElementTypeOrderedList
	ElementTypeHorizontalRule
)

const (
	AttributionPlain Attribution = 1 << iota
	AttributionEmphasis
	AttributionBold
	AttributionCode
	AttributeAnchor
)

const (
	ContextSearchTerm string = "searchTerm"
	ContextHval       string = "hval"
)

type Document struct {
	Node         *html.Node
	Heading      *Element
	Elements     []*Element
	SubDocuments []*Document
	Super        *Document
	SearchTerm   string
}

type ContentSegment struct {
	Raw         string
	Attribution Attribution
	Context     map[string]string
}

type Element struct {
	Type        ElementType
	Tag         string
	Content     []*ContentSegment
	Context     map[string]string
	SubElements []*Element
}

func DocumentFromNode(n *html.Node, filename string) *Document {
	d := Document{}
	d.Node = n
	els := make([]*Element, 0)

	for node := n.FirstChild; node != nil; node = node.NextSibling {
		e := parseElement(node, false)
		if e != nil {
			els = append(els, e)
		}
	}

	if len(els) > 0 && els[0].Tag == "h1" {
		d.Heading = els[0]
		d.SearchTerm = els[0].Context[ContextSearchTerm]
	} else {
		// special h0 for fake heading based on file name - this is to ensure you don't get a heading of equal value within the document
		// fakeHeading :=
		d.Heading = &Element{
			Tag:     "h0",
			Type:    ElementTypeHeading,
			Context: map[string]string{ContextSearchTerm: filename, ContextHval: "0"},
			Content: []*ContentSegment{&ContentSegment{
				Raw:         filename,
				Attribution: AttributionPlain,
			}},
		}
		d.SearchTerm = filename
	}

	d.Elements = els
	if len(els) > 1 {
		d.SubDocuments = extractDocuments(els[1:], &d)
	}
	return &d
}

func zipDocumentAgain(els []*Element, i int) (*Document, int) {
	if len(els) <= i || els[i].Type != ElementTypeHeading {
		return nil, i + 1
	}
	el := els[i]
	hval, _ := strconv.Atoi(el.Context[ContextHval])
	var j int = i + 1
	for ; j < len(els); j++ {
		el2 := els[j]
		if el2.Type == ElementTypeHeading {
			hval2, _ := strconv.Atoi(el2.Context[ContextHval])
			if hval2 <= hval {
				//reached the end
				break
			}
		}
	}

	thisDocEls := els[i:j]
	doc := Document{
		Heading:    el,
		Elements:   thisDocEls,
		SearchTerm: el.Context[ContextSearchTerm],
	}
	doc.SubDocuments = extractDocuments(thisDocEls[1:], &doc)
	return &doc, j
}

func extractDocuments(els []*Element, doc *Document) []*Document {
	res := make([]*Document, 0)
	for i := 0; i < len(els); {
		d, j := zipDocumentAgain(els, i)
		if d != nil {
			d.Super = doc
			res = append(res, d)
			i = j
		} else {
			i++
		}
	}
	return res
}

func parseElement(n *html.Node, includingText bool) *Element {
	var e *Element = nil
	if n.Type == html.ElementNode {
		switch n.Data {
		case "p":
			e = &Element{}
			e.Type = ElementTypeString
			e.Content = parseContent(n)
			e.Tag = "p"
		case "h1", "h2", "h3", "h4", "h5", "h6":
			e = &Element{}
			e.Tag = n.Data
			hvalue := strings.TrimLeft(n.Data, "h")
			e.Type = ElementTypeHeading
			e.Content = parseContent(n)
			plain := parsePlainContent(n)[0]
			context := map[string]string{
				ContextSearchTerm: strings.Trim(plain.Raw, " \n"),
				ContextHval:       hvalue,
			}
			e.Context = context
		case "pre":
			children := childElements(n)
			if len(children) == 1 && children[0].Data == "code" {
				code := children[0]
				e = &Element{}
				e.Tag = code.Data
				e.Type = ElementTypeCode
				e.Content = parsePlainContent(code)
			}
		case "blockquote":
			children := childElements(n)
			if len(children) == 1 && children[0].Data == "p" {
				pEl := children[0]
				e = &Element{}
				e.Tag = "blockquote"
				e.Type = ElementTypeQuote
				e.Content = parsePlainContent(pEl)
			}
		case "ul":
			children := childElements(n)
			e = &Element{}
			e.Tag = "ul"
			e.Type = ElementTypeUnorderedList
			items := make([]*Element, 0)
			for _, c := range children {
				if c.Data == "li" {
					if el := parseElement(c, true); el != nil {
						items = append(items, el)
					}
				}
			}
			e.SubElements = items
		case "li":
			e = &Element{}
			e.Tag = "li"
			e.Type = ElementTypeListItem
			e.Content = parseContent(n)
			subE := make([]*Element, 0)
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if childE := parseElement(c, true); childE != nil {
					subE = append(subE, childE)
				}
			}
			e.SubElements = subE
		case "hr":
			e = &Element{}
			e.Tag = "hr"
			e.Type = ElementTypeHorizontalRule
		default:
			log.Printf("Whaaat? %s\n", n.Data)
		}
	} else if includingText && n.Type == html.TextNode {
		// we parse this element as if it were a <p>.
		// this will be the case for parsing <li> content with only text content
		e = &Element{}
		e.Tag = "p"
		e.Type = ElementTypeString
		e.Content = parseContent(n)
	}
	return e
}

func childElements(node *html.Node) []*html.Node {
	out := make([]*html.Node, 0)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			out = append(out, c)
		}
	}
	return out
}

// returns the string content of a node recursively
func parseContent(n *html.Node) []*ContentSegment {
	segments := make([]*ContentSegment, 0)

	var parser func(*html.Node, Attribution)
	parser = func(node *html.Node, attribution Attribution) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "em":
				attribution = attribution | AttributionEmphasis
			case "strong":
				attribution = attribution | AttributionBold
			case "code":
				attribution = attribution | AttributionCode
			case "a":
				attribution |= AttributeAnchor
				seg := ContentSegment{
					Raw:         node.FirstChild.Data,
					Attribution: attribution,
				}
				seg.Context = make(map[string]string)
				for _, att := range node.Attr {
					if att.Key == "href" {
						seg.Context["href"] = att.Val
						break
					}
				}
				segments = append(segments, &seg)
				return
			}
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				parser(c, attribution)
			}
		} else if node.Type == html.TextNode {
			seg := ContentSegment{
				Raw:         node.Data,
				Attribution: attribution,
			}
			segments = append(segments, &seg)
		}
	}
	parser(n, AttributionPlain)

	// for _, l := range segments {
	// 	log.Printf("Segment %v\n", l)
	// }

	return segments
}

func parsePlainContent(n *html.Node) []*ContentSegment {
	segment := ContentSegment{}
	rawStr := ""
	var parser func(*html.Node)
	parser = func(node *html.Node) {
		if node.Type == html.TextNode {
			rawStr = rawStr + node.Data
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				parser(c)
			}
		}
	}
	parser(n)
	segment.Raw = rawStr
	segment.Attribution = AttributionPlain
	return []*ContentSegment{&segment}
}

func parseContentSegment(n *html.Node) []ContentSegment {
	segs := make([]ContentSegment, 0)

	return segs
}

func (doc *Document) TraverseQuery(q string) *Document {
	st := doc.SearchTerm
	if strings.HasPrefix(q, st) {
		remaining := strings.TrimLeft(" ", strings.TrimPrefix(q, st))
		if remaining != "" {
			return doc
		}
		for _, sub := range doc.SubDocuments {
			dd := sub.TraverseQuery(remaining)
			if dd != nil {
				return dd
			}
		}
	}
	return nil
}

func (doc *Document) SubQueries() [][]string {
	log.Println(doc.Heading.Context)
	st := doc.SearchTerm
	// "this" heading is one query
	// then prepend to all sub-doc queries
	res := [][]string{[]string{st}}
	for _, d := range doc.SubDocuments {
		for _, qq := range d.SubQueries() {
			qq = append([]string{st}, qq...)
			res = append(res, qq)
		}
	}
	return res
}
