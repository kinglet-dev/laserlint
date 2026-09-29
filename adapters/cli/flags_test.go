package cli_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/cli"
)

var cliRun = cli.Run

type brokenPipe struct{}

func (brokenPipe) Write([]byte) (int, error) { return 0, errors.New("broken pipe") }

var _ = Describe("laserlint's flags", func() {
	var w *world

	BeforeEach(func() {
		// Arrange (shared)
		w = newWorld()
	})

	It("takes the line width and gap in any supported unit", func() {
		// Act
		code := w.run("--line", "0.08mm", "--gap", "0.01in", "coaster.svg")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.checker.settings.Line).To(BeNumerically("~", 0.08, 1e-12))
		Expect(w.checker.settings.Gap).To(BeNumerically("~", 0.254, 1e-12))
	})

	It("takes the warning level and limit as percentages", func() {
		// Act
		code := w.run("--warn-close", "2%", "--max-close", "20%", "coaster.svg")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.checker.settings.WarnClose).To(BeNumerically("~", 0.02, 1e-12))
		Expect(w.checker.settings.MaxClose).To(BeNumerically("~", 0.20, 1e-12))
	})

	It("takes the density levels as mm of line per mm²", func() {
		// Act
		code := w.run("--warn-density", "1.2", "--max-density", "3", "coaster.svg")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.checker.settings.WarnDensity).To(Equal(1.2))
		Expect(w.checker.settings.MaxDensity).To(Equal(3.0))
	})

	DescribeTable("refuses bad values with a message that says how to fix them, exiting 2",
		func(args []string, want string) {
			// Act
			code := w.run(append(args, "coaster.svg")...)

			// Assert
			Expect(code).To(Equal(2))
			Expect(w.stderr.String()).To(ContainSubstring(want))
			Expect(w.opened).To(BeEmpty())
		},
		Entry("a length without a unit", []string{"--line", "0.1"}, "--line: length has no unit"),
		Entry("a gap without a unit", []string{"--gap", "0.25"}, "--gap: length has no unit"),
		Entry("a zero line width", []string{"--line", "0mm"}, "--line must be more than zero"),
		Entry("a percentage without %", []string{"--max-close", "10"}, "--max-close: write a percentage such as 10%"),
		Entry("a percentage over 100%", []string{"--warn-close", "150%"}, "--warn-close: write a percentage from 0% to 100%"),
		Entry("a percentage that isn't a number", []string{"--max-close", "lots%"}, "--max-close: write a percentage such as 10%"),
		Entry("a warning level above the limit", []string{"--warn-close", "20%", "--max-close", "10%"}, "--warn-close must not be more than --max-close"),
		Entry("a density that isn't a number", []string{"--max-density", "2mm"}, "--max-density: write a positive number of mm of line per mm², such as 2"),
		Entry("an infinite density", []string{"--max-density", "Inf"}, "--max-density: write a positive number of mm of line per mm², such as 2"),
		Entry("a zero density", []string{"--warn-density", "0"}, "--warn-density: write a positive number of mm of line per mm², such as 2"),
		Entry("a density warning above the limit", []string{"--warn-density", "3", "--max-density", "2"}, "--warn-density must not be more than --max-density"),
		Entry("an unknown flag", []string{"--colour"}, "flag provided but not defined: -colour"),
	)

	It("needs a file", func() {
		// Act
		code := w.run()

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.stderr.String()).To(Equal("laserlint: give one SVG file to check, or - to read standard input\nRun laserlint --help for usage.\n"))
	})

	It("checks only one file at a time", func() {
		// Act
		code := w.run("coaster.svg", "coaster.svg")

		// Assert
		Expect(code).To(Equal(2))
		Expect(w.opened).To(BeEmpty())
		Expect(w.stderr.String()).To(ContainSubstring("give one SVG file to check"))
	})

	It("shows help with examples on --help, exiting 0", func() {
		// Act
		code := w.run("--help")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.stdout.String()).To(ContainSubstring("Usage:"))
		Expect(w.stdout.String()).To(ContainSubstring("Examples:"))
		Expect(w.stdout.String()).To(ContainSubstring("--max-close"))
		Expect(w.stdout.String()).To(ContainSubstring("--max-density"))
		Expect(w.stderr.String()).To(BeEmpty())
	})

	It("shows the version on --version", func() {
		// Act
		code := w.run("--version")

		// Assert
		Expect(code).To(Equal(0))
		Expect(w.stdout.String()).To(Equal("laserlint 0.1.0\n"))
	})
})
