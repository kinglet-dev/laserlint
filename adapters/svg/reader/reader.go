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

// context is what an element inherits from its ancestors: the matrix from
// its user units to millimetres, and its paint.
type context struct {
	toMM  geom.Matrix
	paint paint
}

// walker reads a document element by element, keeping a stack of contexts
// that is pushed at each start tag and popped at each end tag.
type walker struct {
	design design.Design
	stack  []context
}

// Read parses an SVG document.
func Read(r io.Reader) (design.Design, error) {
	dec := xml.NewDecoder(r)
	w := &walker{}
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return w.design, nil
		}
		if err != nil {
			return design.Design{}, fmt.Errorf("%w (%v)", ErrBadXML, err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if err := w.start(t); err != nil {
				return design.Design{}, err
			}
		case xml.EndElement:
			w.stack = w.stack[:len(w.stack)-1]
		}
	}
}

// start works out an element's context, pushes it, and records any shape.
func (w *walker) start(el xml.StartElement) error {
	a := attrs(el)
	ctx, err := w.context(el, a)
	if err != nil {
		return err
	}
	w.stack = append(w.stack, ctx)
	kind, ok := shapeKinds[el.Name.Local]
	if !ok {
		return nil
	}
	path, drawn, err := kind.build(a)
	if err != nil {
		return fmt.Errorf("%s: %w", describe(el), err)
	}
	if !kind.fillable {
		ctx.paint.fill = "none"
	}
	if drawn {
		w.addShape(path, ctx)
	}
	return nil
}

// context combines an element's own attributes with what it inherits. The
// root svg element sets the millimetre scale from the document's size.
func (w *walker) context(el xml.StartElement, a attributes) (context, error) {
	if len(w.stack) == 0 {
		var ctx context
		var err error
		if w.design.Width, w.design.Height, ctx.toMM, err = size(a); err != nil {
			return context{}, err
		}
		ctx.paint = initialPaint.with(a)
		return ctx, nil
	}
	parent := w.stack[len(w.stack)-1]
	ctx := context{toMM: parent.toMM, paint: parent.paint.with(a)}
	if list, ok := a["transform"]; ok {
		m, err := syntax.ParseTransform(list)
		if err != nil {
			return context{}, fmt.Errorf("%s: %w", describe(el), err)
		}
		ctx.toMM = m.Then(parent.toMM)
	}
	return ctx, nil
}

// addShape records a shape in millimetres with what the laser will score.
// Shapes with no fill and no stroke are dropped: nothing is scored.
func (w *walker) addShape(path geom.Path, ctx context) {
	if !ctx.paint.visible() {
		return
	}
	w.design.Shapes = append(w.design.Shapes, design.Shape{
		Path:      path.Transform(ctx.toMM),
		Filled:    ctx.paint.filled(),
		WhiteFill: ctx.paint.white(),
	})
}
