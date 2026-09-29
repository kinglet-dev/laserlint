// Package crossing finds where score lines cross or lie on each other.
// Where lines cross, the laser burns the same spot twice; where a line is
// drawn on top of another, it burns the same stretch twice.
//
// Candidate pairs come from a uniform grid (spatial hash), and each pair is
// classified exactly with orientation predicates (Cormen et al.,
// "Introduction to Algorithms", section 33.1).
package crossing

import (
	"errors"
	"fmt"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Params sets the search.
type Params struct {
	Merge       float64 // crossings nearer than this are one crossing, mm
	Cell        float64 // grid cell size, mm
	MaxSegments int     // limits on work, so huge designs fail clearly
	MaxEntries  int     // limit on segment-in-cell entries
}

// Result lists what was found.
type Result struct {
	Crossings []geom.Point // where lines cross, or one ends on another
	Stacked   float64      // length of line lying on other line, mm
	StackedAt []geom.Point // middle of each stacked stretch
}

// ErrTooComplex means the design has too much line to search within the limits.
var ErrTooComplex = errors.New("the design has too much line to check for crossings")

// Find searches the lines. Pieces that meet only end to end are not
// crossings: neighbouring pieces of one line, and separate lines that join,
// as split paths do.
func Find(lines []geom.Polyline, p Params) (Result, error) {
	segs := collect(lines)
	if len(segs) > p.MaxSegments {
		return Result{}, fmt.Errorf("%w: %d line pieces, the limit is %d", ErrTooComplex, len(segs), p.MaxSegments)
	}
	g, err := newGrid(segs, p.Cell, p.MaxEntries)
	if err != nil {
		return Result{}, err
	}
	var r Result
	var hits []geom.Point
	g.eachPair(func(a, b segment) {
		rel := relate(a, b)
		hits = append(hits, rel.points...)
		if rel.overlap > 0 {
			r.Stacked += rel.overlap
			r.StackedAt = append(r.StackedAt, rel.overlapAt)
		}
	})
	r.Crossings = merge(hits, p.Merge)
	return r, nil
}
