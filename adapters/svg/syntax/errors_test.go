package syntax_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

var _ = DescribeTable("ParsePath refuses malformed path data",
	func(d string, want error) {
		// Arrange: d and want come from the table entry.

		// Act
		_, err := syntax.ParsePath(d)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("a path that doesn't start with a move", "L 1 1", syntax.ErrNoMoveFirst),
	Entry("an unknown command", "M 0 0 X 1 1", syntax.ErrUnknownCommand),
	Entry("a number where a command is expected", "5 5", syntax.ErrNoMoveFirst),
	Entry("a command with a missing number", "M 10", syntax.ErrBadNumber),
	Entry("a lone minus sign", "M 1 -", syntax.ErrBadNumber),
	Entry("NaN spelled out", "M NaN 0", syntax.ErrBadNumber),
	Entry("a number too large for a float", "M 1e999 0", syntax.ErrOutOfRange),
	Entry("a coordinate beyond ±10,000,000", "M 0 0 L 20000000 0", syntax.ErrOutOfRange),
	Entry("a relative move that runs past the range", "M 9000000 0 m 9000000 0", syntax.ErrOutOfRange),
	Entry("a bad number after a horizontal line", "M 0 0 H x", syntax.ErrBadNumber),
	Entry("a bad number after a vertical line", "M 0 0 V", syntax.ErrBadNumber),
	Entry("a line with only one number", "M 0 0 L 5", syntax.ErrBadNumber),
)

var _ = Describe("ParsePath error messages", func() {
	It("say where in the path data the problem is", func() {
		// Arrange
		d := "M 0 0 X 1 1"

		// Act
		_, err := syntax.ParsePath(d)

		// Assert
		Expect(err).To(MatchError(ContainSubstring("at character 7")))
	})
})
