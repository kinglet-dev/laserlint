package reader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Laser software such as xTool Creative Space burns hidden content, so the
// reader keeps it: laserlint checks what will be burned, not what is shown.
var _ = DescribeTable("Read keeps hidden shapes, because the laser burns them",
	func(body string) {
		// Act
		d, err := read(svg(mm96, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(HaveLen(1))
	},
	Entry("a shape hidden by its style", `<path d="M0 0 L1 1" style="display:none;fill:none;stroke:#000"/>`),
	Entry("a shape hidden by an attribute", `<path d="M0 0 L1 1" display="none"/>`),
	Entry("a shape in a hidden layer", `<g style="display:none"><rect width="5" height="5"/></g>`),
)
