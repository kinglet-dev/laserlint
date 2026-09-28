package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

var _ = Describe("Read with a broken document", func() {
	It("reports XML that isn't well formed", func() {
		// Arrange
		doc := `<svg xmlns="http://www.w3.org/2000/svg" width="10mm" height="10mm"><path d="M0 0"`

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(reader.ErrBadXML))
	})

	It("reports bad path data and which element it was in", func() {
		// Arrange
		doc := svg(`width="10mm" height="10mm"`, `<path id="wing" d="M0 0 X 1 1"/>`)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(syntax.ErrUnknownCommand))
		Expect(err).To(MatchError(ContainSubstring(`<path id="wing">`)))
	})
})

var _ = Describe("Read with bad path data in an element without an id", func() {
	It("names the element type", func() {
		// Arrange
		doc := svg(`width="10mm" height="10mm"`, `<path d="L 1 1"/>`)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(ContainSubstring("<path>: path data must start with a move")))
	})
})
