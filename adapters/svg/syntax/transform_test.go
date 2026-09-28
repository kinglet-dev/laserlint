package syntax_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = DescribeTable("ParseTransform moves a point as SVG defines",
	func(transform string, in, want geom.Point) {
		// Arrange: transform, in and want come from the table entry.

		// Act
		m, err := syntax.ParseTransform(transform)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		got := m.Apply(in)
		Expect(got.X).To(BeNumerically("~", want.X, 1e-9))
		Expect(got.Y).To(BeNumerically("~", want.Y, 1e-9))
	},
	Entry("no transform", "", pt(3, 4), pt(3, 4)),
	Entry("translate by x and y", "translate(10 20)", pt(1, 1), pt(11, 21)),
	Entry("translate by x only", "translate(5)", pt(1, 1), pt(6, 1)),
	Entry("scale x and y", "scale(2 3)", pt(1, 1), pt(2, 3)),
	Entry("scale uniformly", "scale(2)", pt(1, 1), pt(2, 2)),
	Entry("rotate by 90°", "rotate(90)", pt(1, 0), pt(0, 1)),
	Entry("rotate about a centre", "rotate(90 10 10)", pt(11, 10), pt(10, 11)),
	Entry("skew along x", "skewX(45)", pt(0, 1), pt(1, 1)),
	Entry("skew along y", "skewY(45)", pt(1, 0), pt(1, 1)),
	Entry("a full matrix", "matrix(1 2 3 4 5 6)", pt(1, 1), pt(9, 12)),
	Entry("a list applies the last transform first", "translate(10 0) scale(2)", pt(1, 1), pt(12, 2)),
	Entry("commas and spaces anywhere", " translate( 10 , 0 ) , scale(2) ", pt(1, 1), pt(12, 2)),
)

var _ = DescribeTable("ParseTransform refuses malformed transforms",
	func(transform string, want error) {
		// Arrange: transform and want come from the table entry.

		// Act
		_, err := syntax.ParseTransform(transform)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("an unknown function", "spin(3)", syntax.ErrUnknownTransform),
	Entry("no arguments", "translate()", syntax.ErrTransformArgs),
	Entry("too few matrix arguments", "matrix(1 2 3)", syntax.ErrTransformArgs),
	Entry("two arguments to rotate", "rotate(1 2)", syntax.ErrTransformArgs),
	Entry("too many arguments", "scale(1 2 3)", syntax.ErrTransformArgs),
	Entry("one argument too many for skew", "skewX(1 2)", syntax.ErrTransformArgs),
	Entry("a missing opening parenthesis", "translate 10", syntax.ErrTransformSyntax),
	Entry("a missing closing parenthesis", "translate(10", syntax.ErrTransformSyntax),
	Entry("a bad number", "translate(ten)", syntax.ErrTransformSyntax),
	Entry("a sign with no digits", "translate(-)", syntax.ErrTransformSyntax),
	Entry("a number out of range", "translate(1e999)", syntax.ErrOutOfRange),
)

var _ = Describe("ParseTransform with a list that overflows", func() {
	It("refuses it instead of returning an infinite matrix", func() {
		// Arrange
		list := strings.Repeat("scale(10000000) ", 50)

		// Act
		_, err := syntax.ParseTransform(list)

		// Assert
		Expect(err).To(MatchError(syntax.ErrOutOfRange))
	})
})
