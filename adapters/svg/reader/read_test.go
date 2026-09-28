package reader_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// svg wraps body in an <svg> element with the given attributes.
func svg(attrs, body string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" ` + attrs + `>` + body + `</svg>`
}

func read(doc string) (design.Design, error) {
	return reader.Read(strings.NewReader(doc))
}

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

var _ = Describe("Read", func() {
	It("reads a path in a document sized in millimetres", func() {
		// Arrange
		doc := svg(`width="100mm" height="50mm" viewBox="0 0 100 50"`, `<path d="M0 0 L10 0"/>`)

		// Act
		d, err := read(doc)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Width).To(BeNumerically("~", 100, 1e-9))
		Expect(d.Height).To(BeNumerically("~", 50, 1e-9))
		Expect(d.Shapes).To(HaveLen(1))
		Expect(d.Shapes[0].Filled).To(BeTrue(), "SVG fills shapes black by default")
		Expect(d.Shapes[0].Path.Subpaths[0].Start).To(Equal(pt(0, 0)))
		Expect(d.Shapes[0].Path.Subpaths[0].Segments[0].To).To(Equal(pt(10, 0)))
	})
})
