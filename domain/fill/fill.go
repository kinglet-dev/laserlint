// Package fill tells what colour a design shows at a point: the colour of
// the topmost filled shape containing it. Containment follows SVG's default
// nonzero rule, counted with the winding number (Sunday, "Inclusion of a
// point in a polygon", 2001; Hormann and Agathos 2001), and edges are
// bucketed by row so a query only looks at edges at its height.
package fill

import (
	"errors"
	"fmt"
	"math"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// row is the height of a bucket of edges, mm: short, since traced art has
// many short edges and a query looks at every edge in its row.
const row = 0.1

// ErrTooLarge means the shapes' edges span more rows than the limit.
var ErrTooLarge = errors.New("the design is too large to work out its fill")

type edge struct {
	a, b  geom.Point
	shape int
}

// Fill answers colour queries for one design.
type Fill struct {
	white []bool
	rows  map[int][]edge
}

// New indexes the filled shapes, following curves within tolerance, mm.
// Shapes later in the list are painted over earlier ones.
func New(shapes []design.Shape, tolerance float64, maxEntries int) (*Fill, error) {
	f := &Fill{white: make([]bool, len(shapes)), rows: map[int][]edge{}}
	entries := 0
	for i, s := range shapes {
		f.white[i] = s.WhiteFill
		if !s.Filled {
			continue
		}
		for _, l := range s.Path.Flatten(tolerance) {
			pts := append(l.Points[:len(l.Points):len(l.Points)], l.Points[0])
			for j := 1; j < len(pts); j++ {
				e := edge{a: pts[j-1], b: pts[j], shape: i}
				lo, hi := key(math.Min(e.a.Y, e.b.Y)), key(math.Max(e.a.Y, e.b.Y))
				if entries += hi - lo + 1; entries > maxEntries {
					return nil, fmt.Errorf("%w: its edges span more than %d rows of %g mm", ErrTooLarge, maxEntries, row)
				}
				for r := lo; r <= hi; r++ {
					f.rows[r] = append(f.rows[r], e)
				}
			}
		}
	}
	return f, nil
}

func key(y float64) int { return int(math.Floor(y / row)) }

// Dark reports whether the topmost filled shape at p is not white.
func (f *Fill) Dark(p geom.Point) bool {
	var winding []count // few shapes cross any one ray, so a list beats a map
	for _, e := range f.rows[key(p.Y)] {
		if (e.a.Y <= p.Y) == (e.b.Y <= p.Y) {
			continue // doesn't cross this height
		}
		if x := e.a.X + (p.Y-e.a.Y)*(e.b.X-e.a.X)/(e.b.Y-e.a.Y); x > p.X {
			if e.b.Y > e.a.Y {
				winding = add(winding, e.shape, 1)
			} else {
				winding = add(winding, e.shape, -1)
			}
		}
	}
	top := -1
	for _, c := range winding {
		if c.w != 0 && c.shape > top {
			top = c.shape
		}
	}
	return top >= 0 && !f.white[top]
}

// count is one shape's winding number around a point.
type count struct{ shape, w int }

func add(counts []count, shape, w int) []count {
	for i := range counts {
		if counts[i].shape == shape {
			counts[i].w += w
			return counts
		}
	}
	return append(counts, count{shape, w})
}
