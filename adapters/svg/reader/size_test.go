package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// startOf reads doc and returns the size and where its first path starts, in mm.
func startOf(doc string) (float64, float64, geom.Point) {
	d, err := read(doc)
	Expect(err).NotTo(HaveOccurred())
	return d.Width, d.Height, d.Shapes[0].Path.Subpaths[0].Start
}

var _ = DescribeTable("Read works out the physical size and scale",
	func(attrs, pathData string, w, h float64, want geom.Point) {
		// Arrange
		doc := svg(attrs, `<path d="`+pathData+`"/>`)

		// Act
		gotW, gotH, got := startOf(doc)

		// Assert
		Expect(gotW).To(BeNumerically("~", w, 1e-9))
		Expect(gotH).To(BeNumerically("~", h, 1e-9))
		Expect(got.X).To(BeNumerically("~", want.X, 1e-9))
		Expect(got.Y).To(BeNumerically("~", want.Y, 1e-9))
	},
	Entry("scales the viewBox to the width and height",
		`width="100mm" height="50mm" viewBox="0 0 200 100"`, "M200 100", 100.0, 50.0, pt(100, 50)),
	Entry("shifts by the viewBox origin",
		`width="100mm" height="50mm" viewBox="10 20 100 50"`, "M10 20", 100.0, 50.0, pt(0, 0)),
	Entry("reads a width without units as CSS pixels",
		`width="96" height="48" viewBox="0 0 96 48"`, "M96 48", 25.4, 12.7, pt(25.4, 12.7)),
	Entry("uses the viewBox as CSS pixels when width and height are missing",
		`viewBox="0 0 96 48"`, "M96 0", 25.4, 12.7, pt(25.4, 0)),
	Entry("uses CSS pixels as user units when there is no viewBox",
		`width="50mm" height="50mm"`, "M96 0", 50.0, 50.0, pt(25.4, 0)),
	Entry("reads points, as potrace writes them",
		`width="72pt" height="144pt" viewBox="0 0 72 144"`, "M72 144", 25.4, 50.8, pt(25.4, 50.8)),
	Entry("centres a viewBox with a different aspect ratio (default xMidYMid meet)",
		`width="100mm" height="50mm" viewBox="0 0 100 100"`, "M0 0", 100.0, 50.0, pt(25, 0)),
	Entry("stretches when preserveAspectRatio is none",
		`width="100mm" height="50mm" viewBox="0 0 100 100" preserveAspectRatio="none"`, "M100 100", 100.0, 50.0, pt(100, 50)),
	Entry("accepts commas in the viewBox",
		`width="10mm" height="10mm" viewBox="0,0,10,10"`, "M10 10", 10.0, 10.0, pt(10, 10)),
)

var _ = DescribeTable("Read refuses sizes it can't work out",
	func(attrs string, want error) {
		// Arrange
		doc := svg(attrs, `<path d="M0 0 L1 1"/>`)

		// Act
		_, err := read(doc)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("no size and no viewBox", ``, reader.ErrNoSize),
	Entry("a percentage width", `width="100%" height="100%" viewBox="0 0 10 10"`, reader.ErrNoSize),
	Entry("a zero width", `width="0mm" height="10mm" viewBox="0 0 10 10"`, reader.ErrNoSize),
	Entry("a height in an unknown unit", `width="10mm" height="10ft" viewBox="0 0 10 10"`, reader.ErrNoSize),
	Entry("a viewBox with three numbers", `width="10mm" height="10mm" viewBox="0 0 10"`, reader.ErrBadViewBox),
	Entry("a viewBox with words", `width="10mm" height="10mm" viewBox="0 0 ten 10"`, reader.ErrBadViewBox),
	Entry("a viewBox with zero height", `width="10mm" height="10mm" viewBox="0 0 10 0"`, reader.ErrBadViewBox),
	Entry("an unsupported aspect ratio setting", `width="10mm" height="5mm" viewBox="0 0 10 10" preserveAspectRatio="xMinYMin slice"`, reader.ErrUnsupported),
)
