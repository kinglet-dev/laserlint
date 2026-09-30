package lineart_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/lineart"
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/fill"
	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// rect is a rectangle shape: filled black or white, or a hairline.
func rect(x, y, w, h float64, filled, white bool) design.Shape {
	return design.Shape{Filled: filled, WhiteFill: white, Path: geom.Path{Subpaths: []geom.Subpath{{
		Start:    pt(x, y),
		Segments: []geom.Segment{geom.LineTo(pt(x+w, y)), geom.LineTo(pt(x+w, y+h)), geom.LineTo(pt(x, y+h))},
		Closed:   true,
	}}}}
}

// hairline is an open straight line.
func hairline(a, b geom.Point) design.Shape {
	return design.Shape{Path: geom.Path{Subpaths: []geom.Subpath{{Start: a, Segments: []geom.Segment{geom.LineTo(b)}}}}}
}

func input(shapes ...design.Shape) check.Input {
	d := design.Design{Width: 100, Height: 100, Shapes: shapes}
	return check.Input{Design: d, Lines: d.ScoredLines(0.005)}
}

func run(shapes ...design.Shape) []check.Finding {
	f, err := lineart.New().Run(input(shapes...), check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The line-art check", func() {
	It("is a check with a stable id and a name for people", func() {
		// Act
		var c check.Check = lineart.New()

		// Assert
		Expect(c.ID()).To(Equal("line-art"))
		Expect(c.Name()).To(Equal("Line art"))
	})

	It("finds nothing in an empty design", func() {
		// Act, Assert
		Expect(run()).To(BeEmpty())
	})

	It("notes a design drawn as thin black strokes, which burn as double lines", func() {
		// Act: a stroke 20 mm long and 0.5 mm wide, drawn as a filled shape.
		f := run(rect(0, 0, 20, 0.5, true, false))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("line-art"))
		Expect(f[0].Severity).To(Equal(check.Info))
		Expect(f[0].Message).To(MatchRegexp(`^9\d% of scored line is the two edges of thin black strokes \(under 0\.70 mm wide\), so each stroke burns as a double line$`))
		Expect(f[0].Fix).To(ContainSubstring("hairline"))
	})

	It("finds nothing where the thin strip between lines isn't black", func() {
		// Act: a thin white shape, and two hairlines 0.5 mm apart.
		f := run(rect(0, 0, 20, 0.5, true, true), hairline(pt(0, 10), pt(20, 10)), hairline(pt(0, 10.5), pt(20, 10.5)))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("finds nothing in bold black shapes", func() {
		// Act
		f := run(rect(0, 0, 20, 10, true, false))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("finds nothing when thin strokes are only a small part of the design", func() {
		// Act: the stroke's 41 mm of line against 200 mm of hairlines.
		f := run(rect(0, 0, 20, 0.5, true, false), hairline(pt(0, 20), pt(100, 20)), hairline(pt(0, 40), pt(100, 40)))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("notes line art only when more than half the line, as shown, is strokes", func() {
		// Act: about 40 mm of stroke edges against 77.5 mm and 79.5 mm of line.
		over := run(rect(0, 0, 20, 0.5, true, false), hairline(pt(0, 20), pt(36.5, 20)))
		half := run(rect(0, 0, 20, 0.5, true, false), hairline(pt(0, 20), pt(38.5, 20)))

		// Assert
		Expect(over).To(HaveLen(1))
		Expect(over[0].Message).To(HavePrefix("51% "))
		Expect(half).To(BeEmpty())
	})

	It("scales the stroke width with the line and gap settings", func() {
		// Arrange: 2 × (0.2 + 0.4) = 1.2 mm.
		s := check.DefaultSettings
		s.Line, s.Gap = 0.2, 0.4

		// Act
		f, err := lineart.New().Run(input(rect(0, 0, 20, 1, true, false)), s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(f[0].Message).To(ContainSubstring("(under 1.20 mm wide)"))
	})

	It("passes on a design too large to measure as an error", func() {
		// Act: 10 m of line at 0.0875 mm steps exceeds the sample limit.
		_, err := lineart.New().Run(input(hairline(pt(0, 0), pt(500_000, 0))), check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(spacing.ErrTooManySamples))
	})

	It("passes on shapes too large to work out their fill as an error", func() {
		// Act
		_, err := lineart.New().Run(input(rect(0, 0, 1, 20_000_000, true, false)), check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(fill.ErrTooLarge))
	})
})
