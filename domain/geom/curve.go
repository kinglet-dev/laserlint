package geom

import "math"

// maxSubdivisions bounds the work for pathological curves or tolerances: a
// segment yields at most 2^16 pieces. Flatness shrinks about fourfold per
// split (4^16 ≈ 4.3 × 10^9), so 16 levels reach 0.01 mm on curves far larger
// than any laser bed.
const maxSubdivisions = 16

// appendCubic adds points along the cubic from p0 (already in out) to p3,
// splitting it in half until each piece is within tolerance of its chord.
func appendCubic(out []Point, p0, c1, c2, p3 Point, tolerance float64, depth int) []Point {
	if depth >= maxSubdivisions || flatness(p0, c1, c2, p3) <= tolerance {
		return append(out, p3)
	}
	// de Casteljau split at t = 0.5
	m01, m12, m23 := mid(p0, c1), mid(c1, c2), mid(c2, p3)
	m012, m123 := mid(m01, m12), mid(m12, m23)
	m := mid(m012, m123)
	out = appendCubic(out, p0, m01, m012, m, tolerance, depth+1)
	return appendCubic(out, m, m123, m23, p3, tolerance, depth+1)
}

// flatness is the largest distance of the control points from the chord
// p0–p3. The curve lies inside its control polygon, so it deviates from the
// chord by no more than this.
func flatness(p0, c1, c2, p3 Point) float64 {
	return math.Max(distanceToSegment(c1, p0, p3), distanceToSegment(c2, p0, p3))
}

func distanceToSegment(q, a, b Point) float64 {
	dx, dy := b.X-a.X, b.Y-a.Y
	t := 0.0
	if l := dx*dx + dy*dy; l > 0 {
		t = math.Max(0, math.Min(1, ((q.X-a.X)*dx+(q.Y-a.Y)*dy)/l))
	}
	return math.Hypot(q.X-(a.X+t*dx), q.Y-(a.Y+t*dy))
}

func mid(a, b Point) Point { return Point{X: (a.X + b.X) / 2, Y: (a.Y + b.Y) / 2} }
