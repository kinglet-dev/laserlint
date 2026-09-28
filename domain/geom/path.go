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
