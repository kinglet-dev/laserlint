package files_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/files"
)

var _ = Describe("Open", func() {
	var dir string

	BeforeEach(func() {
		// Arrange (shared)
		dir = GinkgoT().TempDir()
	})

	It("opens a file for reading", func() {
		// Arrange
		name := filepath.Join(dir, "coaster.svg")
		Expect(os.WriteFile(name, []byte("<svg/>"), 0o600)).To(Succeed())

		// Act
		f, err := files.Open(name)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		b, _ := io.ReadAll(f)
		Expect(string(b)).To(Equal("<svg/>"))
		Expect(f.Close()).To(Succeed())
	})

	It("says a missing file doesn't exist, without repeating its name", func() {
		// Act
		_, err := files.Open(filepath.Join(dir, "missing.svg"))

		// Assert
		Expect(err).To(MatchError(fs.ErrNotExist))
		Expect(err.Error()).NotTo(ContainSubstring("missing.svg"))
	})

	It("refuses a folder, which isn't an SVG file", func() {
		// Act
		_, err := files.Open(dir)

		// Assert
		Expect(err).To(MatchError(files.ErrFolder))
	})
})
