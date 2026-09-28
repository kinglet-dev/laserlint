package reader_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

// roomy is large enough for every document below except where one limit is lowered.
var roomy = reader.Limits{Bytes: 10_000, Depth: 10, Elements: 100, Pieces: 100}

func readWith(doc string, change func(*reader.Limits)) error {
	limits := roomy
	change(&limits)
	_, err := reader.ReadLimited(strings.NewReader(doc), limits)
	return err
}

var _ = Describe("Read's limits on hostile or runaway files", func() {
	It("uses the limits in the threat model by default", func() {
		// Assert
		Expect(reader.DefaultLimits).To(Equal(reader.Limits{
			Bytes: 25_000_000, Depth: 256, Elements: 200_000, Pieces: 2_000_000}))
	})

	It("reads a file of exactly the byte limit", func() {
		// Arrange
		doc := svg(mm96, square)

		// Act
		err := readWith(doc, func(l *reader.Limits) { l.Bytes = int64(len(doc)) })

		// Assert
		Expect(err).NotTo(HaveOccurred())
	})

	It("refuses a file over the byte limit", func() {
		// Arrange
		doc := svg(mm96, square)

		// Act
		err := readWith(doc, func(l *reader.Limits) { l.Bytes = int64(len(doc)) - 1 })

		// Assert
		Expect(err).To(MatchError(reader.ErrTooLarge))
	})

	It("refuses nesting deeper than the limit", func() {
		// Act
		err := readWith(svg(mm96, `<g><g></g></g>`), func(l *reader.Limits) { l.Depth = 2 })

		// Assert
		Expect(err).To(MatchError(reader.ErrTooDeep))
		Expect(err).To(MatchError(ContainSubstring("more than 2 levels")))
	})

	It("counts nesting inside skipped elements too", func() {
		// Act
		err := readWith(svg(mm96, `<defs><g><g/></g></defs>`), func(l *reader.Limits) { l.Depth = 3 })

		// Assert
		Expect(err).To(MatchError(reader.ErrTooDeep))
	})

	It("refuses more elements than the limit, skipped ones included", func() {
		// Act
		err := readWith(svg(mm96, `<defs><g/><g/></defs>`), func(l *reader.Limits) { l.Elements = 3 })

		// Assert
		Expect(err).To(MatchError(reader.ErrTooManyElements))
		Expect(err).To(MatchError(ContainSubstring("more than 3 elements")))
	})

	It("refuses more path pieces in total than the limit", func() {
		// Arrange
		doc := svg(mm96, `<path d="M0 0 L1 0 L1 1"/><path d="M0 0 L1 0 L1 1"/>`)

		// Act
		err := readWith(doc, func(l *reader.Limits) { l.Pieces = 3 })

		// Assert
		Expect(err).To(MatchError(syntax.ErrTooComplex))
		Expect(err).To(MatchError(ContainSubstring("<path>")))
	})

	DescribeTable("counts the pieces of every kind of shape",
		func(body string) {
			// Act
			err := readWith(svg(mm96, body), func(l *reader.Limits) { l.Pieces = 3 })

			// Assert
			Expect(err).To(MatchError(syntax.ErrTooComplex))
		},
		Entry("a circle's four arcs", `<circle r="1"/>`),
		Entry("a polygon's points", `<polygon points="0 0 1 0 1 1 0 1 0 2"/>`),
		Entry("an invisible shape's pieces", `<path d="M0 0 L1 0 L1 1 L2 2" fill="none"/>`),
	)

	It("stops reading a long point list at the limit instead of after it", func() {
		// Act
		err := readWith(svg(mm96, `<polygon points="0 0 1 0 1 1 0 1 0 2 5 5"/>`), func(l *reader.Limits) { l.Pieces = 3 })

		// Assert
		Expect(err).To(MatchError(ContainSubstring("at character 17")))
	})
})

var _ = Describe("Read's defences against XML attacks and extreme numbers", func() {
	It("accepts the plain DOCTYPE line older tools write", func() {
		// Arrange
		doc := `<?xml version="1.0"?><!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 20010904//EN" ` +
			`"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">` + svg(mm96, square)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
	})

	It("refuses a DOCTYPE that defines entities", func() {
		// Arrange
		doc := `<!DOCTYPE svg [<!ENTITY a "aaaa">]>` + svg(mm96, square)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(reader.ErrEntities))
	})

	DescribeTable("refuses coordinates that a transform pushes out of range",
		func(d string) {
			// Act
			_, err := read(svg(mm96, `<path id="p" d="`+d+`" transform="scale(10)"/>`))

			// Assert
			Expect(err).To(MatchError(syntax.ErrOutOfRange))
			Expect(err).To(MatchError(ContainSubstring(`<path id="p">`)))
		},
		Entry("a start point", "M9e6 0"),
		Entry("an end point", "M0 0 L0 -9e6"),
		Entry("a control point", "M0 0 C0 0 9e6 0 1 1"),
	)
})
