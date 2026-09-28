// Package reader turns an SVG document into a design measured in millimetres.
// The document is untrusted: its size, nesting and complexity are bounded,
// and nothing outside it (entities, links, images) is ever loaded.
package reader

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Read parses an SVG document.
func Read(r io.Reader) (design.Design, error) {
	dec := xml.NewDecoder(r)
	var d design.Design
	var toMM geom.Matrix
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return d, nil
		}
		if err != nil {
			return design.Design{}, fmt.Errorf("%w (%v)", ErrBadXML, err)
		}
		el, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch el.Name.Local {
		case "svg":
			if d.Width, d.Height, toMM, err = size(attrs(el)); err != nil {
				return design.Design{}, err
			}
		case "path":
			path, err := syntax.ParsePath(attrs(el)["d"])
			if err != nil {
				return design.Design{}, fmt.Errorf("%s: %w", describe(el), err)
			}
			d.Shapes = append(d.Shapes, design.Shape{Path: path.Transform(toMM), Filled: true})
		}
	}
}

// attrs maps an element's attribute names to values.
func attrs(el xml.StartElement) map[string]string {
	m := make(map[string]string, len(el.Attr))
	for _, a := range el.Attr {
		m[a.Name.Local] = a.Value
	}
	return m
}

// describe names an element for error messages, with its id when it has one.
func describe(el xml.StartElement) string {
	if id := attrs(el)["id"]; id != "" {
		return fmt.Sprintf("<%s id=%q>", el.Name.Local, id)
	}
	return "<" + el.Name.Local + ">"
}
