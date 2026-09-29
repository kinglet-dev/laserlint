package background_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/background"
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// rect is a filled rectangle, white or black.
func rect(x, y, w, h float64, white bool) design.Shape {
	return design.Shape{Filled: true, WhiteFill: white, Path: geom.Path{Subpaths: []geom.Subpath{{
		Start:    geom.Point{X: x, Y: y},
		Segments: []geom.Segment{geom.LineTo(geom.Point{X: x + w, Y: y}), geom.LineTo(geom.Point{X: x + w, Y: y + h}), geom.LineTo(geom.Point{X: x, Y: y + h})},
		Closed:   true,
	}}}}
}

// run checks a 100 × 50 mm design.
func run(shapes ...design.Shape) []check.Finding {
	in := check.Input{Design: design.Design{Width: 100, Height: 50, Shapes: shapes}}
	f, err := background.New().Run(in, check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The background check", func() {
	It("is a check with a stable id and a name for people", func() {
		// Act
		var c check.Check = background.New()

		// Assert
		Expect(c.ID()).To(Equal("background"))
		Expect(c.Name()).To(Equal("Background shape"))
	})

	It("finds nothing without white shapes", func() {
		// Act
		f := run(rect(0, 0, 100, 50, false))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("ignores a white shape with no outline", func() {
		// Act: as from <path d="" fill="white"/>.
		f := run(design.Shape{Filled: true, WhiteFill: true})

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("warns about a white shape covering the design, and says to delete it", func() {
		// Act
		f := run(rect(0, 0, 100, 50, true), rect(10, 10, 5, 5, false))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("background"))
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(Equal("a white background shape (100.0 × 50.0 mm) is scored around its edge like any other shape"))
		Expect(f[0].Fix).To(ContainSubstring("Delete"))
		Expect(f[0].Fix).To(ContainSubstring("hidden"))
		Expect(f[0].Locations).To(Equal([]geom.Point{{X: 50, Y: 25}}))
	})

	It("takes a white shape covering 90% of the design as a background", func() {
		// Act: 90 × 50 mm of 100 × 50 mm.
		f := run(rect(5, 0, 90, 50, true))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
	})

	It("notes smaller white shapes as information, largest first", func() {
		// Act: 89% of the design, then a small square.
		f := run(rect(80, 40, 2, 2, true), rect(0, 0, 89, 50, true))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Severity).To(Equal(check.Info))
		Expect(f[0].Message).To(Equal("2 white-filled shapes are scored like any other shape: white doesn't mean not burned"))
		Expect(f[0].Fix).NotTo(BeEmpty())
		Expect(f[0].Locations[0]).To(Equal(geom.Point{X: 44.5, Y: 25}))
	})

	It("words one small white shape and several backgrounds", func() {
		// Act
		f := run(rect(0, 0, 100, 50, true), rect(-5, -5, 110, 60, true), rect(1, 1, 2, 2, true))

		// Assert
		Expect(f).To(HaveLen(2))
		Expect(f[0].Message).To(Equal("2 white background shapes (largest 110.0 × 60.0 mm) are scored around their edges like any other shape"))
		Expect(f[1].Message).To(HavePrefix("1 white-filled shape is scored like any other shape"))
	})
})
