package geom_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// cubicAt evaluates a cubic Bézier directly, independently of the code under test.
func cubicAt(p0, p1, p2, p3 geom.Point, t float64) geom.Point {
	u := 1 - t
	a, b, c, d := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
	return geom.Point{X: a*p0.X + b*p1.X + c*p2.X + d*p3.X, Y: a*p0.Y + b*p1.Y + c*p2.Y + d*p3.Y}
}

// distanceToPolyline is the shortest distance from q to any piece of the polyline.
func distanceToPolyline(q geom.Point, pts []geom.Point) float64 {
	best := math.Inf(1)
	for i := 1; i < len(pts); i++ {
		a, b := pts[i-1], pts[i]
		dx, dy := b.X-a.X, b.Y-a.Y
		t := 0.0
		if l := dx*dx + dy*dy; l > 0 {
			t = math.Max(0, math.Min(1, ((q.X-a.X)*dx+(q.Y-a.Y)*dy)/l))
		}
		best = math.Min(best, math.Hypot(q.X-(a.X+t*dx), q.Y-(a.Y+t*dy)))
	}
	return best
}

var _ = Describe("Path.Flatten with a cubic curve", func() {
	var p0, p1, p2, p3 geom.Point
	var lines []geom.Polyline
	const tolerance = 0.01

	BeforeEach(func() {
		// Arrange (shared): an S-shaped curve 100 units long
		p0, p1, p2, p3 = geom.Point{X: 0, Y: 0}, geom.Point{X: 30, Y: 60}, geom.Point{X: 70, Y: -60}, geom.Point{X: 100, Y: 0}
		path := geom.Path{Subpaths: []geom.Subpath{{Start: p0, Segments: []geom.Segment{geom.CubicTo(p1, p2, p3)}}}}

		// Act (shared)
		lines = path.Flatten(tolerance)
	})

	It("stays within the tolerance of every point on the true curve", func() {
		// Arrange
		worst := 0.0

		// Act
		for i := 0; i <= 1000; i++ {
			worst = math.Max(worst, distanceToPolyline(cubicAt(p0, p1, p2, p3, float64(i)/1000), lines[0].Points))
		}

		// Assert
		Expect(worst).To(BeNumerically("<=", tolerance))
	})

	It("ends exactly at the curve's end point", func() {
		// Arrange
		pts := lines[0].Points

		// Act
		last := pts[len(pts)-1]

		// Assert
		Expect(last).To(Equal(p3))
	})
})

var _ = Describe("Path.Flatten with a tolerance that can't be met", func() {
	It("stops splitting a curve after 16 levels, bounding the points to 65,537", func() {
		// Arrange
		curve := geom.CubicTo(geom.Point{X: 0, Y: 100}, geom.Point{X: 100, Y: 100}, geom.Point{X: 100, Y: 0})
		path := geom.Path{Subpaths: []geom.Subpath{{Start: geom.Point{}, Segments: []geom.Segment{curve}}}}

		// Act
		lines := path.Flatten(0)

		// Assert
		Expect(lines[0].Points).To(HaveLen(1<<16 + 1))
	})
})
