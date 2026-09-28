package geom

// Path is a shape outline made of subpaths, as in SVG path data.
type Path struct{ Subpaths []Subpath }

// Subpath starts at Start and follows its segments; Closed joins the end back to Start.
type Subpath struct {
	Start    Point
	Segments []Segment
	Closed   bool
}

// Segment is a straight line or a cubic Bézier curve ending at To.
type Segment struct {
	C1, C2, To Point
	Curved     bool
}

// Polyline is a flattened subpath: straight pieces between Points.
type Polyline struct {
	Points []Point
	Closed bool
}

// LineTo returns a straight segment to p.
func LineTo(p Point) Segment { return Segment{To: p} }

// Flatten turns the path into polylines that stay within tolerance of every curve.
func (p Path) Flatten(tolerance float64) []Polyline {
	lines := make([]Polyline, 0, len(p.Subpaths))
	for _, sub := range p.Subpaths {
		points := []Point{sub.Start}
		for _, seg := range sub.Segments {
			if seg.Curved {
				points = appendCubic(points, points[len(points)-1], seg.C1, seg.C2, seg.To, tolerance, 0)
				continue
			}
			points = append(points, seg.To)
		}
		lines = append(lines, Polyline{Points: points, Closed: sub.Closed})
	}
	return lines
}

// CubicTo returns a cubic Bézier segment with control points c1, c2, ending at to.
func CubicTo(c1, c2, to Point) Segment { return Segment{C1: c1, C2: c2, To: to, Curved: true} }

// QuadTo returns the cubic equivalent of the quadratic Bézier from p0 via control q to p2.
func QuadTo(p0, q, p2 Point) Segment {
	// Degree elevation: each cubic control point lies 2/3 of the way to q.
	return CubicTo(
		Point{X: p0.X + 2.0/3*(q.X-p0.X), Y: p0.Y + 2.0/3*(q.Y-p0.Y)},
		Point{X: p2.X + 2.0/3*(q.X-p2.X), Y: p2.Y + 2.0/3*(q.Y-p2.Y)},
		p2,
	)
}

// Transform returns the path with m applied to every point. Affine transforms
// map Bézier curves exactly by mapping their control points.
func (p Path) Transform(m Matrix) Path {
	out := Path{Subpaths: make([]Subpath, len(p.Subpaths))}
	for i, sub := range p.Subpaths {
		segs := make([]Segment, len(sub.Segments))
		for j, s := range sub.Segments {
			segs[j] = Segment{C1: m.Apply(s.C1), C2: m.Apply(s.C2), To: m.Apply(s.To), Curved: s.Curved}
		}
		out.Subpaths[i] = Subpath{Start: m.Apply(sub.Start), Segments: segs, Closed: sub.Closed}
	}
	return out
}
