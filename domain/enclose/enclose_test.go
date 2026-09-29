package enclose_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/enclose"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func expectCircle(c enclose.Circle, x, y, diameter float64) {
	Expect(c.Centre.X).To(BeNumerically("~", x, 1e-9))
	Expect(c.Centre.Y).To(BeNumerically("~", y, 1e-9))
	Expect(c.Diameter).To(BeNumerically("~", diameter, 1e-9))
}

var _ = Describe("Smallest finds the smallest circle holding every point", func() {
	It("gives an empty circle for no points", func() {
		// Act
		c := enclose.Smallest(nil)

		// Assert
		Expect(c.Diameter).To(BeZero())
	})

	It("gives a circle of no size around a single point", func() {
		// Act
		c := enclose.Smallest([]geom.Point{pt(3, 4)})

		// Assert
		expectCircle(c, 3, 4, 0)
	})

	It("spans two points", func() {
		// Act
		c := enclose.Smallest([]geom.Point{pt(0, 0), pt(6, 8)})

		// Assert
		expectCircle(c, 3, 4, 10)
	})

	It("passes through the far corners of a square", func() {
		// Act
		c := enclose.Smallest([]geom.Point{pt(0, 0), pt(1, 0), pt(1, 1), pt(0, 1)})

		// Assert
		expectCircle(c, 0.5, 0.5, math.Sqrt2)
	})

	It("spans only the longest side of a blunt triangle", func() {
		// Act
		c := enclose.Smallest([]geom.Point{pt(0, 0), pt(5, 1), pt(10, 0)})

		// Assert
		expectCircle(c, 5, 0, 10)
	})

	It("passes through all three corners of a sharp triangle", func() {
		// Act: an equilateral triangle with 3 mm sides.
		h := 3 * math.Sqrt(3) / 2
		c := enclose.Smallest([]geom.Point{pt(0, 0), pt(3, 0), pt(1.5, h)})

		// Assert
		expectCircle(c, 1.5, h/3, 2*math.Sqrt(3))
	})

	It("spans the ends of points in a straight line, repeats included", func() {
		// Act
		c := enclose.Smallest([]geom.Point{pt(1, 0), pt(0, 0), pt(3, 0), pt(1, 0), pt(2, 0)})

		// Assert
		expectCircle(c, 1.5, 0, 3)
	})

	It("fits the circle through many points on a circle, in any order", func() {
		// Arrange: 100 points on a circle of radius 3 around (7, -2).
		var points, reversed []geom.Point
		for i := 0; i < 100; i++ {
			a := 2 * math.Pi * float64(i) / 100
			points = append(points, pt(7+3*math.Cos(a), -2+3*math.Sin(a)))
		}
		for i := len(points) - 1; i >= 0; i-- {
			reversed = append(reversed, points[i])
		}

		// Act
		c, r := enclose.Smallest(points), enclose.Smallest(reversed)

		// Assert
		expectCircle(c, 7, -2, 6)
		expectCircle(r, 7, -2, 6)
	})

	It("copes with points a rounding error apart", func() {
		// Arrange: found by fuzzing; without a tolerance the circle came out NaN.
		points := []geom.Point{pt(94, 0), pt(0.15000000000000002, 2), pt(0.75, 4), pt(0.15, 2)}

		// Act
		c := enclose.Smallest(points)

		// Assert
		Expect(math.IsNaN(c.Diameter)).To(BeFalse())
		Expect(c.Diameter).To(BeNumerically("~", math.Hypot(94-0.15, 2), 0.1))
	})

	It("doesn't reorder the caller's points", func() {
		// Arrange
		points := []geom.Point{pt(0, 0), pt(1, 0), pt(1, 1), pt(0, 1)}

		// Act
		enclose.Smallest(points)

		// Assert
		Expect(points).To(Equal([]geom.Point{pt(0, 0), pt(1, 0), pt(1, 1), pt(0, 1)}))
	})
})
