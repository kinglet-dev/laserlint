// Package geom holds the plane geometry laserlint measures: points, affine
// transforms, paths and polylines. Units are whatever the caller uses; the
// reader converts everything to millimetres.
package geom

// Point is a position in the plane.
type Point struct{ X, Y float64 }

// Matrix is an affine transform in SVG order: x' = A·x + C·y + E, y' = B·x + D·y + F.
type Matrix struct{ A, B, C, D, E, F float64 }

// Apply transforms p.
func (m Matrix) Apply(p Point) Point {
	return Point{X: m.A*p.X + m.C*p.Y + m.E, Y: m.B*p.X + m.D*p.Y + m.F}
}

// Identity returns the transform that changes nothing.
func Identity() Matrix { return Matrix{A: 1, D: 1} }

// Then returns the transform that applies m first and n second.
func (m Matrix) Then(n Matrix) Matrix {
	return Matrix{
		A: n.A*m.A + n.C*m.B,
		B: n.B*m.A + n.D*m.B,
		C: n.A*m.C + n.C*m.D,
		D: n.B*m.C + n.D*m.D,
		E: n.A*m.E + n.C*m.F + n.E,
		F: n.B*m.E + n.D*m.F + n.F,
	}
}
