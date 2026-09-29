// Package spacing measures how much scored line runs too close to other
// scored line. Laser lines closer than their own width plus a clear gap
// burn together, so fine detail turns into a dark smudge.
//
// A point is too close when scored line within MinDistance (centre to
// centre) belongs to another line, or to the same line after a turn tighter
// than a half circle of diameter MinDistance: the tightest U-turn that
// keeps the gap. So sharp corners count near their tip, while right angles
// and gentle curves never do. A closed loop smaller than MinDistance has no
// such turn and is left to the small-details check.
package spacing

import (
	"errors"
	"fmt"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Params sets the measurement.
type Params struct {
	MinDistance float64 // centre-to-centre distance below which lines are too close, mm
	Step        float64 // length of line each sample stands for, mm
	MaxSamples  int     // limit on samples, so huge designs fail clearly instead of running for ever
}

// Spot is one stretch of line that is too close to other line.
type Spot struct {
	At       geom.Point // where the other line is nearest
	Distance float64    // nearest centre-to-centre distance, mm
	Length   float64    // length of the stretch, mm
}

// Result is the measurement for a whole design.
type Result struct {
	Scored float64 // total scored length, mm
	Close  float64 // scored length that is too close, mm
	Spots  []Spot  // nearest first
}

// ErrTooManySamples means the design has too much line to measure within the limit.
var ErrTooManySamples = errors.New("the design has too much line to measure")

// Measure finds scored line that is too close to other scored line.
func Measure(lines []geom.Polyline, p Params) (Result, error) {
	idx, err := build(lines, p)
	if err != nil {
		return Result{}, err
	}
	var r Result
	near := make([]float64, len(idx.pieces))
	for i, pc := range idx.pieces {
		r.Scored += pc.length
		near[i], _ = idx.nearest(pc)
		if near[i] < p.MinDistance {
			r.Close += pc.length
		}
	}
	r.Spots = spots(idx.pieces, near, p.MinDistance)
	return r, nil
}

// build indexes the lines, refusing more samples than the limit.
func build(lines []geom.Polyline, p Params) (*index, error) {
	if n := sampleCount(lines, p.Step); n > p.MaxSamples {
		return nil, fmt.Errorf("%w: it needs %d samples, the limit is %d", ErrTooManySamples, n, p.MaxSamples)
	}
	return newIndex(lines, p.Step, p.MinDistance), nil
}

// Facing calls visit for every sample of line with other line within
// MinDistance, the same rule as Measure: from is the sample, to the
// nearest point of other line, and length how much line the sample stands for.
func Facing(lines []geom.Polyline, p Params, visit func(from, to geom.Point, length float64)) error {
	idx, err := build(lines, p)
	if err != nil {
		return err
	}
	for _, pc := range idx.pieces {
		if d, q := idx.nearest(pc); d < p.MinDistance {
			visit(pc.mid, q, pc.length)
		}
	}
	return nil
}
