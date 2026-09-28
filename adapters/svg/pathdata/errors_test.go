package pathdata_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/svg/pathdata"
)

var _ = DescribeTable("Parse refuses malformed path data",
	func(d string, want error) {
		// Arrange: d and want come from the table entry.

		// Act
		_, err := pathdata.Parse(d)

		// Assert
		Expect(err).To(MatchError(want))
	},
	Entry("a path that doesn't start with a move", "L 1 1", pathdata.ErrNoMoveFirst),
	Entry("an unknown command", "M 0 0 X 1 1", pathdata.ErrUnknownCommand),
	Entry("a number where a command is expected", "5 5", pathdata.ErrNoMoveFirst),
	Entry("a command with a missing number", "M 10", pathdata.ErrBadNumber),
	Entry("a lone minus sign", "M 1 -", pathdata.ErrBadNumber),
	Entry("NaN spelled out", "M NaN 0", pathdata.ErrBadNumber),
	Entry("a number too large for a float", "M 1e999 0", pathdata.ErrOutOfRange),
	Entry("a coordinate beyond ±10,000,000", "M 0 0 L 20000000 0", pathdata.ErrOutOfRange),
	Entry("a relative move that runs past the range", "M 9000000 0 m 9000000 0", pathdata.ErrOutOfRange),
	Entry("a bad number after a horizontal line", "M 0 0 H x", pathdata.ErrBadNumber),
	Entry("a bad number after a vertical line", "M 0 0 V", pathdata.ErrBadNumber),
	Entry("a line with only one number", "M 0 0 L 5", pathdata.ErrBadNumber),
)

var _ = Describe("Parse error messages", func() {
	It("say where in the path data the problem is", func() {
		// Arrange
		d := "M 0 0 X 1 1"

		// Act
		_, err := pathdata.Parse(d)

		// Assert
		Expect(err).To(MatchError(ContainSubstring("at character 7")))
	})
})
