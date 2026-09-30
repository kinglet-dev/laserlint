package fill_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/fill"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// square is a closed square subpath, anticlockwise or not.
func square(x, y, side float64, reverse bool) geom.Subpath {
	pts := []geom.Point{pt(x+side, y), pt(x+side, y+side), pt(x, y+side)}
	if reverse {
		pts = []geom.Point{pt(x, y+side), pt(x+side, y+side), pt(x+side, y)}
	}
	var segs []geom.Segment
	for _, p := range pts {
		segs = append(segs, geom.LineTo(p))
	}
	return geom.Subpath{Start: pt(x, y), Segments: segs, Closed: true}
}

func shape(filled, white bool, subpaths ...geom.Subpath) design.Shape {
	return design.Shape{Filled: filled, WhiteFill: white, Path: geom.Path{Subpaths: subpaths}}
}

func dark(p geom.Point, shapes ...design.Shape) bool {
	f, err := fill.New(shapes, 0.01, 1_000_000)
	Expect(err).NotTo(HaveOccurred())
	return f.Dark(p)
}

var _ = Describe("Fill tells whether the design shows dark fill at a point", func() {
	It("shows nothing where there are no shapes", func() {
		// Act, Assert
		Expect(dark(pt(1, 1))).To(BeFalse())
	})

	It("shows dark inside a black shape and nothing outside it", func() {
		// Arrange
		black := shape(true, false, square(0, 0, 10, false))

		// Act, Assert
		Expect(dark(pt(5, 5), black)).To(BeTrue())
		Expect(dark(pt(15, 5), black)).To(BeFalse())
		Expect(dark(pt(-5, 5), black)).To(BeFalse())
		Expect(dark(pt(5, 15), black)).To(BeFalse())
	})

	It("leaves a hole drawn the other way round empty, as the nonzero rule does", func() {
		// Arrange
		ring := shape(true, false, square(0, 0, 10, false), square(3, 3, 4, true))
		double := shape(true, false, square(0, 0, 10, false), square(3, 3, 4, false))

		// Act, Assert
		Expect(dark(pt(5, 5), ring)).To(BeFalse())
		Expect(dark(pt(1, 5), ring)).To(BeTrue())
		Expect(dark(pt(5, 5), double)).To(BeTrue())
	})

	It("shows the colour of the shape painted last", func() {
		// Arrange
		black, white := shape(true, false, square(0, 0, 10, false)), shape(true, true, square(2, 2, 6, false))

		// Act, Assert
		Expect(dark(pt(5, 5), black, white)).To(BeFalse())
		Expect(dark(pt(1, 5), black, white)).To(BeTrue())
		Expect(dark(pt(5, 5), white, black)).To(BeTrue())
	})

	It("ignores shapes with no fill, which paint only their line", func() {
		// Act, Assert
		Expect(dark(pt(5, 5), shape(false, false, square(0, 0, 10, false)))).To(BeFalse())
	})

	It("handles points on a row boundary and in a row with a flat side", func() {
		// Arrange: a shape spanning several rows; its top side lies in row 0.
		black := shape(true, false, square(0, 0, 10, false))

		// Act, Assert
		Expect(dark(pt(5, 2), black)).To(BeTrue())
		Expect(dark(pt(5, 0.5), black)).To(BeTrue())
		Expect(dark(pt(5, 0.2), black)).To(BeTrue())
		Expect(dark(pt(5, -0.2), black)).To(BeFalse())
	})

	It("gives each point on a shared side to exactly one of two neighbouring shapes", func() {
		// Arrange: a black square with a white one beside it on the right, and one below.
		black := shape(true, false, square(0, 0, 10, false))
		right, below := shape(true, true, square(10, 0, 10, false)), shape(true, true, square(0, 10, 10, false))

		// Act, Assert: left and top sides belong to a shape, right and bottom sides don't.
		Expect(dark(pt(0, 5), black)).To(BeTrue())
		Expect(dark(pt(10, 5), black)).To(BeFalse())
		Expect(dark(pt(5, 0), black)).To(BeTrue())
		Expect(dark(pt(5, 10), black)).To(BeFalse())
		Expect(dark(pt(10, 5), right, black)).To(BeFalse())
		Expect(dark(pt(5, 10), black, below)).To(BeFalse())
	})

	It("refuses shapes spread over more rows than the limit", func() {
		// Act
		_, err := fill.New([]design.Shape{shape(true, false, square(0, 0, 10000, false))}, 0.01, 1000)

		// Assert
		Expect(err).To(MatchError(fill.ErrTooLarge))
	})
})
