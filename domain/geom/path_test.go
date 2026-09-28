package geom_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = Describe("Path.Flatten", func() {
	Context("with straight lines only", func() {
		It("returns the corner points unchanged", func() {
			// Arrange
			path := geom.Path{Subpaths: []geom.Subpath{{
				Start: geom.Point{X: 0, Y: 0},
				Segments: []geom.Segment{
					geom.LineTo(geom.Point{X: 10, Y: 0}),
					geom.LineTo(geom.Point{X: 10, Y: 5}),
				},
			}}}

			// Act
			lines := path.Flatten(0.01)

			// Assert
			Expect(lines).To(HaveLen(1))
			Expect(lines[0].Points).To(Equal([]geom.Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 5}}))
			Expect(lines[0].Closed).To(BeFalse())
		})

		It("marks a closed subpath as closed", func() {
			// Arrange
			path := geom.Path{Subpaths: []geom.Subpath{{
				Start:    geom.Point{X: 0, Y: 0},
				Segments: []geom.Segment{geom.LineTo(geom.Point{X: 4, Y: 0}), geom.LineTo(geom.Point{X: 0, Y: 3})},
				Closed:   true,
			}}}

			// Act
			lines := path.Flatten(0.01)

			// Assert
			Expect(lines[0].Closed).To(BeTrue())
		})
	})
})

var _ = Describe("Path.Transform", func() {
	It("moves the start, control and end points of every segment", func() {
		// Arrange
		path := geom.Path{Subpaths: []geom.Subpath{{
			Start:    geom.Point{X: 1, Y: 1},
			Segments: []geom.Segment{geom.CubicTo(geom.Point{X: 2, Y: 2}, geom.Point{X: 3, Y: 3}, geom.Point{X: 4, Y: 4})},
			Closed:   true,
		}}}
		shift := geom.Matrix{A: 1, D: 1, E: 10, F: 20}

		// Act
		moved := path.Transform(shift)

		// Assert
		sub := moved.Subpaths[0]
		Expect(sub.Start).To(Equal(geom.Point{X: 11, Y: 21}))
		Expect(sub.Segments[0]).To(Equal(geom.CubicTo(geom.Point{X: 12, Y: 22}, geom.Point{X: 13, Y: 23}, geom.Point{X: 14, Y: 24})))
		Expect(sub.Closed).To(BeTrue())
	})
})

var _ = Describe("QuadTo", func() {
	It("returns a cubic that traces the same curve as the quadratic", func() {
		// Arrange
		p0, q, p2 := geom.Point{X: 0, Y: 0}, geom.Point{X: 50, Y: 100}, geom.Point{X: 100, Y: 0}

		// Act
		seg := geom.QuadTo(p0, q, p2)

		// Assert
		for _, t := range []float64{0, 0.25, 0.5, 0.75, 1} {
			u := 1 - t
			want := geom.Point{X: u*u*p0.X + 2*u*t*q.X + t*t*p2.X, Y: u*u*p0.Y + 2*u*t*q.Y + t*t*p2.Y}
			got := cubicAt(p0, seg.C1, seg.C2, seg.To, t)
			Expect(got.X).To(BeNumerically("~", want.X, 1e-9))
			Expect(got.Y).To(BeNumerically("~", want.Y, 1e-9))
		}
	})
})
