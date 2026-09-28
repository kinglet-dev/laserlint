package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = DescribeTable("ParseNumber reads one attribute number",
	func(input string, want float64) {
		// Act
		got, err := syntax.ParseNumber(input)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(want))
	},
	Entry("an integer", "12", 12.0),
	Entry("a signed decimal with spaces around it", " -1.5e1 ", -15.0),
)

var _ = DescribeTable("ParseNumber refuses anything but one plain number",
	func(input string, want error) {
		// Act
		_, err := syntax.ParseNumber(input)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("empty", "", syntax.ErrBadNumber),
	Entry("a unit", "10mm", syntax.ErrUnitInAttribute),
	Entry("a percentage", "50%", syntax.ErrUnitInAttribute),
	Entry("two numbers", "1 2", syntax.ErrBadNumber),
	Entry("too large", "1e8", syntax.ErrOutOfRange),
)

var _ = Describe("ParsePoints", func() {
	It("reads coordinate pairs separated by commas and spaces", func() {
		// Act
		got, err := syntax.ParsePoints(" 0,0 10 0,\n10-5 ")

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal([]geom.Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: -5}}))
	})

	It("refuses an odd number of coordinates", func() {
		// Act
		_, err := syntax.ParsePoints("0 0 10")

		// Assert
		Expect(err).To(MatchError(syntax.ErrOddCoordinates))
	})

	It("refuses a bad number and says where it is", func() {
		// Act
		_, err := syntax.ParsePoints("0 0 x 1")

		// Assert
		Expect(err).To(MatchError(ContainSubstring("expected a number at character 5")))
	})

	It("refuses a bad y coordinate", func() {
		// Act
		_, err := syntax.ParsePoints("0 0 1 +")

		// Assert
		Expect(err).To(MatchError(ContainSubstring("expected a number at character 7")))
	})

	It("stops at a limit on the number of points, so memory stays bounded", func() {
		// Act
		within, errWithin := syntax.ParsePointsLimited("0 0 1 1", 2)
		_, errOver := syntax.ParsePointsLimited("0 0 1 1 2 2", 2)

		// Assert
		Expect(errWithin).NotTo(HaveOccurred())
		Expect(within).To(HaveLen(2))
		Expect(errOver).To(MatchError(syntax.ErrTooComplex))
		Expect(errOver).To(MatchError(ContainSubstring("at character 9")))
	})
})
