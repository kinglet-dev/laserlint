package check_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

func pt(x, y float64) geom.Point { return geom.Point{X: x, Y: y} }

var _ = Describe("Places lists where a finding is, for people to look", func() {
	It("keeps the given order, listing points within 1 mm of a listed one once", func() {
		// Act
		p := check.Places([]geom.Point{pt(0, 0), pt(0.9, 0), pt(5, 5), pt(0, 1)})

		// Assert
		Expect(p).To(Equal([]geom.Point{pt(0, 0), pt(5, 5), pt(0, 1)}))
	})

	It("lists at most five places", func() {
		// Act
		p := check.Places([]geom.Point{pt(0, 0), pt(2, 0), pt(4, 0), pt(6, 0), pt(8, 0), pt(10, 0)})

		// Assert
		Expect(p).To(HaveLen(5))
		Expect(p[4]).To(Equal(pt(8, 0)))
	})
})
