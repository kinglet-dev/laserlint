package parallel_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/parallel"
)

var _ = Describe("For runs a function for every index, spread over the processors", func() {
	It("calls it exactly once for each index", func() {
		// Arrange
		calls := make([]int, 1000)

		// Act
		parallel.For(len(calls), func(i int) { calls[i]++ })

		// Assert
		for i := range calls {
			Expect(calls[i]).To(Equal(1), "index %d", i)
		}
	})

	It("does nothing for no indexes", func() {
		// Act, Assert
		parallel.For(0, func(int) { Fail("called") })
	})
})
