package geom_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Out-of-range arc parameters, handled as the SVG specification requires.
var _ = Describe("ArcTo with out-of-range parameters", func() {
	It("draws nothing when the end point equals the start point", func() {
		// Arrange
		p := geom.Point{X: 3, Y: 4}

		// Act
		segments := geom.ArcTo(p, 5, 5, 0, false, true, p)

		// Assert
		Expect(segments).To(BeEmpty())
	})

	DescribeTable("draws a straight line when a radius is zero",
		func(rx, ry float64) {
			// Arrange
			start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

			// Act
			segments := geom.ArcTo(start, rx, ry, 0, false, true, end)

			// Assert
			Expect(segments).To(Equal([]geom.Segment{geom.LineTo(end)}))
		},
		Entry("rx is zero", 0.0, 5.0),
		Entry("ry is zero", 5.0, 0.0),
	)

	It("scales radii that are too small up until the arc fits", func() {
		// Arrange: radius 1 can't span 10 units, so it grows to 5, a half circle.
		start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

		// Act
		pts := flattenArc(start, 1, 1, 0, false, true, end)

		// Assert
		for _, p := range pts {
			Expect(math.Hypot(p.X-5, p.Y)).To(BeNumerically("~", 5, 0.002))
		}
	})

	It("uses the absolute value of negative radii", func() {
		// Arrange
		start, end := geom.Point{X: 0, Y: 0}, geom.Point{X: 10, Y: 0}

		// Act
		y := extremeY(flattenArc(start, -5, -5, 0, false, true, end))

		// Assert
		Expect(y).To(BeNumerically("~", -5, 0.002))
	})
})

var _ = Describe("ArcTo with an x-axis rotation", func() {
	It("turns the ellipse, so a 10 × 5 ellipse rotated 90° stands upright", func() {
		// Arrange: half of an upright ellipse centred at the origin, 5 wide and 10 tall.
		start, end := geom.Point{X: 0, Y: -10}, geom.Point{X: 0, Y: 10}

		// Act
		pts := flattenArc(start, 10, 5, 90, false, true, end)

		// Assert
		for _, p := range pts {
			Expect(p.X*p.X/25 + p.Y*p.Y/100).To(BeNumerically("~", 1, 0.001))
		}
	})
})
