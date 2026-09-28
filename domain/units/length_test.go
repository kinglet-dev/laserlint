package units_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/units"
)

var _ = Describe("ParseLength", func() {
	Context("when the value is in millimetres", func() {
		It("returns the same number of millimetres", func() {
			// Arrange
			input := "0.25mm"

			// Act
			length, err := units.ParseLength(input)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(length.Millimetres()).To(BeNumerically("~", 0.25, 1e-12))
		})
	})
})

var _ = Describe("ParseLength with other units", func() {
	DescribeTable("converts the value to millimetres",
		func(input string, expectedMM float64) {
			// Arrange: the input and expected value come from the table entry.

			// Act
			length, err := units.ParseLength(input)

			// Assert
			Expect(err).NotTo(HaveOccurred())
			Expect(length.Millimetres()).To(BeNumerically("~", expectedMM, 1e-9))
		},
		Entry("centimetres", "2.5cm", 25.0),
		Entry("inches", "4in", 101.6),
		Entry("points, 72 per inch", "72pt", 25.4),
		Entry("CSS pixels, 96 per inch", "96px", 25.4),
	)
})

var _ = Describe("ParseLength with invalid input", func() {
	DescribeTable("refuses the value with an error that says what is wrong",
		func(input string, expected error) {
			// Arrange: the input and expected error come from the table entry.

			// Act
			_, err := units.ParseLength(input)

			// Assert
			Expect(err).To(MatchError(expected))
		},
		Entry("a bare number, so no unit is guessed", "0.25", units.ErrMissingUnit),
		Entry("an empty value", "", units.ErrMissingUnit),
		Entry("an unknown unit", "3ft", units.ErrUnknownUnit),
		Entry("a unit with no number", "mm", units.ErrInvalidNumber),
		Entry("text instead of a number", "abcmm", units.ErrInvalidNumber),
		Entry("not a number", "NaNmm", units.ErrInvalidNumber),
		Entry("an infinite value", "Infmm", units.ErrInvalidNumber),
		Entry("a negative length", "-0.1mm", units.ErrNegative),
		Entry("a length that overflows when converted", "1e308in", units.ErrTooLarge),
		Entry("a length beyond 10 km", "10000001mm", units.ErrTooLarge),
	)
})
