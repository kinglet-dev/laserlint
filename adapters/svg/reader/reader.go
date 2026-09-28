// Package reader turns an SVG document into a design measured in millimetres.
// The document is untrusted: its size, nesting and complexity are bounded,
// and nothing outside it (entities, links, images) is ever loaded.
package reader

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// context is what an element inherits from its ancestors: the matrix from
// its user units to millimetres, its paint, and whether it is being skipped.
type context struct {
	toMM  geom.Matrix
	paint paint
	skip  bool
}

// walker reads a document element by element, keeping a stack of contexts
// that is pushed at each start tag and popped at each end tag.
type walker struct {
	design   design.Design
	stack    []context
	limits   Limits
	elements int
	pieces   int
}

// Read parses an SVG document within DefaultLimits.
func Read(r io.Reader) (design.Design, error) { return ReadLimited(r, DefaultLimits) }

// ReadLimited parses an SVG document within the given limits.
func ReadLimited(r io.Reader, limits Limits) (design.Design, error) {
	dec := xml.NewDecoder(&limitedReader{r: r, left: limits.Bytes, limit: limits.Bytes})
	w := &walker{limits: limits}
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) && w.elements == 0 {
			return design.Design{}, fmt.Errorf("%w: it has no elements", ErrNotSVG)
		}
		if errors.Is(err, io.EOF) {
			return w.design, nil
		}
		if err != nil {
			return design.Design{}, badXML(err)
		}
		if err := w.token(tok); err != nil {
			return design.Design{}, err
		}
	}
}

// badXML explains a decoder error; hitting the size limit is reported as is.
func badXML(err error) error {
	if errors.Is(err, ErrTooLarge) {
		return err
	}
	return fmt.Errorf("%w (%v)", ErrBadXML, err)
}

// token handles one XML token.
func (w *walker) token(tok xml.Token) error {
	switch t := tok.(type) {
	case xml.StartElement:
		return w.element(t)
	case xml.EndElement:
		w.stack = w.stack[:len(w.stack)-1]
	case xml.Directive:
		return directive(t)
	}
	return nil
}

// element reads, skips or refuses one element.
func (w *walker) element(el xml.StartElement) error {
	if err := w.count(); err != nil {
		return err
	}
	if len(w.stack) > 0 && w.stack[len(w.stack)-1].skip {
		w.stack = append(w.stack, context{skip: true})
		return nil
	}
	skip, err := w.classify(el)
	if err != nil || skip {
		w.stack = append(w.stack, context{skip: true})
		return err
	}
	return w.start(el)
}
