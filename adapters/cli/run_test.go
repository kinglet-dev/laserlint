package cli_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("laserlint FILE", func() {
	var w *world

	BeforeEach(func() {
		// Arrange (shared)
		w = newWorld()
	})

	It("checks the file and writes the text report, exiting 0 when ready to burn", func() {
		// Act
		code := w.run("coaster.svg")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.opened).To(Equal([]string{"coaster.svg"}))
		Expect(w.checker.design.Width).To(Equal(float64(len("from file"))))
		Expect(w.stdout.String()).To(HavePrefix("laserlint 0.1.0 · coaster.svg · "))
		Expect(w.stderr.String()).To(BeEmpty())
	})

	It("closes the file after reading it", func() {
		// Act
		w.run("coaster.svg")

		// Assert
		Expect(w.files["coaster.svg"].closed).To(BeTrue())
	})

	It("exits 1 when there are problems", func() {
		// Arrange
		w.checker.report = problem

		// Act
		code := w.run("coaster.svg")

		// Assert
		Expect(code).To(Equal(1))
		Expect(w.stdout.String()).To(ContainSubstring("Not ready to burn"))
	})

	It("writes JSON with --json", func() {
		// Act
		code := w.run("--json", "coaster.svg")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.stdout.String()).To(ContainSubstring(`"schema": "kinglet.report/v1"`))
	})

	It("reads standard input for -", func() {
		// Arrange
		w.stdin = "piped svg"

		// Act
		code := w.run("-")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.opened).To(BeEmpty())
		Expect(w.checker.design.Width).To(Equal(float64(len("piped svg"))))
		Expect(w.stdout.String()).To(HavePrefix("laserlint 0.1.0 · standard input · "))
	})

	It("uses the default limits", func() {
		// Act
		w.run("coaster.svg")

		// Assert
		Expect(w.checker.settings.Line).To(Equal(0.10))
		Expect(w.checker.settings.Gap).To(Equal(0.25))
		Expect(w.checker.settings.WarnClose).To(Equal(0.01))
		Expect(w.checker.settings.MaxClose).To(Equal(0.10))
	})
})

var _ = Describe("When laserlint can't check the file", func() {
	var w *world

	BeforeEach(func() {
		// Arrange (shared)
		w = newWorld()
	})

	It("exits 2 and says which file couldn't be opened", func() {
		// Act
		code := w.run("missing.svg")

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.stderr.String()).To(Equal("laserlint: can't open missing.svg: file not found\n"))
		Expect(w.stdout.String()).To(BeEmpty())
	})

	It("exits 2 and explains why the SVG couldn't be read", func() {
		// Arrange
		w.readErr = errors.New("<use>: clones can't be measured")

		// Act
		code := w.run("coaster.svg")

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.stderr.String()).To(Equal("laserlint: coaster.svg: <use>: clones can't be measured\n"))
	})

	It("exits 2 when a check can't run", func() {
		// Arrange
		w.checker.err = errors.New("Lines too close: the design has too much line to measure")

		// Act
		code := w.run("coaster.svg")

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.stderr.String()).To(ContainSubstring("laserlint: coaster.svg: Lines too close: "))
	})

	It("exits 2 when the report can't be written", func() {
		// Arrange
		d := w.deps()
		d.Stdout = brokenPipe{}

		// Act
		code := cliRun([]string{"coaster.svg"}, d)

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.stderr.String()).To(ContainSubstring("laserlint: can't write the report: broken pipe"))
	})
})
