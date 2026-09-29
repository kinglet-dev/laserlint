package spacing

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// nearest is the distance from a piece's midpoint to the nearest other
// scored line within reach, or +Inf when there is none. Line on the
// same line within the window either side (measured along the line) is
// not "other" line.
func (idx *index) nearest(pc piece) float64 {
	p := pc.mid
	at := pc.arc + pc.length/2
	best := math.Inf(1)
	for x := idx.key(p.X) - 1; x <= idx.key(p.X)+1; x++ {
		for y := idx.key(p.Y) - 1; y <= idx.key(p.Y)+1; y++ {
			for _, id := range idx.cells[[2]int{x, y}] {
				if q := &idx.pieces[id]; !idx.skip(p, at, pc.line, q) {
					best = math.Min(best, idx.distance(p, at, pc.line, *q))
				}
			}
		}
	}
	return best
}

// skip is a cheap early exit: q is out of reach even at its ends, or lies
// wholly inside the same-line window (on a closed line the window only
// grows by wrapping round, so this holds there too).
func (idx *index) skip(p geom.Point, at float64, line int, q *piece) bool {
	dx, dy, r := q.mid.X-p.X, q.mid.Y-p.Y, idx.reach+q.length/2
	if dx*dx+dy*dy >= r*r {
		return true
	}
	return q.line == line && q.arc >= at-idx.window && q.arc+q.length <= at+idx.window
}

// distance is the nearest distance from p to the part of q within reach
// that counts as other line, or +Inf.
func (idx *index) distance(p geom.Point, at float64, line int, q piece) float64 {
	u0, u1, ok := within(p, q.a, q.b, idx.reach)
	if !ok {
		return math.Inf(1)
	}
	parts := [][2]float64{{q.arc + u0*q.length, q.arc + u1*q.length}}
	if q.line == line {
		parts = idx.outsideWindow(parts, at, idx.lines[line])
	}
	best := math.Inf(1)
	for _, part := range parts {
		u := clamp(closestU(p, q.a, q.b), (part[0]-q.arc)/q.length, (part[1]-q.arc)/q.length)
		c := lerp(q.a, q.b, u)
		best = math.Min(best, math.Hypot(c.X-p.X, c.Y-p.Y))
	}
	return best
}

// outsideWindow removes the stretch of the same line around at, going
// round the loop for a closed line.
func (idx *index) outsideWindow(parts [][2]float64, at float64, l lineInfo) [][2]float64 {
	shifts := []float64{0}
	if l.closed {
		shifts = []float64{-l.perimeter, 0, l.perimeter}
	}
	for _, s := range shifts {
		parts = minus(parts, at+s-idx.window, at+s+idx.window)
	}
	return parts
}

// minus removes (lo, hi) from each interval.
func minus(parts [][2]float64, lo, hi float64) [][2]float64 {
	var out [][2]float64
	for _, p := range parts {
		if p[0] < lo {
			out = append(out, [2]float64{p[0], math.Min(p[1], lo)})
		}
		if p[1] > hi {
			out = append(out, [2]float64{math.Max(p[0], hi), p[1]})
		}
	}
	return out
}

// within finds the part [u0, u1] of segment ab (0 ≤ u ≤ 1) nearer to p than r.
func within(p, a, b geom.Point, r float64) (u0, u1 float64, ok bool) {
	vx, vy, wx, wy := b.X-a.X, b.Y-a.Y, a.X-p.X, a.Y-p.Y
	qa, qb, qc := vx*vx+vy*vy, 2*(vx*wx+vy*wy), wx*wx+wy*wy-r*r
	disc := qb*qb - 4*qa*qc
	if disc <= 0 {
		return 0, 0, false
	}
	u0 = math.Max(0, (-qb-math.Sqrt(disc))/(2*qa))
	u1 = math.Min(1, (-qb+math.Sqrt(disc))/(2*qa))
	return u0, u1, u0 < u1
}

// closestU is where on the infinite line through a and b p is nearest.
func closestU(p, a, b geom.Point) float64 {
	vx, vy := b.X-a.X, b.Y-a.Y
	return ((p.X-a.X)*vx + (p.Y-a.Y)*vy) / (vx*vx + vy*vy)
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
