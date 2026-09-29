package spacing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

// facing collects what Facing reports, with a 0.7 mm reach.
func facing(lines ...geom.Polyline) (from, to []geom.Point, length float64) {
	p := spacing.Params{MinDistance: 0.7, Step: 0.01, MaxSamples: 1_000_000}
	err := spacing.Facing(lines, p, func(f, t geom.Point, l float64) {
		from, to, length = append(from, f), append(to, t), length+l
	})
	Expect(err).NotTo(HaveOccurred())
	return from, to, length
}

var _ = Describe("Facing finds, for each stretch of line, the nearest other line within reach", func() {
	It("reports nothing for a line on its own", func() {
		// Act
		from, _, _ := facing(line(pt(0, 0), pt(10, 0)))

		// Assert
		Expect(from).To(BeEmpty())
	})

	It("pairs each point of two parallel lines with the point opposite", func() {
		// Act
		from, to, length := facing(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.5), pt(10, 0.5)))

		// Assert
		Expect(length).To(BeNumerically("~", 20, 1e-9))
		for i := range from {
			Expect(to[i].X).To(BeNumerically("~", from[i].X, 1e-9))
			Expect(to[i].Y).To(BeNumerically("~", 0.5-from[i].Y, 1e-9))
		}
	})

	It("pairs the two long sides of one thin closed shape", func() {
		// Act: a 10 × 0.5 mm strip drawn as one outline.
		_, _, length := facing(loop(pt(0, 0), pt(10, 0), pt(10, 0.5), pt(0, 0.5)))

		// Assert: both long sides, less a little near each end.
		Expect(length).To(BeNumerically(">", 17))
		Expect(length).To(BeNumerically("<=", 21))
	})

	It("reports nothing for lines further apart than the reach", func() {
		// Act
		from, _, _ := facing(line(pt(0, 0), pt(10, 0)), line(pt(0, 0.8), pt(10, 0.8)))

		// Assert
		Expect(from).To(BeEmpty())
	})

	It("refuses more samples than the limit", func() {
		// Act
		err := spacing.Facing([]geom.Polyline{line(pt(0, 0), pt(10, 0))},
			spacing.Params{MinDistance: 0.7, Step: 0.01, MaxSamples: 10}, func(_, _ geom.Point, _ float64) {})

		// Assert
		Expect(err).To(MatchError(spacing.ErrTooManySamples))
	})
})
