package details_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/check/details"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// square is a closed square with its top-left corner at (x, y).
func square(x, y, side float64) geom.Polyline {
	return geom.Polyline{Points: []geom.Point{pt(x, y), pt(x+side, y), pt(x+side, y+side), pt(x, y+side)}, Closed: true}
}

// specks is n squares of the given side, 5 mm apart.
func specks(n int, side float64) []geom.Polyline {
	var out []geom.Polyline
	for i := 0; i < n; i++ {
		out = append(out, square(5*float64(i), 0, side))
	}
	return out
}

func run(lines ...geom.Polyline) []check.Finding {
	f, err := details.New().Run(check.Input{Lines: lines}, check.DefaultSettings)
	Expect(err).NotTo(HaveOccurred())
	return f
}

var _ = Describe("The small-details check", func() {
	It("is a check with a stable id and a name for people", func() {
		// Act
		var c check.Check = details.New()

		// Assert
		Expect(c.ID()).To(Equal("small-details"))
		Expect(c.Name()).To(Equal("Small details"))
	})

	It("finds nothing when every closed shape is large enough", func() {
		// Act
		f := run(square(0, 0, 10), square(20, 20, 0.36))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("warns about a closed shape smaller than the detail size, and where it is", func() {
		// Act: a 0.3 mm square is 0.42 mm across, corner to corner.
		f := run(square(0, 0, 10), square(20, 30, 0.3))

		// Assert
		Expect(f).To(HaveLen(1))
		Expect(f[0].Check).To(Equal("small-details"))
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(Equal("1 closed shape under 0.50 mm across burns as a dot (problem above 20); smallest 0.42 mm"))
		Expect(f[0].Fix).To(ContainSubstring("0.50 mm"))
		Expect(f[0].Locations).To(HaveLen(1))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 20.15, 1e-9))
		Expect(f[0].Locations[0].Y).To(BeNumerically("~", 30.15, 1e-9))
	})

	It("compares sizes as shown, so one that shows as the detail size passes", func() {
		// Act: 0.3535 mm sides are 0.49992 mm across, shown as 0.50.
		f := run(square(0, 0, 0.3535))

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("ignores open lines and closed lines with no length", func() {
		// Act
		f := run(geom.Polyline{Points: []geom.Point{pt(0, 0), pt(0.1, 0)}},
			geom.Polyline{Points: []geom.Point{pt(5, 5), pt(5, 5)}, Closed: true})

		// Assert
		Expect(f).To(BeEmpty())
	})

	It("allows small shapes up to the limit as a warning", func() {
		// Act
		f := run(specks(20, 0.1)...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Warning))
		Expect(f[0].Message).To(HavePrefix("20 closed shapes under 0.50 mm across"))
	})

	It("reports more small shapes than the limit as a problem", func() {
		// Act
		f := run(specks(21, 0.1)...)

		// Assert
		Expect(f[0].Severity).To(Equal(check.Problem))
	})

	It("lists at most five places, the smallest shapes first, each place once", func() {
		// Arrange: six 0.3 mm squares, then two tiny ones 0.5 mm apart.
		lines := append(specks(6, 0.3), square(100, 100, 0.05), square(100.5, 100, 0.06))

		// Act
		f := run(lines...)

		// Assert
		Expect(f[0].Locations).To(HaveLen(5))
		Expect(f[0].Locations[0].X).To(BeNumerically("~", 100.025, 1e-9))
		Expect(f[0].Locations[1].Y).To(BeNumerically("~", 0.15, 1e-9)) // one of the six equal squares
	})

	It("takes the detail size and limit from the settings", func() {
		// Arrange
		s := check.DefaultSettings
		s.Detail, s.MaxDetails = 2, 0

		// Act
		f, err := details.New().Run(check.Input{Lines: []geom.Polyline{square(0, 0, 1)}}, s)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(f[0].Severity).To(Equal(check.Problem))
		Expect(f[0].Message).To(Equal("1 closed shape under 2.00 mm across burns as a dot (problem above 0); smallest 1.41 mm"))
	})
})
