package tooclose_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/tooclose"
	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(a, b geom.Point) geom.Polyline { return geom.Polyline{Points: []geom.Point{a, b}} }

// pair is two 1 mm lines 0.3 mm apart, starting at x.
func pair(x float64) []geom.Polyline {
	return []geom.Polyline{line(pt(x, 0), pt(x+1, 0)), line(pt(x, 0.3), pt(x+1, 0.3))}
}

func run(lines ...geom.Polyline) []check.Finding {
	f, err := tooclose.New().Run(check.Input{Lines: lines}, check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The lines-too-close check", func() {
	It("has a stable id for JSON output", func() {
		// Act
		id := tooclose.New().ID()

		// Assert
		Expect(id).To(Equal("lines-too-close"))
	})

	It("is a check with a name for people", func() {
		// Act
		var c check.Check = tooclose.New()

		// Assert
		Expect(c.Name()).To(Equal("Lines too close"))
	})

	It("finds nothing when every line keeps its distance", func() {
		// Act
		f := run(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.4), pt(10, 0.4)))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("warns about a few close spots, under the limit", func() {
		// Arrange: 2 mm of close line among 102 mm.
		lines := append(pair(0), line(pt(0, 50), pt(100, 50)))

		// Act
		f := run(lines...)

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("lines-too-close"))
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(Equal("2.0% of scored line is within 0.35 mm of other line, centre to centre (limit 10%); nearest 0.30 mm"))
		Expect(f[0].Fix).NotTo(BeEmpty())
	})

	It("only notes close spots below the warning level", func() {
		// Arrange: 2 mm of close line among 302 mm is 0.7%.
		lines := append(pair(0), line(pt(0, 50), pt(300, 50)))

		// Act
		f := run(lines...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Info))
		Expect(f[0].Locations).NotTo(BeEmpty())
	})

	It("warns from exactly the warning level", func() {
		// Arrange: 2 mm close of 200 mm in all is exactly 1%.
		lines := append(pair(0), line(pt(0, 50), pt(198, 50)))

		// Act
		f := run(lines...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
	})

	It("takes the warning level from the settings", func() {
		// Arrange: 2.0% close, with warnings only from 5%.
		s := check.DefaultSettings
		s.WarnClose = 0.05

		// Act
		f, err := tooclose.New().Run(check.Input{Lines: append(pair(0), line(pt(0, 50), pt(100, 50)))}, s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(f[0].Severity).To(Equal(check.Info))
	})

	It("reports a problem when the close share is over the limit", func() {
		// Act
		f := run(pair(0)...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Problem))
		Expect(f[0].Message).To(HavePrefix("100.0% of scored line"))
	})

	It("allows a share exactly at the limit", func() {
		// Arrange: 2 mm close of 20 mm in all is exactly 10%.
		lines := append(pair(0), line(pt(0, 50), pt(18, 50)))

		// Act
		f := run(lines...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
	})

	It("uses the line width plus the gap as the minimum distance", func() {
		// Arrange
		s := check.DefaultSettings
		s.Line, s.Gap = 0.05, 0.2

		// Act
		f, err := tooclose.New().Run(check.Input{Lines: pair(0)}, s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(f).To(BeEmpty())
	})

	It("locates the five nearest spots, nearest first", func() {
		// Arrange: eight close pairs, the one at x = 30 the nearest.
		var lines []geom.Polyline
		for i := 0; i < 8; i++ {
			lines = append(lines, pair(float64(10*i))...)
		}
		lines[7] = line(pt(30, 0.2), pt(31, 0.2))

		// Act
		f := run(lines...)

		// Assert
		Expect(f[0].Locations).To(HaveLen(5))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 30.5, 0.5))
	})

	It("lists each place once, though both lines of a close pair report it", func() {
		// Act
		f := run(pair(0)...)

		// Assert
		Expect(f[0].Locations).To(HaveLen(1))
	})

	It("passes on a design too large to measure as an error", func() {
		// Act
		_, err := tooclose.New().Run(check.Input{Lines: []geom.Polyline{line(pt(0, 0), pt(300_000, 0))}}, check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(spacing.ErrTooManySamples))
	})
})
