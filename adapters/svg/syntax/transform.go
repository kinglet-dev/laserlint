package syntax

import (
	"errors"
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// transformFunc builds a matrix from a transform function's arguments.
type transformFunc struct {
	counts []int // allowed argument counts
	build  func(a []float64) geom.Matrix
}

// transformFuncs lists the SVG transform functions (Strategy pattern, like
// the path command table).
var transformFuncs = map[string]transformFunc{
	"matrix": {[]int{6}, func(a []float64) geom.Matrix {
		return geom.Matrix{A: a[0], B: a[1], C: a[2], D: a[3], E: a[4], F: a[5]}
	}},
	"translate": {[]int{1, 2}, func(a []float64) geom.Matrix { return translate(a[0], second(a, 0)) }},
	"scale": {[]int{1, 2}, func(a []float64) geom.Matrix {
		return geom.Matrix{A: a[0], D: second(a, a[0])}
	}},
	"rotate": {[]int{1, 3}, func(a []float64) geom.Matrix {
		cx, cy := 0.0, 0.0
		if len(a) == 3 {
			cx, cy = a[1], a[2]
		}
		sin, cos := math.Sincos(a[0] * math.Pi / 180)
		turn := geom.Matrix{A: cos, B: sin, C: -sin, D: cos}
		return translate(-cx, -cy).Then(turn).Then(translate(cx, cy))
	}},
	"skewX": {[]int{1}, func(a []float64) geom.Matrix {
		return geom.Matrix{A: 1, C: math.Tan(a[0] * math.Pi / 180), D: 1}
	}},
	"skewY": {[]int{1}, func(a []float64) geom.Matrix {
		return geom.Matrix{A: 1, B: math.Tan(a[0] * math.Pi / 180), D: 1}
	}},
}

// ParseTransform reads an SVG transform list such as "translate(10 0) scale(2)".
// The functions apply right to left: the last one acts on a point first.
func ParseTransform(list string) (geom.Matrix, error) {
	s := &scanner{d: list}
	result := geom.Identity()
	for !s.done() {
		at := s.pos
		name := s.name()
		fn, known := transformFuncs[name]
		if !known {
			return geom.Matrix{}, s.fail(ErrUnknownTransform, at)
		}
		args, err := s.arguments()
		if err != nil {
			return geom.Matrix{}, err
		}
		if !allowed(fn.counts, len(args)) {
			return geom.Matrix{}, s.fail(ErrTransformArgs, at)
		}
		result = fn.build(args).Then(result)
		if !finite(result) {
			return geom.Matrix{}, s.fail(ErrOutOfRange, at)
		}
	}
	return result, nil
}

// finite reports whether every entry of m is a finite number.
func finite(m geom.Matrix) bool {
	for _, v := range [...]float64{m.A, m.B, m.C, m.D, m.E, m.F} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

// name reads a run of letters.
func (s *scanner) name() string {
	start := s.pos
	for s.pos < len(s.d) && (s.d[s.pos]|0x20 >= 'a' && s.d[s.pos]|0x20 <= 'z') {
		s.pos++
	}
	return s.d[start:s.pos]
}

// arguments reads "(numbers)".
func (s *scanner) arguments() ([]float64, error) {
	s.skipSeparators()
	if s.pos >= len(s.d) || s.d[s.pos] != '(' {
		return nil, s.fail(ErrTransformSyntax, s.pos)
	}
	s.pos++
	var args []float64
	for s.numberNext() {
		v, err := s.number()
		if errors.Is(err, ErrBadNumber) {
			return nil, s.fail(ErrTransformSyntax, s.pos)
		}
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	if s.pos >= len(s.d) || s.d[s.pos] != ')' {
		return nil, s.fail(ErrTransformSyntax, s.pos)
	}
	s.pos++
	return args, nil
}

func translate(x, y float64) geom.Matrix { return geom.Matrix{A: 1, D: 1, E: x, F: y} }

// second returns the second argument, or fallback when there is only one.
func second(a []float64, fallback float64) float64 {
	if len(a) > 1 {
		return a[1]
	}
	return fallback
}

func allowed(counts []int, n int) bool {
	for _, c := range counts {
		if c == n {
			return true
		}
	}
	return false
}
