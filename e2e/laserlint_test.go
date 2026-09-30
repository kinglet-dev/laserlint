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
		Expect(out).To(ContainSubstring("INFO     Density"))
		Expect(out).To(ContainSubstring("OK       Crossings"))
		Expect(out).To(ContainSubstring("OK       Small details"))
		Expect(out).To(ContainSubstring("OK       Background shape"))
		Expect(out).To(ContainSubstring("OK       Line art"))
		Expect(out).To(HaveSuffix("Ready to burn, with 1 info.\n"))
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
		Expect(out).To(HaveSuffix("Not ready to burn: 1 problem and 1 info.\n"))
	})

	It("warns about crossing lines but still passes, exiting 0", func() {
		// Arrange: an X of two hairlines.
		path := file("cross.svg", svg(`<g fill="none" stroke="#000"><path d="M2 2 L18 18"/><path d="M2 18 L18 2"/></g>`))

		// Act
		out, _, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("WARNING  Crossings"))
		Expect(out).To(ContainSubstring("Where: (10.0, 10.0) mm from the top-left"))
	})

	It("warns about a speck but still passes, exiting 0", func() {
		// Arrange: a 0.2 mm dot beside a clean square.
		path := file("speck.svg", svg(`<rect x="2" y="2" width="10" height="10" fill="none" stroke="#000"/><circle cx="16" cy="16" r="0.1"/>`))

		// Act
		out, _, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("WARNING  Small details"))
		Expect(out).To(ContainSubstring("Where: (16.0, 16.0) mm from the top-left"))
	})

	It("warns about a white background but still passes, exiting 0", func() {
		// Arrange
		path := file("background.svg", svg(`<rect width="20" height="20" fill="rgb(255,255,255)"/><rect x="2" y="2" width="10" height="10" fill="none" stroke="#000"/>`))

		// Act
		out, _, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("WARNING  Background shape"))
		Expect(out).To(ContainSubstring("a white background shape (20.0 × 20.0 mm)"))
	})

	It("notes line art drawn as thin black strokes, exiting 0", func() {
		// Arrange: a black stroke 16 mm long and 0.5 mm wide.
		path := file("stroke.svg", svg(`<rect x="2" y="5" width="16" height="0.5"/>`))

		// Act
		out, _, code := laserlint("", path)

		// Assert
		Expect(code).To(Equal(0))
		Expect(out).To(ContainSubstring("INFO     Line art"))
		Expect(out).To(ContainSubstring("so each stroke burns as a double line"))
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
