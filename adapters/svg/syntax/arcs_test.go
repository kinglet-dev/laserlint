package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = DescribeTable("ParsePath with arc commands",
	func(d string, want []geom.Segment) {
		// Arrange: d and want come from the table entry.

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths[0].Segments).To(Equal(want))
	},
	Entry("an absolute arc", "M0 0 A 5 5 0 0 1 10 0",
		geom.ArcTo(pt(0, 0), 5, 5, 0, false, true, pt(10, 0))),
	Entry("a relative arc", "M1 1 a 5 5 0 0 1 10 0",
		geom.ArcTo(pt(1, 1), 5, 5, 0, false, true, pt(11, 1))),
	Entry("flags written without separators", "M0 0 a5 5 0 0110 0",
		geom.ArcTo(pt(0, 0), 5, 5, 0, false, true, pt(10, 0))),
	Entry("large-arc flag set, rotation given", "M0 0 A 10 5 30 1 0 10 0",
		geom.ArcTo(pt(0, 0), 10, 5, 30, true, false, pt(10, 0))),
	Entry("repeated arc arguments make more arcs", "M0 0 A 5 5 0 0 1 10 0 5 5 0 0 1 20 0",
		append(geom.ArcTo(pt(0, 0), 5, 5, 0, false, true, pt(10, 0)), geom.ArcTo(pt(10, 0), 5, 5, 0, false, true, pt(20, 0))...)),
)

var _ = Describe("ParsePath with an arc back to its own start", func() {
	It("draws nothing but keeps the pen at the end point", func() {
		// Arrange
		d := "M 3 3 A 5 5 0 0 1 3 3 L 6 3"

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths[0].Segments).To(Equal([]geom.Segment{geom.LineTo(pt(6, 3))}))
	})
})

var _ = DescribeTable("ParsePath refuses malformed arcs",
	func(d string, want error) {
		// Arrange: d and want come from the table entry.

		// Act
		_, err := syntax.ParsePath(d)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("a flag that isn't 0 or 1", "M0 0 A 5 5 0 2 1 10 0", syntax.ErrBadFlag),
	Entry("a bad sweep flag", "M0 0 A 5 5 0 0 x 10 0", syntax.ErrBadFlag),
	Entry("a missing flag", "M0 0 A 5 5 0", syntax.ErrBadFlag),
	Entry("a missing radius", "M0 0 A 5", syntax.ErrBadNumber),
	Entry("a missing rotation", "M0 0 A 5 5", syntax.ErrBadNumber),
	Entry("a missing end point", "M0 0 A 5 5 0 0 1 10", syntax.ErrBadNumber),
)
