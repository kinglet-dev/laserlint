package reader

import "strings"

// paint is the fill and stroke in force for an element, after inheritance.
type paint struct {
	fill, stroke string
}

// initialPaint is SVG's starting point: black fill, no stroke.
var initialPaint = paint{fill: "black", stroke: "none"}

// with returns the paint for an element: the inherited paint, overridden
// by the element's fill and stroke attributes, then by its style attribute.
func (p paint) with(a attributes) paint {
	p.set("fill", a["fill"])
	p.set("stroke", a["stroke"])
	for _, decl := range strings.Split(a["style"], ";") {
		name, value, _ := strings.Cut(decl, ":")
		p.set(strings.ToLower(strings.TrimSpace(name)), value)
	}
	return p
}

// set applies one property; an empty value or "inherit" keeps the parent's.
func (p *paint) set(name, value string) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || value == "inherit" {
		return
	}
	switch name {
	case "fill":
		p.fill = value
	case "stroke":
		p.stroke = value
	}
}

// visible reports whether the laser software would score the shape at all.
func (p paint) visible() bool { return p.filled() || p.stroke != "none" }

// filled reports whether the shape's outline is scored as a filled area.
func (p paint) filled() bool { return p.fill != "none" }

// white reports a white fill, which usually means a background shape.
// Values are already lower case; spaces inside rgb() don't matter.
func (p paint) white() bool {
	switch strings.ReplaceAll(p.fill, " ", "") {
	case "#fff", "#ffffff", "white", "rgb(255,255,255)":
		return true
	}
	return false
}
