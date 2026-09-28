// Package pathdata parses SVG path data (the d attribute) into geometry.
// Input is untrusted: the grammar is strict and the amount of work is bounded.
package pathdata

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// handler runs one command with its arguments; rel is true for lower-case (relative) commands.
type handler func(st *state, s *scanner, rel bool) error

// handlers maps each upper-case command letter to its handler (Strategy pattern):
// adding a command means adding an entry, not editing the parse loop.
var handlers = map[byte]handler{
	'M': moveTo, 'L': lineTo, 'H': horizontal, 'V': vertical, 'Z': closePath,
	'C': cubicTo, 'S': smoothCubicTo, 'Q': quadTo, 'T': smoothQuadTo,
}

// state is the pen position while parsing.
type state struct {
	path    geom.Path
	current geom.Point // where the pen is
	start   geom.Point // start of the current subpath, where Z returns to

	// The last curve's final control point, which S and T mirror.
	prevKind, lastKind byte // 'C' after C/S, 'Q' after Q/T, 0 otherwise
	prevCtrl, lastCtrl geom.Point
}

// Parse reads SVG path data such as "M 10 20 L 30 40 Z".
func Parse(d string) (geom.Path, error) {
	s := &scanner{d: d}
	st := &state{}
	for !s.done() {
		at := s.pos
		letter := s.command()
		upper, rel := letter&^0x20, letter >= 'a'
		run, known := handlers[upper]
		if len(st.path.Subpaths) == 0 && upper != 'M' {
			return geom.Path{}, s.fail(ErrNoMoveFirst, at)
		}
		if !known {
			return geom.Path{}, s.fail(ErrUnknownCommand, at)
		}
		for {
			st.prevKind, st.prevCtrl, st.lastKind = st.lastKind, st.lastCtrl, 0
			if err := run(st, s, rel); err != nil {
				return geom.Path{}, err
			}
			if upper == 'Z' || !s.numberNext() {
				break
			}
			if upper == 'M' { // extra pairs after a move are lines
				run = lineTo
			}
		}
	}
	return st.path, nil
}

// resolve turns a relative point into an absolute one and checks its range.
func (st *state) resolve(s *scanner, x, y float64, rel bool) (geom.Point, error) {
	if rel {
		x, y = st.current.X+x, st.current.Y+y
	}
	if math.Abs(x) > maxCoordinate || math.Abs(y) > maxCoordinate {
		return geom.Point{}, s.fail(ErrOutOfRange, s.pos-1)
	}
	return geom.Point{X: x, Y: y}, nil
}

// add appends a segment, starting a new subpath at the pen if the last one is closed.
func (st *state) add(seg geom.Segment) {
	n := len(st.path.Subpaths)
	if n == 0 || st.path.Subpaths[n-1].Closed {
		st.path.Subpaths = append(st.path.Subpaths, geom.Subpath{Start: st.current})
		n++
	}
	sub := &st.path.Subpaths[n-1]
	sub.Segments = append(sub.Segments, seg)
	st.current = seg.To
}
