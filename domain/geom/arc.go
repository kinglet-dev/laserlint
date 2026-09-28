package geom

import "math"

// ArcTo converts an SVG elliptical arc from p0 to p1 into cubic segments.
// rotation is the x-axis rotation in degrees; large and sweep are the SVG
// flags. It follows the endpoint-to-centre conversion in the SVG
// specification's implementation notes, then approximates each piece of at
// most 90° with one cubic.
func ArcTo(p0 Point, rx, ry, rotation float64, large, sweep bool, p1 Point) []Segment {
	// Out-of-range parameters, as the specification requires: no arc between
	// identical points, a straight line for a zero radius, absolute radii.
	if p0 == p1 {
		return nil
	}
	rx, ry = math.Abs(rx), math.Abs(ry)
	if rx == 0 || ry == 0 {
		return []Segment{LineTo(p1)}
	}
	phi := rotation * math.Pi / 180
	cos, sin := math.Cos(phi), math.Sin(phi)

	// Step 1: the start point in the ellipse's own frame, relative to the chord's middle.
	dx, dy := (p0.X-p1.X)/2, (p0.Y-p1.Y)/2
	x1, y1 := cos*dx+sin*dy, -sin*dx+cos*dy

	// Radii too small to span the chord are scaled up until they just fit.
	if lambda := x1*x1/(rx*rx) + y1*y1/(ry*ry); lambda > 1 {
		rx, ry = rx*math.Sqrt(lambda), ry*math.Sqrt(lambda)
	}

	// Step 2: the centre in that frame. Rounding can make the radicand
	// slightly negative for a half ellipse, so it is clamped at zero.
	rx2, ry2, x12, y12 := rx*rx, ry*ry, x1*x1, y1*y1
	coef := math.Sqrt(math.Max(0, (rx2*ry2-rx2*y12-ry2*x12)/(rx2*y12+ry2*x12)))
	if large == sweep {
		coef = -coef
	}
	cxp, cyp := coef*rx*y1/ry, -coef*ry*x1/rx

	// Step 3: the centre in user space.
	cx := cos*cxp - sin*cyp + (p0.X+p1.X)/2
	cy := sin*cxp + cos*cyp + (p0.Y+p1.Y)/2

	// Step 4: start angle and sweep.
	theta := angle(1, 0, (x1-cxp)/rx, (y1-cyp)/ry)
	delta := angle((x1-cxp)/rx, (y1-cyp)/ry, (-x1-cxp)/rx, (-y1-cyp)/ry)
	if !sweep && delta > 0 {
		delta -= 2 * math.Pi
	} else if sweep && delta < 0 {
		delta += 2 * math.Pi
	}

	toUser := func(ux, uy float64) Point {
		return Point{X: cx + rx*ux*cos - ry*uy*sin, Y: cy + rx*ux*sin + ry*uy*cos}
	}
	pieces := int(math.Ceil(math.Abs(delta) / (math.Pi / 2)))
	step := delta / float64(pieces)
	k := 4.0 / 3 * math.Tan(step/4)
	segments := make([]Segment, 0, pieces)
	for i := 0; i < pieces; i++ {
		a1 := theta + float64(i)*step
		a2 := a1 + step
		c1 := toUser(math.Cos(a1)-k*math.Sin(a1), math.Sin(a1)+k*math.Cos(a1))
		c2 := toUser(math.Cos(a2)+k*math.Sin(a2), math.Sin(a2)-k*math.Cos(a2))
		segments = append(segments, CubicTo(c1, c2, toUser(math.Cos(a2), math.Sin(a2))))
	}
	segments[len(segments)-1].To = p1 // land exactly on the end point
	return segments
}

// angle is the signed angle from vector (ux, uy) to (vx, vy).
func angle(ux, uy, vx, vy float64) float64 {
	return math.Atan2(ux*vy-uy*vx, ux*vx+uy*vy)
}
