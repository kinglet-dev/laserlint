package reader

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Limits bound the work a document can cause, so a hostile or runaway
// file fails with a clear message instead of exhausting memory or time.
type Limits struct {
	Bytes    int64 // size of the file
	Depth    int   // levels of nested elements, the root included
	Elements int   // elements in total, skipped ones included
	Pieces   int   // straight and curved path pieces in total, over all shapes
}

// DefaultLimits are the limits in laserlint's threat model.
var DefaultLimits = Limits{Bytes: 25_000_000, Depth: 256, Elements: 200_000, Pieces: 2_000_000}

// maxMM bounds every coordinate once converted to millimetres (10 km).
const maxMM = 10_000_000

// limitedReader fails with ErrTooLarge once more than left bytes are read.
type limitedReader struct {
	r     io.Reader
	left  int64
	limit int64
}

func (l *limitedReader) Read(p []byte) (int, error) {
	n, err := l.r.Read(p)
	if l.left -= int64(n); l.left < 0 {
		return 0, fmt.Errorf("%w: the limit is %g MB", ErrTooLarge, float64(l.limit)/1e6)
	}
	return n, err
}

// count records one more element and checks the element and depth limits.
func (w *walker) count() error {
	if w.elements++; w.elements > w.limits.Elements {
		return fmt.Errorf("%w: more than %d elements", ErrTooManyElements, w.limits.Elements)
	}
	if len(w.stack) >= w.limits.Depth {
		return fmt.Errorf("%w: more than %d levels", ErrTooDeep, w.limits.Depth)
	}
	return nil
}

// budget is how many path pieces the remaining shapes may still use.
func (w *walker) budget() int { return w.limits.Pieces - w.pieces }

// spend records a shape's pieces, failing once the total is over the limit.
func (w *walker) spend(p geom.Path) error {
	for _, s := range p.Subpaths {
		w.pieces += len(s.Segments)
	}
	if w.pieces > w.limits.Pieces {
		return fmt.Errorf("%w: the limit is %d in total", syntax.ErrTooComplex, w.limits.Pieces)
	}
	return nil
}

// directive refuses a DOCTYPE that defines entities.
func directive(d xml.Directive) error {
	if bytes.Contains(d, []byte("<!ENTITY")) {
		return ErrEntities
	}
	return nil
}

// inRange reports whether every point of a path in millimetres is finite
// and within maxMM.
func inRange(p geom.Path) bool {
	ok := func(q geom.Point) bool { return math.Abs(q.X) <= maxMM && math.Abs(q.Y) <= maxMM }
	for _, s := range p.Subpaths {
		if !ok(s.Start) {
			return false
		}
		for _, seg := range s.Segments {
			if !ok(seg.To) || !ok(seg.C1) || !ok(seg.C2) {
				return false
			}
		}
	}
	return true
}
