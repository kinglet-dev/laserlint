package spacing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

var _ = Describe("Measure treats separate lines that meet end to end as one line", func() {
	It("finds nothing close where two lines meet at a right angle", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 0), pt(10, 10)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("joins lines drawn in opposite directions", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 10), pt(10, 0)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("joins lines that start at the same point", func() {
		// Act
		r := measure(line(pt(10, 0), pt(0, 0)), line(pt(10, 0), pt(10, 10)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("treats a square drawn as four separate sides as a closed loop", func() {
		// Act
		r := measure(
			line(pt(0, 0), pt(10, 0)), line(pt(10, 10), pt(0, 10)),
			line(pt(10, 0), pt(10, 10)), line(pt(0, 0), pt(0, 10)),
		)

		// Assert
		Expect(r.Close).To(BeZero())
		Expect(r.Scored).To(BeNumerically("~", 40, 1e-9))
	})

	It("joins a chain whatever order its pieces come in", func() {
		// Arrange: a U of three pieces, the middle one first.
		bottom, left, right := line(pt(0, 10), pt(10, 10)), line(pt(0, 0), pt(0, 10)), line(pt(10, 10), pt(10, 0))

		// Act
		r := measure(bottom, right, left)

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("treats a line that ends where it starts as a closed loop", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("still counts the corner of two lines that only nearly meet", func() {
		// Arrange: a 0.02 mm gap is a real gap, not a join.
		a, b := line(pt(0, 0), pt(10, 0)), line(pt(10, 0.02), pt(10, 10))

		// Act
		r := measure(a, b)

		// Assert
		Expect(r.Close).To(BeNumerically(">", 0))
	})

	It("counts the stretch near a point where three lines meet", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 0), pt(20, 0)), line(pt(10, 0), pt(10, 10)))

		// Assert: on each of the three arms, the 0.35 mm nearest the junction
		// is within reach of another arm.
		Expect(r.Close).To(BeNumerically("~", 3*0.35, tol))
	})

	It("doesn't join ends one and a half millionths of a millimetre apart", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 0.0000015), pt(10, 10)))

		// Assert
		Expect(r.Close).To(BeNumerically(">", 0))
	})

	It("joins ends that lie either side of a grid cell boundary", func() {
		// Arrange: 10 mm is a boundary of the millionth-of-a-millimetre grid.
		pairs := [][2][2]float64{{{9.9999998, 0}, {10.0000002, 0}}, {{0, 9.9999998}, {0, 10.0000002}}}

		for _, p := range pairs {
			a, b := pt(p[0][0], p[0][1]), pt(p[1][0], p[1][1])

			// Act: a right angle whose corner straddles the boundary.
			r := measure(line(pt(a.X-5*a.Y/10, a.Y-5*a.X/10), a), line(b, pt(b.X+5, b.Y+5)))

			// Assert
			Expect(r.Close).To(BeZero())
		}
	})

	It("joins two lines whose ends differ by less than a millionth of a millimetre", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 0.0000004), pt(10, 10)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("measures three lines meeting at a point the same when their ends are a few millionths of a millimetre apart", func() {
		// Arrange: the first line's start meets the second line's start, which
		// also meets the third line's end, but the first and third don't meet.
		exact := []float64{0, 0, 0}
		spread := []float64{0, 0.0000007, 0.0000014}
		junction := func(x []float64) float64 {
			return measure(
				line(pt(x[0], 0), pt(x[0], 10)),
				line(pt(x[1], 0), pt(10, 0)),
				line(pt(x[2], -10), pt(x[2], 0)),
			).Close
		}

		// Act
		got := junction(spread)

		// Assert
		Expect(got).To(BeNumerically("~", junction(exact), tol))
	})
	It("ignores a lone point where two lines meet", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(10, 0)), line(pt(10, 0), pt(10, 10)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("leaves a closed loop whole when an open line ends at its corner", func() {
		// Act
		r := measure(loop(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10)), line(pt(0, 0), pt(-10, 0)))

		// Assert
		Expect(r.Scored).To(BeNumerically("~", 50, 1e-9))
	})

	It("measures an empty line as no line", func() {
		// Act
		r := measure(line(), line(pt(0, 0), pt(10, 0)))

		// Assert
		Expect(r.Scored).To(BeNumerically("~", 10, 1e-9))
	})
	It("leaves the lines it is given unchanged", func() {
		// Arrange
		lines := []geom.Polyline{line(pt(0, 0), pt(10, 0)), line(pt(10, 10), pt(10, 0))}

		// Act
		_, err := spacing.Measure(lines, params)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(lines[1].Points).To(Equal([]geom.Point{pt(10, 10), pt(10, 0)}))
	})
})
