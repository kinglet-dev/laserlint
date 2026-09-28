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
