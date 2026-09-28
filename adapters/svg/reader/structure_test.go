package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mm96 is a document where one user unit is one millimetre.
const mm96 = `width="100mm" height="100mm" viewBox="0 0 100 100"`

var _ = Describe("Read with groups and transforms", func() {
	It("applies a group's transform to the shapes inside it", func() {
		// Arrange
		doc := svg(mm96, `<g transform="translate(10 20)"><path d="M1 1"/></g>`)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes[0].Path.Subpaths[0].Start).To(Equal(pt(11, 21)))
	})

	It("combines nested transforms, outer group first", func() {
		// Arrange
		doc := svg(mm96, `<g transform="translate(10 0)"><g transform="scale(2)"><path d="M1 1" transform="translate(1 0)"/></g></g>`)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes[0].Path.Subpaths[0].Start).To(Equal(pt(14, 2)))
	})

	It("stops applying a group's transform after the group ends", func() {
		// Arrange
		doc := svg(mm96, `<g transform="translate(50 50)"><path d="M0 0"/></g><path d="M1 1"/>`)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes[1].Path.Subpaths[0].Start).To(Equal(pt(1, 1)))
	})

	It("reports a bad transform and the element it was on", func() {
		// Arrange
		doc := svg(mm96, `<g id="layer1" transform="spin(3)"><path d="M0 0"/></g>`)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(ContainSubstring(`<g id="layer1">: unknown transform function`)))
	})
})

var _ = Describe("Read with paint on the root element", func() {
	It("passes the root svg's fill down to its shapes", func() {
		// Arrange
		doc := svg(mm96+` fill="none" stroke="#000"`, `<path d="M0 0 L1 1"/>`)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes[0].Filled).To(BeFalse())
	})

	It("rejects broken path data even on a shape that won't be scored", func() {
		// Arrange
		doc := svg(mm96, `<path d="M0 0 X" fill="none"/>`)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(HaveOccurred())
	})
})

var _ = DescribeTable("Read works out what the laser will score from fill and stroke",
	func(body string, visible, filled, white bool) {
		// Arrange
		doc := svg(mm96, body)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		if !visible {
			Expect(d.Shapes).To(BeEmpty())
			return
		}
		Expect(d.Shapes).To(HaveLen(1))
		Expect(d.Shapes[0].Filled).To(Equal(filled))
		Expect(d.Shapes[0].WhiteFill).To(Equal(white))
	},
	Entry("black by default", `<path d="M0 0 L1 1"/>`, true, true, false),
	Entry("a hairline: stroke, no fill", `<path d="M0 0 L1 1" fill="none" stroke="#000"/>`, true, false, false),
	Entry("invisible: no fill, no stroke", `<path d="M0 0 L1 1" fill="none"/>`, false, false, false),
	Entry("a white fill", `<path d="M0 0 L1 1" fill="#ffffff"/>`, true, true, true),
	Entry("a short white hex", `<path d="M0 0 L1 1" fill="#FFF"/>`, true, true, true),
	Entry("white by name", `<path d="M0 0 L1 1" fill="white"/>`, true, true, true),
	Entry("the style attribute", `<path d="M0 0 L1 1" style="fill:none;stroke:#000000;stroke-width:0.1"/>`, true, false, false),
	Entry("style wins over the attribute", `<path d="M0 0 L1 1" fill="#000" style="fill: none; stroke: red"/>`, true, false, false),
	Entry("fill inherited from a group", `<g fill="none" stroke="#000"><path d="M0 0 L1 1"/></g>`, true, false, false),
	Entry("a shape overriding its group", `<g fill="none"><path d="M0 0 L1 1" fill="#000"/></g>`, true, true, false),
	Entry("inherit keeps the group's fill", `<g fill="#fff"><path d="M0 0 L1 1" fill="inherit"/></g>`, true, true, true),
	Entry("style names in any case", `<path d="M0 0 L1 1" style="FILL:none;Stroke:#000"/>`, true, false, false),
)
