package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

var _ = DescribeTable("ParsePath reads numbers in every form the SVG grammar allows",
	func(d string, x, y float64) {
		// Arrange: d, x and y come from the table entry.

		// Act
		path, err := syntax.ParsePath(d)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(path.Subpaths[0].Start).To(Equal(pt(x, y)))
	},
	Entry("an exponent", "M1e2 0", 100.0, 0.0),
	Entry("an explicit plus sign", "M+5 -5", 5.0, -5.0),
	Entry("two fractions run together", "M.5.5", 0.5, 0.5),
	Entry("signed exponents in either case", "M1.5e-1 2E+1", 0.15, 20.0),
	Entry("a minus sign that starts the next number", "M-1-2", -1.0, -2.0),
	Entry("a trailing decimal point", "M3. 4.", 3.0, 4.0),
)

var _ = DescribeTable("ParsePath refuses these number and range errors",
	func(d string, want error) {
		// Arrange: d and want come from the table entry.

		// Act
		_, err := syntax.ParsePath(d)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("an 'e' with no exponent digits is not part of the number", "M 0 0 L 2 3e 4", syntax.ErrUnknownCommand),
	Entry("a relative horizontal line past the range", "M 9000000 0 h 9000000", syntax.ErrOutOfRange),
	Entry("a relative vertical line past the range", "M 0 9000000 v 9000000", syntax.ErrOutOfRange),
)
