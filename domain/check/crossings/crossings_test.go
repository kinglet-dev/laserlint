package crossings_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/crossings"
	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

func line(a, b geom.Point) geom.Polyline { return geom.Polyline{Points: []geom.Point{a, b}} }

// ladder is a 1000 mm line crossed by n short lines 5 mm apart.
func ladder(n int) []geom.Polyline {
	lines := []geom.Polyline{line(pt(0, 10), pt(1000, 10))}
	for i := 0; i < n; i++ {
		x := 2.5 + 5*float64(i)
		lines = append(lines, line(pt(x, 5), pt(x, 15)))
	}
	return lines
}

// twice is a straight line of the given length drawn twice, starting at x.
func twice(x, length float64) []geom.Polyline {
	return []geom.Polyline{line(pt(x, 0), pt(x+length, 0)), line(pt(x+length, 0), pt(x, 0))}
}

func run(lines []geom.Polyline) []check.Finding {
	f, err := crossings.New().Run(check.Input{Lines: lines}, check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The crossings check", func() {
	It("is a check with a stable id and a name for people", func() {
		// Act
		var c check.Check = crossings.New()

		// Assert
		Expect(c.ID()).To(Equal("crossings"))
		Expect(c.Name()).To(Equal("Crossings"))
	})

	It("finds nothing where no lines cross", func() {
		// Act
		f := run([]geom.Polyline{line(pt(0, 0), pt(10, 0)), line(pt(0, 1), pt(10, 1))})

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("warns about a single crossing, and where it is", func() {
		// Act
		f := run([]geom.Polyline{line(pt(0, 0), pt(10, 10)), line(pt(0, 10), pt(10, 0))})

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("crossings"))
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(Equal("1 place where score lines cross or one ends on another, burning twice (problem above 50)"))
		Expect(f[0].Fix).To(ContainSubstring("Path → Union"))
		Expect(f[0].Locations).To(HaveLen(1))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 5, 1e-9))
	})

	It("allows crossings up to the limit as a warning", func() {
		// Act
		f := run(ladder(50))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(HavePrefix("50 places where score lines cross"))
	})

	It("reports more crossings than the limit as a problem", func() {
		// Act
		f := run(ladder(51))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Problem))
	})

	It("counts crossings closer than a line's width as one", func() {
		// Act: two crossings 0.05 mm apart burn as one spot.
		f := run([]geom.Polyline{line(pt(0, 5), pt(10, 5)), line(pt(5, 0), pt(5, 10)), line(pt(5.05, 0), pt(5.05, 10))})

		// Assert
		Expect(f[0].Message).To(HavePrefix("1 place "))
	})

	It("lists at most five places, the busiest first, each place once", func() {
		// Arrange: a busy corner of six crossings within 1 mm, then five single crossings.
		lines := ladder(5)
		for i := 0; i < 6; i++ {
			y := 50 + 0.15*float64(i)
			lines = append(lines, line(pt(100, y), pt(101, y)))
		}
		lines = append(lines, line(pt(100.5, 49), pt(100.5, 52)))

		// Act
		f := run(lines)

		// Assert
		Expect(f[0].Locations).To(HaveLen(5))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 100.5, 1e-9))
	})

	It("warns about line drawn on top of other line up to the limit", func() {
		// Act: 2.04 mm shows as 2.0 mm, the limit.
		f := run(twice(0, 2.04))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(Equal("2.0 mm of score line lies on other line and burns twice (problem above 2 mm)"))
		Expect(f[0].Fix).To(ContainSubstring("duplicate"))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 1.02, 1e-9))
	})

	It("reports more stacked line than the limit as a problem", func() {
		// Act: 2.06 mm shows as 2.1 mm.
		f := run(twice(0, 2.06))

		// Assert
		Expect(f[0].Severity).To(Equal(check.Problem))
	})

	It("ignores stacked line too short to show", func() {
		// Act
		f := run(twice(0, 0.04))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("lists at most five stacked stretches, the longest first", func() {
		// Arrange
		var lines []geom.Polyline
		for i := 0; i < 6; i++ {
			lines = append(lines, twice(float64(20*i), 0.1)...)
		}
		lines = append(lines, twice(200, 1)...)

		// Act
		f := run(lines)

		// Assert
		Expect(f[0].Locations).To(HaveLen(5))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 200.5, 1e-9))
	})

	It("reports crossings and stacked line as separate findings", func() {
		// Act
		f := run(append(twice(50, 3), line(pt(0, 0), pt(10, 10)), line(pt(0, 10), pt(10, 0))))

		// Assert
		Expect(f).To(HaveLen(2))
		Expect(f[0].Message).To(HavePrefix("1 place"))
		Expect(f[1].Severity).To(Equal(check.Problem))
	})

	It("takes its limits from the settings", func() {
		// Arrange
		s := check.DefaultSettings
		s.MaxCrossings, s.MaxStacked = 0, 0.5

		// Act
		f, err := crossings.New().Run(check.Input{Lines: append(twice(50, 0.6), ladder(1)...)}, s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(f[0].Severity).To(Equal(check.Problem))
		Expect(f[0].Message).To(HaveSuffix("(problem above 0)"))
		Expect(f[1].Severity).To(Equal(check.Problem))
		Expect(f[1].Message).To(HaveSuffix("(problem above 0.5 mm)"))
	})

	It("passes on a design too large to search as an error", func() {
		// Act
		_, err := crossings.New().Run(check.Input{Lines: []geom.Polyline{line(pt(0, 0), pt(5000, 5000))}}, check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(crossing.ErrTooComplex))
	})
})
