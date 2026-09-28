package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

var _ = Describe("ParsePathLimited", func() {
	It("accepts a path with exactly as many pieces as allowed", func() {
		// Arrange: one move and two lines are three pieces.
		d := "M0 0 L1 1 L2 2"

		// Act
		_, err := syntax.ParsePathLimited(d, 3)

		// Assert
		Expect(err).NotTo(HaveOccurred())
	})

	It("refuses a path with more pieces than allowed", func() {
		// Arrange
		d := "M0 0 L1 1 L2 2 L3 3"

		// Act
		_, err := syntax.ParsePathLimited(d, 3)

		// Assert
		Expect(err).To(MatchError(syntax.ErrTooComplex))
	})

	It("counts each move as a piece", func() {
		// Arrange
		d := "M0 0 M1 1 M2 2 M3 3"

		// Act
		_, err := syntax.ParsePathLimited(d, 3)

		// Assert
		Expect(err).To(MatchError(syntax.ErrTooComplex))
	})
})

var _ = Describe("ParsePath", func() {
	It("allows up to DefaultMaxPieces pieces", func() {
		// Arrange / Act / Assert: the default is the threat model's 2,000,000.
		Expect(syntax.DefaultMaxPieces).To(Equal(2_000_000))
	})
})
