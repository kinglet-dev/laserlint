package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
)

// square is a visible shape used to show whether content is read.
const square = `<rect width="5" height="5"/>`

var _ = DescribeTable("Read skips elements that are never drawn, with everything inside them",
	func(body string) {
		// Act
		d, err := read(svg(mm96, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(BeEmpty())
	},
	Entry("definitions", `<defs>`+square+`</defs>`),
	Entry("a symbol", `<symbol id="s">`+square+`</symbol>`),
	Entry("a clip path", `<clipPath id="c">`+square+`</clipPath>`),
	Entry("a mask", `<mask id="m">`+square+`</mask>`),
	Entry("a pattern", `<pattern id="p">`+square+`</pattern>`),
	Entry("a marker", `<marker id="k">`+square+`</marker>`),
	Entry("a title and description", `<title>Coaster</title><desc>A bird</desc>`),
	Entry("metadata", `<metadata><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"/></metadata>`),
	Entry("an editor's own elements", `<ink:view xmlns:ink="urn:editor">`+square+`</ink:view>`),
)

var _ = Describe("Read with namespaces", func() {
	It("ignores attributes from other namespaces", func() {
		// Act
		s := only(`<rect xmlns:ink="urn:editor" ink:x="3" width="5" height="5"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(0, 0)))
	})

	It("refuses a file whose root element isn't svg", func() {
		// Act
		_, err := read(`<html><svg ` + mm96 + `/></html>`)

		// Assert
		Expect(err).To(MatchError(reader.ErrNotSVG))
		Expect(err).To(MatchError(ContainSubstring("<html>")))
	})

	It("reports broken XML inside a skipped element", func() {
		// Act
		_, err := read(svg(mm96, `<defs><path></defs>`))

		// Assert
		Expect(err).To(MatchError(reader.ErrBadXML))
	})

	It("reads a document that leaves out the SVG namespace", func() {
		// Act
		d, err := read(`<svg ` + mm96 + `>` + square + `</svg>`)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(HaveLen(1))
	})
})

var _ = DescribeTable("Read refuses content it can't measure and says how to fix it",
	func(body, want string) {
		// Act
		_, err := read(svg(mm96, body))

		// Assert
		Expect(err).To(MatchError(reader.ErrUnsupported))
		Expect(err).To(MatchError(ContainSubstring(want)))
	},
	Entry("a clone", `<use id="u1" href="#a"/>`, `<use id="u1">: clones can't be measured: in Inkscape, Edit → Clone → Unlink Clone`),
	Entry("text", `<text>Hi</text>`, `<text>: text can't be measured: in Inkscape, Path → Object to Path`),
	Entry("an image", `<image href="a.png"/>`, `<image>: images can't be measured: remove the image`),
	Entry("a style sheet", `<style>.a{fill:none}</style>`, `<style>: CSS style sheets aren't supported`),
	Entry("a nested svg", `<svg>`+square+`</svg>`, `<svg>: nested <svg> elements aren't supported`),
	Entry("an element it doesn't know", `<switch>`+square+`</switch>`, `<switch>: this element isn't supported`),
)
