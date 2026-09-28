package geom_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

var _ = Describe("Matrix", func() {
	It("moves a point by the translation part (e, f)", func() {
		// Arrange
		m := geom.Matrix{A: 1, D: 1, E: 3, F: -2}

		// Act
		p := m.Apply(geom.Point{X: 1, Y: 1})

		// Assert
		Expect(p).To(Equal(geom.Point{X: 4, Y: -1}))
	})
})

var _ = Describe("Identity", func() {
	It("leaves a point where it is", func() {
		// Arrange
		p := geom.Point{X: 2.5, Y: -7}

		// Act
		got := geom.Identity().Apply(p)

		// Assert
		Expect(got).To(Equal(p))
	})
})

var _ = Describe("Matrix.Then", func() {
	It("applies the first transform, then the second", func() {
		// Arrange
		scale := geom.Matrix{A: 2, D: 2}
		shift := geom.Matrix{A: 1, D: 1, E: 1}

		// Act
		p := scale.Then(shift).Apply(geom.Point{X: 1, Y: 1})

		// Assert
		Expect(p).To(Equal(geom.Point{X: 3, Y: 2}))
	})

	It("gives a different result when the order is reversed", func() {
		// Arrange
		scale := geom.Matrix{A: 2, D: 2}
		shift := geom.Matrix{A: 1, D: 1, E: 1}

		// Act
		p := shift.Then(scale).Apply(geom.Point{X: 1, Y: 1})

		// Assert
		Expect(p).To(Equal(geom.Point{X: 4, Y: 2}))
	})
})
