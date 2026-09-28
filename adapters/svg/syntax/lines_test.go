package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

var _ = Describe("ParsePath with move and line commands", func() {
	It("reads an absolute move followed by an absolute line", func() {
		// Arrange
		d := "M 10 20 L 30 40"

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths).To(Equal([]geom.Subpath{{Start: pt(10, 20), Segments: []geom.Segment{geom.LineTo(pt(30, 40))}}}))
	})
})

// line is shorthand for a straight segment to (x, y).
func line(x, y float64) geom.Segment { return geom.LineTo(pt(x, y)) }

var _ = DescribeTable("ParsePath with the other straight-line forms",
	func(d string, want []geom.Subpath) {
		// Arrange: d and want come from the table entry.

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths).To(Equal(want))
	},
	Entry("relative move and line", "m 10 20 l 5 5",
		[]geom.Subpath{{Start: pt(10, 20), Segments: []geom.Segment{line(15, 25)}}}),
	Entry("horizontal and vertical lines", "M 1 1 H 5 V 7 h -2 v -3",
		[]geom.Subpath{{Start: pt(1, 1), Segments: []geom.Segment{line(5, 1), line(5, 7), line(3, 7), line(3, 4)}}}),
	Entry("extra pairs after a move are lines", "M 0 0 10 0 10 10",
		[]geom.Subpath{{Start: pt(0, 0), Segments: []geom.Segment{line(10, 0), line(10, 10)}}}),
	Entry("extra pairs after a relative move are relative lines", "m 1 1 2 0 0 2",
		[]geom.Subpath{{Start: pt(1, 1), Segments: []geom.Segment{line(3, 1), line(3, 3)}}}),
	Entry("a repeated line command", "M 0 0 L 1 1 2 2",
		[]geom.Subpath{{Start: pt(0, 0), Segments: []geom.Segment{line(1, 1), line(2, 2)}}}),
	Entry("close path, then a relative move from the closed subpath's start", "M 5 5 L 9 5 Z m 1 1 l 1 0",
		[]geom.Subpath{
			{Start: pt(5, 5), Segments: []geom.Segment{line(9, 5)}, Closed: true},
			{Start: pt(6, 6), Segments: []geom.Segment{line(7, 6)}},
		}),
	Entry("a line straight after close starts a new subpath at the old start", "M 0 0 L 4 0 Z L 0 4",
		[]geom.Subpath{
			{Start: pt(0, 0), Segments: []geom.Segment{line(4, 0)}, Closed: true},
			{Start: pt(0, 0), Segments: []geom.Segment{line(0, 4)}},
		}),
	Entry("no path data at all", "  ", []geom.Subpath(nil)),
)
