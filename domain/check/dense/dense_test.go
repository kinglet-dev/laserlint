package dense_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/dense"
	"github.com/kinglet-dev/laserlint/domain/density"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(a, b geom.Point) geom.Polyline { return geom.Polyline{Points: []geom.Point{a, b}} }

// hatch is n 10 mm lines evenly spaced down a 10 mm square: n mm of line per 10 mm².
func hatch(n int) []geom.Polyline {
	var out []geom.Polyline
	for i := 0; i < n; i++ {
		y := (float64(i) + 0.5) * 10 / float64(n)
		out = append(out, line(pt(0, y), pt(10, y)))
	}
	return out
}

func run(lines []geom.Polyline) []check.Finding {
	f, err := dense.New().Run(check.Input{Lines: lines}, check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The density check", func() {
	It("is a check with a stable id and a name for people", func() {
		// Act
		var c check.Check = dense.New()

		// Assert
		Expect(c.ID()).To(Equal("density"))
		Expect(c.Name()).To(Equal("Density"))
	})

	It("finds nothing in an empty design", func() {
		// Act
		f := run(nil)

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("notes the densest area of a sparse design as information", func() {
		// Act
		f := run(hatch(1))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("density"))
		Expect(f[0].Severity).To(Equal(check.Info))
		Expect(f[0].Message).To(Equal("densest 10 mm area has 0.10 mm of line per mm² (lines about 10.00 mm apart); warning above 0.9, problem above 2"))
		Expect(f[0].Fix).NotTo(BeEmpty())
		Expect(f[0].Locations).To(HaveLen(1))
	})

	It("warns above the warning level", func() {
		// Act: 10 lines 1 mm apart.
		f := run(hatch(10))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(HavePrefix("densest 10 mm area has 1.00 mm of line per mm² (lines about 1.00 mm apart)"))
	})

	It("allows a density exactly at the limit", func() {
		// Act: 20 lines 0.5 mm apart.
		f := run(hatch(20))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
	})

	It("reports a problem above the limit, located at the densest area", func() {
		// Act: 30 lines 0.33 mm apart.
		f := run(hatch(30))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Problem))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 5, 0.5))
		Expect(f[0].Locations[0].Y).To(BeNumerically("~", 5, 0.5))
	})

	It("takes its levels from the settings", func() {
		// Arrange
		s := check.DefaultSettings
		s.WarnDensity, s.MaxDensity = 0.05, 0.5

		// Act
		sparse, err1 := dense.New().Run(check.Input{Lines: hatch(1)}, s)
		busy, err2 := dense.New().Run(check.Input{Lines: hatch(10)}, s)

		// Assert
		Expect(err1).NotTo(HaveOccurred())
		Expect(err2).NotTo(HaveOccurred())
		Expect(sparse[0].Severity).To(Equal(check.Warning))
		Expect(busy[0].Severity).To(Equal(check.Problem))
		Expect(busy[0].Message).To(HaveSuffix("warning above 0.05, problem above 0.5"))
	})

	It("passes on a design too large to measure as an error", func() {
		// Act
		_, err := dense.New().Run(check.Input{Lines: []geom.Polyline{line(pt(0, 0), pt(2000, 2000))}}, check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(density.ErrTooLarge))
	})
})
