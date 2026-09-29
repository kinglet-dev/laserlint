package crossing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = Describe("Group gathers nearby points into places", func() {
	It("counts the points at each place, in a fixed order", func() {
		// Arrange: three points within 1 mm of the first, and one far away.
		points := []geom.Point{pt(9, 9), pt(0.5, 0), pt(0, 0.5), pt(0, 0), pt(0.9, 0.3)}

		// Act
		places := crossing.Group(points, 1)

		// Assert
		Expect(places).To(Equal([]crossing.Place{{At: pt(0, 0), Count: 4}, {At: pt(9, 9), Count: 1}}))
	})
})
