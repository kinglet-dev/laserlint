package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("The laserlint binary", func() {
	It("passes a clean design, exiting 0", func() {
		// Arrange
		path := file("clean.svg", clean)

		// Act
		out, errOut, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("laserlint 0.0.0-test · "))
		Expect(out).To(ContainSubstring("· 20.0 × 20.0 mm · line 0.10 mm · gap 0.25 mm"))
		Expect(out).To(ContainSubstring("OK       Lines too close"))
		Expect(out).To(HaveSuffix("Ready to burn.\n"))
		Expect(errOut).To(BeEmpty())
	})

	It("fails a design with lines too close, exiting 1", func() {
		// Arrange
		path := file("close.svg", tooClose)

		// Act
		out, _, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(1))
		Expect(out).To(ContainSubstring("PROBLEM  Lines too close"))
		Expect(out).To(ContainSubstring("100.0% of scored line is within 0.35 mm"))
		Expect(out).To(HaveSuffix("Not ready to burn: 1 problem.\n"))
	})

	It("reads standard input and writes JSON", func() {
		// Act
		out, _, code := laserlint(tooClose, "--json", "-")

		// Assert
		Expect(code).To(Equal(1))
		Expect(out).To(ContainSubstring(`"schema": "kinglet.report/v1"`))
		Expect(out).To(ContainSubstring(`"name": "standard input"`))
		Expect(out).To(ContainSubstring(`"severity": "problem"`))
	})

	It("explains a file it can't check, exiting 2", func() {
		// Arrange
		path := file("clone.svg", svg(`<use href="#a"/>`))

		// Act
		out, errOut, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(2))
		Expect(out).To(BeEmpty())
		Expect(errOut).To(ContainSubstring("clones can't be measured: in Inkscape, Edit → Clone → Unlink Clone"))
	})

	It("explains a missing file, exiting 2", func() {
		// Act
		_, errOut, code := laserlint("", "no-such-file.svg")

		// Assert
		Expect(code).To(Equal(2))
		Expect(errOut).To(HavePrefix("laserlint: can't open no-such-file.svg: "))
	})

	It("shows help, exiting 0", func() {
		// Act
		out, _, code := laserlint("", "--help")

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("Examples:"))
	})
})
