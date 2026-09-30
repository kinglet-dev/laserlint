package e2e_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// readmeExample is the output shown in the README for examples/coaster.svg,
// from the line after the command to the end of the code block.
func readmeExample() []string {
	b, err := os.ReadFile(filepath.Join("..", "README.md"))
	Expect(err).NotTo(HaveOccurred())
	text := strings.ReplaceAll(string(b), "\r\n", "\n") // Windows checkouts may use CRLF
	_, after, found := strings.Cut(text, "$ laserlint coaster.svg\n")
	Expect(found).To(BeTrue(), "the README's example command")
	block, _, found := strings.Cut(after, "\n```")
	Expect(found).To(BeTrue(), "the end of the README's example")
	return strings.Split(block, "\n")
}

var _ = Describe("The README", func() {
	It("shows the real output for the example design, which has problems", func() {
		// Arrange
		want := readmeExample()

		// Act
		out, _, code := laserlint("", filepath.Join("..", "examples", "coaster.svg"))

		// Assert: the first line differs only in version and path.
		got := strings.Split(strings.TrimRight(out, "\n"), "\n")
		Expect(code).To(Equal(1))
		Expect(got[0]).To(HaveSuffix(" · 60.0 × 60.0 mm · line 0.10 mm · gap 0.25 mm"))
		Expect(want[0]).To(Equal("laserlint 0.1.0 · coaster.svg · 60.0 × 60.0 mm · line 0.10 mm · gap 0.25 mm"))
		Expect(got[1:]).To(Equal(want[1:]))
	})
})
