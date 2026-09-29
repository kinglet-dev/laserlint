package crossing

import (
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// eps is how near, in mm, a point must be to a line to count as on it:
// far below anything a laser can resolve, far above rounding error.
const eps = 1e-6

// relation is how two segments meet.
type relation struct {
	points    []geom.Point // crossings
	overlap   float64      // length lying on each other
	overlapAt geom.Point
}

// relate classifies a pair of segments with orientation tests.
func relate(p, q segment) relation {
	dc, dd := side(p, q.a), side(p, q.b)
	if math.Abs(dc) < eps && math.Abs(dd) < eps {
		return collinear(p, q)
	}
	var r relation
	for _, t := range []struct {
		s  segment
		pt geom.Point
		d  float64
	}{{p, q.a, dc}, {p, q.b, dd}, {q, p.a, side(q, p.a)}, {q, p.b, side(q, p.b)}} {
		if math.Abs(t.d) < eps && inside(t.s, t.pt) {
			r.points = append(r.points, t.pt) // one ends on the other's middle: a T
		}
	}
	da, db := side(q, p.a), side(q, p.b)
	if dc*dd < 0 && da*db < 0 {
		// A proper crossing. An end within eps of the other piece is also
		// found as a T above; merge folds the two into one.
		t := dc / (dc - dd)
		r.points = append(r.points, lerp(q.a, q.b, t))
	}
	return r
}

// side is the signed distance from pt to the line through s.
func side(s segment, pt geom.Point) float64 {
	return ((s.b.X-s.a.X)*(pt.Y-s.a.Y) - (s.b.Y-s.a.Y)*(pt.X-s.a.X)) / length(s)
}

// inside reports whether pt, on s's line, lies strictly between s's ends.
func inside(s segment, pt geom.Point) bool {
	t := along(s, pt)
	return t > eps && t < length(s)-eps
}

// along is how far pt's projection lies along s from its start, mm.
func along(s segment, pt geom.Point) float64 {
	return ((pt.X-s.a.X)*(s.b.X-s.a.X) + (pt.Y-s.a.Y)*(s.b.Y-s.a.Y)) / length(s)
}

// collinear measures how much of two segments on one line lie on each other.
func collinear(p, q segment) relation {
	t := []float64{along(p, q.a), along(p, q.b)}
	sort.Float64s(t)
	lo, hi := math.Max(0, t[0]), math.Min(length(p), t[1])
	if hi-lo <= eps {
		return relation{}
	}
	mid := (lo + hi) / 2 / length(p)
	return relation{overlap: hi - lo, overlapAt: lerp(p.a, p.b, mid)}
}

func length(s segment) float64 { return math.Hypot(s.b.X-s.a.X, s.b.Y-s.a.Y) }

func lerp(a, b geom.Point, t float64) geom.Point {
	return geom.Point{X: a.X + t*(b.X-a.X), Y: a.Y + t*(b.Y-a.Y)}
}
