package crossing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = Describe("Find measures line drawn on top of other line, within its limits", func() {
	It("measures a line drawn twice as stacked, not as crossings", func() {
		// Act
		// Act: a 10 mm diagonal, so the pair shares cells in several rows and columns.
		r := find(line(pt(0, 0), pt(6, 8)), line(pt(6, 8), pt(0, 0)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
		Expect(r.Stacked).To(BeNumerically("~", 10, 1e-9))
		Expect(r.StackedAt).To(HaveLen(1))
		Expect(r.StackedAt[0].X).To(BeNumerically("~", 3, 1e-9))
		Expect(r.StackedAt[0].Y).To(BeNumerically("~", 4, 1e-9))
	})

	It("measures only the part of two lines that lies on each other", func() {
		// Act
		r := find(line(pt(0, 0), pt(10, 0)), line(pt(6, 0), pt(20, 0)), line(pt(30, 0), pt(40, 0)))

		// Assert
		Expect(r.Stacked).To(BeNumerically("~", 4, 1e-9))
	})

	It("doesn't count a straight line split into pieces that continue each other", func() {
		// Act
		r := find(line(pt(0, 0), pt(5, 0)), line(pt(5, 0), pt(9, 0)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
		Expect(r.Stacked).To(BeZero())
	})

	It("doesn't count pieces that overlap by a rounding error as stacked", func() {
		// Act
		r := find(line(pt(0, 0), pt(5, 0)), line(pt(5-1e-9, 0), pt(9, 0)))

		// Assert
		Expect(r.StackedAt).To(BeEmpty())
	})

	It("doesn't mistake a repeated point on a stacked line for a crossing", func() {
		// Act: the repeated point's line comes first, so its empty piece is
		// compared against the other line rather than the other way round.
		r := find(line(pt(2, 0), pt(4, 0), pt(4, 0), pt(6, 0)), line(pt(0, 0), pt(10, 0)))

		// Assert
		Expect(r.Crossings).To(BeEmpty())
		Expect(r.Stacked).To(BeNumerically("~", 4, 1e-9))
	})

	It("measures a line that folds straight back on itself as stacked", func() {
		// Act
		r := find(line(pt(0, 0), pt(10, 0), pt(6, 0)))

		// Assert
		Expect(r.Stacked).To(BeNumerically("~", 4, 1e-9))
	})

	It("refuses more line pieces than the limit", func() {
		// Arrange
		p := params
		p.MaxSegments = 2

		// Act
		_, err := crossing.Find([]geom.Polyline{line(pt(0, 0), pt(1, 0), pt(2, 0), pt(3, 0))}, p)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})

	It("refuses lines spread over more grid cells than the limit", func() {
		// Arrange: one 100 mm diagonal covers 100 × 100 cells.
		p := params
		p.MaxEntries = 1000

		// Act
		_, err := crossing.Find([]geom.Polyline{line(pt(0, 0), pt(100, 100))}, p)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})
})
