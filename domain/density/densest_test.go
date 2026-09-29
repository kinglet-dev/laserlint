package density_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/density"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// params: a 10 mm window moved in 0.5 mm steps.
var params = density.Params{Window: 10, Step: 0.5, MaxCells: 1_000_000}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(points ...geom.Point) geom.Polyline { return geom.Polyline{Points: points} }

func densest(lines ...geom.Polyline) density.Result {
	r, err := density.Densest(lines, params)
	Expect(err).NotTo(HaveOccurred())
	return r
}

// hatch is twenty 10 mm lines, 0.5 mm apart, filling the square at (x, y).
func hatch(x, y float64) []geom.Polyline {
	var out []geom.Polyline
	for i := 0; i < 20; i++ {
		yy := y + 0.25 + 0.5*float64(i)
		out = append(out, line(pt(x, yy), pt(x+10, yy)))
	}
	return out
}

var _ = Describe("Densest finds the most scored line in any window", func() {
	It("finds nothing in an empty design", func() {
		// Act
		r := densest()

		// Assert
		Expect(r.Density).To(BeZero())
	})

	It("measures a short line as its length over the window's area", func() {
		// Act
		r := densest(line(pt(1, 1), pt(9, 1)))

		// Assert: 8 mm of line in 100 mm².
		Expect(r.Density).To(BeNumerically("~", 0.08, 1e-9))
	})

	It("counts only the part of a long line inside the window", func() {
		// Act
		r := densest(line(pt(0, 1), pt(30, 1)))

		// Assert: 10 mm of line in 100 mm².
		Expect(r.Density).To(BeNumerically("~", 0.10, 1e-9))
	})

	It("measures a slanted line by its true length", func() {
		// Act
		r := densest(line(pt(1, 1), pt(4, 5)))

		// Assert: a 3-4-5 line, 5 mm long.
		Expect(r.Density).To(BeNumerically("~", 0.05, 1e-9))
	})

	DescribeTable("measures a line the same whichever way it runs",
		func(a, b geom.Point) {
			// Act
			r := densest(line(a, b))

			// Assert: 5 mm of line in 100 mm².
			Expect(r.Density).To(BeNumerically("~", 0.05, 1e-9))
		},
		Entry("right to left", pt(6.3, 2.2), pt(1.3, 2.2)),
		Entry("downwards", pt(2.2, 1.3), pt(2.2, 6.3)),
		Entry("upwards", pt(2.2, 6.3), pt(2.2, 1.3)),
		Entry("up and to the left", pt(4.7, 5.9), pt(1.7, 1.9)),
	)

	DescribeTable("splits a line exactly where it crosses cell edges",
		func(a, b geom.Point) {
			// Act
			r := densest(line(a, b))

			// Assert: windows sit on the 0.5 mm grid, so any window holds at
			// most 9.75 mm of this 10 mm line.
			Expect(r.Density).To(BeNumerically("~", 0.0975, 1e-9))
		},
		Entry("left to right", pt(0.25, 1.1), pt(10.25, 1.1)),
		Entry("right to left", pt(10.25, 1.1), pt(0.25, 1.1)),
		Entry("downwards", pt(1.1, 0.25), pt(1.1, 10.25)),
		Entry("upwards", pt(1.1, 10.25), pt(1.1, 0.25)),
	)

	It("includes the closing side of a closed line, in a design smaller than the window", func() {
		// Act
		r := densest(geom.Polyline{Points: []geom.Point{pt(1, 1), pt(5, 1), pt(5, 5), pt(1, 5)}, Closed: true})

		// Assert: a 4 mm square, 16 mm of line.
		Expect(r.Density).To(BeNumerically("~", 0.16, 1e-9))
	})

	It("finds the densest area and says where its centre is", func() {
		// Arrange: a sparse line far from a dense hatch.
		lines := append(hatch(40, 40), line(pt(0, 0), pt(10, 0)))

		// Act
		r := densest(lines...)

		// Assert: 200 mm of line in 100 mm².
		Expect(r.Density).To(BeNumerically("~", 2.0, 1e-9))
		Expect(r.At.X).To(BeNumerically("~", 45, 0.5))
		Expect(r.At.Y).To(BeNumerically("~", 45, 0.5))
	})

	// Found by fuzzing: rounding at a cell edge stepped the traversal outside
	// the grid, an index out of range. More such inputs are kept in
	// testdata/fuzz/FuzzDensest and run with every test run.
	DescribeTable("stays inside the grid when rounding lands on a cell edge",
		func(a, b geom.Point) {
			// Act
			r := densest(line(a, b))

			// Assert: at most the whole line's length over the window's area.
			Expect(r.Density).To(BeNumerically(">", 0))
			Expect(r.Density).To(BeNumerically("<=", math.Hypot(b.X-a.X, b.Y-a.Y)/100+1e-9))
		},
		Entry("stepping vertically", pt(-0.075, -7.5), pt(18, -60)),
	)

	It("refuses a design spread over more cells than the limit", func() {
		// Arrange: 2 km of line needs 4000 × 1 cells at 0.5 mm.
		p := params
		p.MaxCells = 1000

		// Act
		_, err := density.Densest([]geom.Polyline{line(pt(0, 0), pt(2000, 0))}, p)

		// Assert
		Expect(err).To(MatchError(density.ErrTooLarge))
	})
})
