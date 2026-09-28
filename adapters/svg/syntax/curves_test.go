package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func cubic(x1, y1, x2, y2, x, y float64) geom.Segment {
	return geom.CubicTo(pt(x1, y1), pt(x2, y2), pt(x, y))
}

var _ = DescribeTable("ParsePath with curve commands",
	func(d string, want []geom.Segment) {
		// Arrange: d and want come from the table entry.

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths[0].Segments).To(Equal(want))
	},
	Entry("an absolute cubic", "M0 0 C 1 2 3 4 5 6",
		[]geom.Segment{cubic(1, 2, 3, 4, 5, 6)}),
	Entry("a relative cubic", "m1 1 c 1 1 2 2 3 3",
		[]geom.Segment{cubic(2, 2, 3, 3, 4, 4)}),
	Entry("a smooth cubic mirrors the previous second control point", "M0 0 C 0 10 10 10 10 0 S 20 -10 20 0",
		[]geom.Segment{cubic(0, 10, 10, 10, 10, 0), cubic(10, -10, 20, -10, 20, 0)}),
	Entry("a relative smooth cubic", "M0 0 C 0 10 10 10 10 0 s 10 -10 10 0",
		[]geom.Segment{cubic(0, 10, 10, 10, 10, 0), cubic(10, -10, 20, -10, 20, 0)}),
	Entry("a smooth cubic with no cubic before it starts at the pen", "M0 0 S 5 5 10 0",
		[]geom.Segment{cubic(0, 0, 5, 5, 10, 0)}),
	Entry("an absolute quadratic", "M0 0 Q 5 10 10 0",
		[]geom.Segment{geom.QuadTo(pt(0, 0), pt(5, 10), pt(10, 0))}),
	Entry("a relative quadratic", "M1 1 q 5 10 10 0",
		[]geom.Segment{geom.QuadTo(pt(1, 1), pt(6, 11), pt(11, 1))}),
	Entry("a smooth quadratic mirrors the previous control point", "M0 0 Q 5 10 10 0 T 20 0",
		[]geom.Segment{geom.QuadTo(pt(0, 0), pt(5, 10), pt(10, 0)), geom.QuadTo(pt(10, 0), pt(15, -10), pt(20, 0))}),
	Entry("a relative smooth quadratic", "M0 0 Q 5 10 10 0 t 10 0",
		[]geom.Segment{geom.QuadTo(pt(0, 0), pt(5, 10), pt(10, 0)), geom.QuadTo(pt(10, 0), pt(15, -10), pt(20, 0))}),
	Entry("a smooth quadratic with no quadratic before it is straight", "M0 0 T 10 0",
		[]geom.Segment{geom.QuadTo(pt(0, 0), pt(0, 0), pt(10, 0))}),
	Entry("a smooth cubic after a quadratic doesn't mirror it", "M0 0 Q 5 10 10 0 S 15 5 20 0",
		[]geom.Segment{geom.QuadTo(pt(0, 0), pt(5, 10), pt(10, 0)), cubic(10, 0, 15, 5, 20, 0)}),
	Entry("repeated cubic arguments make more cubics", "M0 0 C 1 1 2 2 3 3 4 4 5 5 6 6",
		[]geom.Segment{cubic(1, 1, 2, 2, 3, 3), cubic(4, 4, 5, 5, 6, 6)}),
)

var _ = DescribeTable("ParsePath refuses curves with missing numbers",
	func(d string) {
		// Arrange: d comes from the table entry.

		// Act
		_, err := syntax.ParsePath(d)

		// Assert
		Expect(err).To(MatchError(syntax.ErrBadNumber))
	},
	Entry("cubic with no numbers", "M0 0 C"),
	Entry("cubic", "M0 0 C 1 2 3 4 5"),
	Entry("cubic missing its second control", "M0 0 C 1 2"),
	Entry("smooth cubic", "M0 0 S 1 2 3"),
	Entry("smooth cubic missing its control", "M0 0 S 1"),
	Entry("quadratic", "M0 0 Q 1 2 3"),
	Entry("quadratic missing its control", "M0 0 Q 1"),
	Entry("smooth quadratic", "M0 0 T 1"),
)
