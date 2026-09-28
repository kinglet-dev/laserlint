package spacing_test

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

// params are the defaults: 0.10 mm line + 0.25 mm gap, sampled every 0.01 mm.
var params = spacing.Params{MinDistance: 0.35, Step: 0.01, MaxSamples: 1_000_000}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(points ...geom.Point) geom.Polyline { return geom.Polyline{Points: points} }

func loop(points ...geom.Point) geom.Polyline { return geom.Polyline{Points: points, Closed: true} }

// circle approximates a circle with n straight pieces.
func circle(cx, cy, r float64, n int) geom.Polyline {
	c := geom.Polyline{Closed: true}
	for i := 0; i < n; i++ {
		a := 2 * math.Pi * float64(i) / float64(n)
		c.Points = append(c.Points, pt(cx+r*math.Cos(a), cy+r*math.Sin(a)))
	}
	return c
}

// vee is two 5 mm arms meeting at the origin with the given angle between them.
func vee(degrees float64) geom.Polyline {
	a := degrees * math.Pi / 180
	return line(pt(5, 0), pt(0, 0), pt(5*math.Cos(a), 5*math.Sin(a)))
}

func measure(lines ...geom.Polyline) spacing.Result {
	r, err := spacing.Measure(lines, params)
	if err != nil {
		panic(err)
	}
	return r
}
