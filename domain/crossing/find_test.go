package crossing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var params = crossing.Params{Merge: 0.1, Cell: 1, MaxSegments: 100_000, MaxEntries: 1_000_000}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(points ...geom.Point) geom.Polyline { return geom.Polyline{Points: points} }

func loop(points ...geom.Point) geom.Polyline { return geom.Polyline{Points: points, Closed: true} }

func find(lines ...geom.Polyline) crossing.Result {
	r, err := crossing.Find(lines, params)
	Expect(err).NotTo(HaveOccurred())
	return r
}

var _ = Describe("Find locates where score lines cross, touch or lie on each other", func() {
	It("finds nothing in an empty design", func() {
		// Act
		r := find()

		// Assert
		Expect(r.Crossings).To(BeEmpty())
		Expect(r.Stacked).To(BeZero())
	})

	It("skips closed lines with no points or a single point", func() {
		// Act
		r := find(loop(), loop(pt(1, 1)), line(pt(0, 0), pt(2, 2)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
	})

	It("finds where two lines cross", func() {
		// Act
		r := find(line(pt(0, 0), pt(10, 10)), line(pt(0, 10), pt(10, 0)))

		// Assert
		Expect(r.Crossings).To(HaveLen(1))
		Expect(r.Crossings[0].X).To(BeNumerically("~", 5, 1e-9))
		Expect(r.Crossings[0].Y).To(BeNumerically("~", 5, 1e-9))
	})

	It("finds nothing between lines that don't meet", func() {
		// Act
		r := find(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.3), pt(10, 0.3)), line(pt(0, 5), pt(4, 1)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
	})

	It("finds nothing where a line would reach another only if it were longer", func() {
		// Act: a short steep line passes the end of a flat one; the pair is
		// given in both orders.
		r := find(line(pt(0, 0), pt(10, 0)), line(pt(10.5, 1), pt(10.6, -1)),
			line(pt(30.5, 1), pt(30.6, -1)), line(pt(20, 0), pt(30, 0)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
	})

	It("counts crossings nearer than the merge distance as one, across grid cells", func() {
		// Act: two crossings 0.04 mm apart side by side, two more one above the other.
		r := find(line(pt(0, 5), pt(1, 5)), line(pt(0.29, 4), pt(0.29, 6)), line(pt(0.33, 4), pt(0.33, 6)),
			line(pt(5, 0), pt(5, 1)), line(pt(4, 0.29), pt(6, 0.29)), line(pt(4, 0.33), pt(6, 0.33)))

		// Assert
		Expect(r.Crossings).To(HaveLen(2))
	})

	It("counts a line ending on another as a crossing", func() {
		// Act: a T.
		r := find(line(pt(0, 0), pt(10, 0)), line(pt(5, 5), pt(5, 0)))

		// Assert
		Expect(r.Crossings).To(HaveLen(1))
	})

	It("doesn't count separate lines that join end to end, as split paths do", func() {
		// Act
		r := find(line(pt(0, 0), pt(5, 0)), line(pt(5, 0), pt(5, 5)), line(pt(5, 5), pt(9, 5)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
		Expect(r.Stacked).To(BeZero())
	})

	It("doesn't count the corners of a line, closing corner included", func() {
		// Act
		r := find(loop(pt(0, 0), pt(10, 0), pt(10, 10), pt(0, 10)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
	})

	It("finds where another line crosses a closed line's closing side", func() {
		// Act: the triangle's closing side runs from (10, 10) back to (0, 0).
		r := find(loop(pt(0, 0), pt(10, 0), pt(10, 10)), line(pt(2, 8), pt(8, 2)))

		// Assert
		Expect(r.Crossings).To(HaveLen(1))
		Expect(r.Crossings[0].X).To(BeNumerically("~", 5, 1e-9))
	})

	It("finds where a line crosses itself", func() {
		// Act: a figure of eight.
		r := find(loop(pt(0, 0), pt(10, 10), pt(10, 0), pt(0, 10)))

		// Assert
		Expect(r.Crossings).To(HaveLen(1))
	})

	It("counts a crossing at a line's corner once", func() {
		// Act: a line crossing a V exactly at its tip.
		r := find(line(pt(0, 0), pt(5, 5), pt(10, 0)), line(pt(5, 0), pt(5, 10)))

		// Assert
		Expect(r.Crossings).To(HaveLen(1))
	})

	It("counts each crossing once, however many cells the lines share", func() {
		// Arrange: a grid of 10 horizontal and 10 vertical lines, each 50 mm long.
		var lines []geom.Polyline
		for i := 0; i < 10; i++ {
			f := 2.5 + 5*float64(i)
			lines = append(lines, line(pt(0, f), pt(50, f)), line(pt(f, 0), pt(f, 50)))
		}

		// Act
		r := find(lines...)

		// Assert
		Expect(r.Crossings).To(HaveLen(100))
	})
})
