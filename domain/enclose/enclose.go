// Package enclose finds the smallest circle that holds a set of points,
// with Welzl's randomised incremental algorithm (Welzl 1991; de Berg et al.,
// "Computational Geometry", section 4.7). Points are shuffled with a fixed
// seed, so the expected time is linear and every run gives the same answer.
package enclose

import (
	"math"
	"math/rand/v2"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Circle is a circle in mm.
type Circle struct {
	Centre   geom.Point
	Diameter float64
}

// disc is a circle by its radius, the form the algorithm works in.
type disc struct {
	c geom.Point
	r float64
}

// holds reports whether p is in the disc, allowing for rounding.
func (d disc) holds(p geom.Point) bool {
	tol := 1e-9 * (1 + d.r + math.Abs(d.c.X) + math.Abs(d.c.Y))
	return math.Hypot(p.X-d.c.X, p.Y-d.c.Y) <= d.r+tol
}

// Smallest returns the smallest circle holding every point.
func Smallest(points []geom.Point) Circle {
	if len(points) == 0 {
		return Circle{}
	}
	p := append([]geom.Point(nil), points...)
	rand.New(rand.NewPCG(1, 2)).Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
	d := disc{c: p[0]}
	for i := 1; i < len(p); i++ {
		if !d.holds(p[i]) {
			d = withOne(p[:i], p[i])
		}
	}
	return Circle{Centre: d.c, Diameter: 2 * d.r}
}

// withOne is the smallest disc holding points with q on its edge.
func withOne(points []geom.Point, q geom.Point) disc {
	d := disc{c: q}
	for j, a := range points {
		if !d.holds(a) {
			d = withTwo(points[:j], q, a)
		}
	}
	return d
}

// withTwo is the smallest disc holding points with q and a on its edge.
func withTwo(points []geom.Point, q, a geom.Point) disc {
	d := spanning(q, a)
	for _, b := range points {
		if !d.holds(b) {
			d = through(q, a, b)
		}
	}
	return d
}

// spanning is the disc with a and b at opposite ends.
func spanning(a, b geom.Point) disc {
	return disc{c: geom.Point{X: (a.X + b.X) / 2, Y: (a.Y + b.Y) / 2}, r: math.Hypot(b.X-a.X, b.Y-a.Y) / 2}
}

// through is the disc through three points. The algorithm only calls it
// when b lies outside the disc spanning q and a, both of which must be on
// the edge: so the three are never in a straight line.
func through(a, b, c geom.Point) disc {
	bx, by, cx, cy := b.X-a.X, b.Y-a.Y, c.X-a.X, c.Y-a.Y
	det := 2 * (bx*cy - by*cx)
	b2, c2 := bx*bx+by*by, cx*cx+cy*cy
	ux, uy := (cy*b2-by*c2)/det, (bx*c2-cx*b2)/det
	return disc{c: geom.Point{X: a.X + ux, Y: a.Y + uy}, r: math.Hypot(ux, uy)}
}
