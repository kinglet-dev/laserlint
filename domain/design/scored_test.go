package design_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

// open is a path with one unclosed subpath through the given points.
func open(points ...geom.Point) geom.Path {
	sub := geom.Subpath{Start: points[0]}
	for _, p := range points[1:] {
		sub.Segments = append(sub.Segments, geom.LineTo(p))
	}
	return geom.Path{Subpaths: []geom.Subpath{sub}}
}

var _ = Describe("ScoredLines lists the lines the laser will score", func() {
	It("scores a hairline along its path, open as drawn", func() {
		// Arrange
		d := design.Design{Shapes: []design.Shape{{Path: open(pt(0, 0), pt(10, 0))}}}

		// Act
		lines := d.ScoredLines(0.01)

		// Assert
		Expect(lines).To(Equal([]geom.Polyline{{Points: []geom.Point{pt(0, 0), pt(10, 0)}}}))
	})

	It("scores a filled shape's outline closed, because the fill closes it", func() {
		// Arrange
		d := design.Design{Shapes: []design.Shape{{Path: open(pt(0, 0), pt(10, 0), pt(10, 10)), Filled: true}}}

		// Act
		lines := d.ScoredLines(0.01)

		// Assert
		Expect(lines).To(HaveLen(1))
		Expect(lines[0].Closed).To(BeTrue())
	})

	It("keeps a closed hairline closed", func() {
		// Arrange
		square := open(pt(0, 0), pt(1, 0), pt(1, 1))
		square.Subpaths[0].Closed = true
		d := design.Design{Shapes: []design.Shape{{Path: square}}}

		// Act
		lines := d.ScoredLines(0.01)

		// Assert
		Expect(lines[0].Closed).To(BeTrue())
	})

	It("scores every subpath of every shape", func() {
		// Arrange
		two := open(pt(0, 0), pt(1, 0))
		two.Subpaths = append(two.Subpaths, open(pt(5, 5), pt(6, 5)).Subpaths...)
		d := design.Design{Shapes: []design.Shape{{Path: two}, {Path: open(pt(9, 9), pt(9, 8))}}}

		// Act
		lines := d.ScoredLines(0.01)

		// Assert
		Expect(lines).To(HaveLen(3))
	})

	It("follows curves within the tolerance", func() {
		// Arrange
		curve := geom.Path{Subpaths: []geom.Subpath{{Start: pt(0, 0),
			Segments: []geom.Segment{geom.CubicTo(pt(0, 10), pt(10, 10), pt(10, 0))}}}}
		d := design.Design{Shapes: []design.Shape{{Path: curve}}}

		// Act
		coarse, fine := d.ScoredLines(1), d.ScoredLines(0.001)

		// Assert
		Expect(len(fine[0].Points)).To(BeNumerically(">", len(coarse[0].Points)))
	})
})
