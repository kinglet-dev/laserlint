package spacing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

// tol allows for the length each sample stands for at either end of a run.
const tol = 0.05

var _ = Describe("Measure finds scored line that runs too close to other scored line", func() {
	It("measures the total scored length", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(0, 5), pt(0, 8)))

		// Assert
		Expect(r.Scored).To(BeNumerically("~", 13, 1e-9))
	})

	It("includes the closing side of a closed line", func() {
		// Act
		r := measure(loop(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10)))

		// Assert
		Expect(r.Scored).To(BeNumerically("~", 40, 1e-9))
	})

	It("handles repeated points, such as a return to the start before closing", func() {
		// Act
		r := measure(loop(pt(0, 0), pt(10, 0), pt(10, 0), pt(10, 10), pt(0, 10), pt(0, 0)))

		// Assert
		Expect(r.Scored).To(BeNumerically("~", 40, 1e-9))
		Expect(r.Close).To(BeZero())
	})

	It("finds nothing close on a single straight line", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("counts both of two parallel lines nearer than the minimum distance", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.3), pt(10, 0.3)))

		// Assert
		Expect(r.Close).To(BeNumerically("~", 20, tol))
	})

	It("finds nothing between parallel lines farther apart than the minimum", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.4), pt(10, 0.4)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("counts only the stretch around a crossing", func() {
		// Act
		r := measure(line(pt(-5, 0), pt(5, 0)), line(pt(0, -5), pt(0, 5)))

		// Assert: each line is within 0.35 mm of the other for 0.35 mm either side.
		Expect(r.Close).To(BeNumerically("~", 4*0.35, tol))
	})

	It("counts a line that turns back on itself nearer than the minimum", func() {
		// Act
		r := measure(line(pt(0, 0), pt(10, 0), pt(10, 0.2), pt(0, 0.2)))

		// Assert: both 10 mm legs, less the bend where the line is still turning.
		Expect(r.Close).To(BeNumerically(">", 19))
	})

	It("never counts right angles, so rectangles are clean", func() {
		// Act
		r := measure(loop(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10)))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("never counts a gentle curve", func() {
		// Act
		r := measure(circle(0, 0, 1, 400))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("never counts an obtuse corner", func() {
		// Act
		r := measure(vee(100))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("counts the tip of a sharp corner, and only the tip", func() {
		// Act
		r := measure(vee(30))

		// Assert: worked out from the geometry, each arm is close from 0.109 mm
		// (where the half-circle allowance ends) to 0.35/sin 30° = 0.7 mm from the tip.
		Expect(r.Close).To(BeNumerically("~", 2*(0.7-0.1094), tol))
	})

	It("leaves a loop smaller than the minimum to the small-details check", func() {
		// Act
		r := measure(circle(0, 0, 0.1, 40))

		// Assert
		Expect(r.Close).To(BeZero())
	})

	It("finds other line whose end is within reach even when its midpoint is far, at a coarse step", func() {
		// Arrange: with a 1 mm step the second line is one piece with its
		// midpoint 0.82 mm away, but its end is 0.34 mm from the first line.
		p := params
		p.Step = 1

		// Act
		r, err := spacing.Measure([]geom.Polyline{line(pt(-0.15, 0), pt(0.15, 0)), line(pt(0.34, 0), pt(1.3, 0))}, p)

		// Assert: the first line's single sample, at the origin, is too close.
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Close).To(BeNumerically("~", 0.3, 1e-9))
	})

	It("refuses more samples than the limit", func() {
		// Arrange
		p := params
		p.MaxSamples = 100

		// Act
		_, err := spacing.Measure([]geom.Polyline{line(pt(0, 0), pt(10, 0))}, p)

		// Assert
		Expect(err).To(MatchError(spacing.ErrTooManySamples))
	})
})

var _ = Describe("Measure's spots show where lines are too close", func() {
	It("reports each close stretch with its length and nearest distance", func() {
		// Act: the second line is offset along x, so samples don't line up across.
		r := measure(line(pt(0, 0), pt(10, 0)), line(pt(0.005, 0.3), pt(10.005, 0.3)))

		// Assert
		Expect(r.Spots).To(HaveLen(2))
		Expect(r.Spots[0].Distance).To(BeNumerically("~", 0.3, 1e-9))
		Expect(r.Spots[0].Length).To(BeNumerically("~", 10, tol))
	})

	It("reports separate close stretches of one line separately", func() {
		// Act
		r := measure(line(pt(0, 0), pt(30, 0)), line(pt(0, 0.3), pt(5, 0.3)), line(pt(20, 0.3), pt(25, 0.3)))

		// Assert: two stretches on the long line, one on each short line.
		Expect(r.Spots).To(HaveLen(4))
	})

	It("lists the nearest spots first, located where they are nearest", func() {
		// Act
		r := measure(
			line(pt(0, 0), pt(10, 0)), line(pt(0, 0.3), pt(10, 0.3)),
			line(pt(0, 20), pt(10, 20)), line(pt(0, 20.2), pt(10, 20.1)))

		// Assert
		Expect(r.Spots[0].Distance).To(BeNumerically("~", 0.1, 0.01))
		Expect(r.Spots[0].At.X).To(BeNumerically("~", 10, 0.02))
		Expect(r.Spots[len(r.Spots)-1].Distance).To(BeNumerically("~", 0.3, 1e-9))
	})
})
