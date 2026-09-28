package syntax

import (
	"strings"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// ParseNumber reads a shape attribute such as a rect's width: exactly one
// number in user units, with optional surrounding whitespace.
func ParseNumber(value string) (float64, error) {
	s := &scanner{d: strings.TrimSpace(value)}
	n, err := s.number()
	if err != nil {
		return 0, err
	}
	if s.pos < len(s.d) {
		if isUnitStart(s.d[s.pos]) {
			return 0, s.fail(ErrUnitInAttribute, s.pos)
		}
		return 0, s.fail(ErrBadNumber, s.pos)
	}
	return n, nil
}

// DefaultMaxPoints bounds ParsePoints, matching the limit on path pieces.
const DefaultMaxPoints = DefaultMaxPieces

// ParsePoints reads the points of a polyline or polygon: x,y pairs
// separated by whitespace and commas.
func ParsePoints(value string) ([]geom.Point, error) {
	return ParsePointsLimited(value, DefaultMaxPoints)
}

// ParsePointsLimited is ParsePoints with a limit on the number of points.
func ParsePointsLimited(value string, maxPoints int) ([]geom.Point, error) {
	s := &scanner{d: value}
	var points []geom.Point
	for !s.done() {
		if len(points) == maxPoints {
			return nil, s.fail(ErrTooComplex, s.pos)
		}
		x, err := s.number()
		if err != nil {
			return nil, err
		}
		if s.done() {
			return nil, ErrOddCoordinates
		}
		y, err := s.number()
		if err != nil {
			return nil, err
		}
		points = append(points, geom.Point{X: x, Y: y})
	}
	return points, nil
}

func isUnitStart(c byte) bool {
	return c == '%' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
