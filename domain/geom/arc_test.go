package geom_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// flattenArc flattens a single arc from p0 so tests can inspect its points.
func flattenArc(p0 geom.Point, rx, ry, rotation float64, large, sweep bool, p1 geom.Point) []geom.Point {
	path := geom.Path{Subpaths: []geom.Subpath{{Start: p0, Segments: geom.ArcTo(p0, rx, ry, rotation, large, sweep, p1)}}}
	return path.Flatten(0.001)[0].Points
}

// extremeY returns the point furthest from the x axis.
func extremeY(pts []geom.Point) float64 {
	best := 0.0
	for _, p := range pts {
		if math.Abs(p.Y) > math.Abs(best) {
			best = p.Y
		}
	}
	return best
}

var _ = Describe("ArcTo", func() {
	Context("for a half circle of radius 5 from (0,0) to (10,0)", func() {
		It("keeps every point 5 from the centre (5,0)", func() {
			// Arrange
			start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

			// Act
			pts := flattenArc(start, 5, 5, 0, false, true, end)

			// Assert
			Expect(len(pts)).To(BeNumerically(">", 2))
			for _, p := range pts {
				Expect(math.Hypot(p.X-5, p.Y)).To(BeNumerically("~", 5, 0.002))
			}
			Expect(pts[len(pts)-1].X).To(BeNumerically("~", 10, 1e-9))
		})

		It("bulges to negative y when the sweep flag is set", func() {
			// Arrange
			start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

			// Act
			y := extremeY(flattenArc(start, 5, 5, 0, false, true, end))

			// Assert
			Expect(y).To(BeNumerically("~", -5, 0.002))
		})

		It("bulges to positive y when the sweep flag is clear", func() {
			// Arrange
			start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

			// Act
			y := extremeY(flattenArc(start, 5, 5, 0, false, false, end))

			// Assert
			Expect(y).To(BeNumerically("~", 5, 0.002))
		})
	})
})

var _ = Describe("ArcTo with the large-arc flag", func() {
	// Two circles of radius 10 pass through (0,0) and (10,0); their centres
	// are at (5, ±8.66). The small arc bulges 10 − 8.66 = 1.34 from the chord,
	// the large arc 10 + 8.66 = 18.66.
	DescribeTable("chooses the arc the flags ask for",
		func(large, sweep bool, wantY float64) {
			// Arrange
			start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

			// Act
			y := extremeY(flattenArc(start, 10, 10, 0, large, sweep, end))

			// Assert
			Expect(y).To(BeNumerically("~", wantY, 0.002))
		},
		Entry("small arc, sweep set", false, true, -(10-5*math.Sqrt(3))),
		Entry("small arc, sweep clear", false, false, 10-5*math.Sqrt(3)),
		Entry("large arc, sweep set", true, true, -(10+5*math.Sqrt(3))),
		Entry("large arc, sweep clear", true, false, 10+5*math.Sqrt(3)),
	)
})
