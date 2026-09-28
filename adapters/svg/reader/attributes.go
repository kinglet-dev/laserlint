package reader

import (
	"encoding/xml"
	"fmt"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

// attributes maps an element's attribute names (without namespace) to values.
type attributes map[string]string

// attrs collects an element's SVG attributes, leaving out those in other
// namespaces (an editor's own settings).
func attrs(el xml.StartElement) attributes {
	m := make(attributes, len(el.Attr))
	for _, a := range el.Attr {
		if a.Name.Space == "" {
			m[a.Name.Local] = a.Value
		}
	}
	return m
}

// numbers reads plain-number attributes; a missing one is 0, as in SVG.
func (a attributes) numbers(names ...string) ([]float64, error) {
	out := make([]float64, len(names))
	for i, name := range names {
		v, ok := a[name]
		if !ok {
			continue
		}
		n, err := syntax.ParseNumber(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		out[i] = n
	}
	return out, nil
}

// radii reads rx and ry the SVG 2 way: a missing or negative radius takes
// the other one's value, and both default to 0.
func (a attributes) radii() (rx, ry float64, err error) {
	v, err := a.numbers("rx", "ry")
	if err != nil {
		return 0, 0, err
	}
	rxSet, rySet := a.has("rx") && v[0] >= 0, a.has("ry") && v[1] >= 0
	rx, ry = max(v[0], 0), max(v[1], 0)
	if !rxSet {
		rx = ry
	}
	if !rySet {
		ry = rx
	}
	return rx, ry, nil
}

func (a attributes) has(name string) bool {
	_, ok := a[name]
	return ok
}

// describe names an element for error messages, with its id when it has one.
func describe(el xml.StartElement) string {
	if id := attrs(el)["id"]; id != "" {
		return fmt.Sprintf("<%s id=%q>", el.Name.Local, id)
	}
	return "<" + el.Name.Local + ">"
}
