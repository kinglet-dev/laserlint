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

// row is the height of a bucket of edges, mm.
const row = 0.5

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
	winding := map[int]int{}
	for _, e := range f.rows[key(p.Y)] {
		if (e.a.Y <= p.Y) == (e.b.Y <= p.Y) {
			continue // doesn't cross this height
		}
		if x := e.a.X + (p.Y-e.a.Y)*(e.b.X-e.a.X)/(e.b.Y-e.a.Y); x > p.X {
			if e.b.Y > e.a.Y {
				winding[e.shape]++
			} else {
				winding[e.shape]--
			}
		}
	}
	top := -1
	for s, w := range winding {
		if w != 0 && s > top {
			top = s
		}
	}
	return top >= 0 && !f.white[top]
}
