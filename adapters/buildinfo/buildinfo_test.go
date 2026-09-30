package buildinfo_test

import (
	"runtime/debug"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/buildinfo"
)

// module reports a build of the given module version, as go install records it.
func module(version string) func() (*debug.BuildInfo, bool) {
	return func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Path: "github.com/kinglet-dev/laserlint", Version: version}}, true
	}
}

func unreadable() (*debug.BuildInfo, bool) { return nil, false }

var _ = Describe("Version says which version of laserlint is running", func() {
	It("uses the version written in at release time", func() {
		// Act
		v := buildinfo.Version("0.1.0", module("v0.0.9"))

		// Assert
		Expect(v).To(Equal("0.1.0"))
	})

	It("uses the module version go install recorded, without its v", func() {
		// Act
		v := buildinfo.Version("dev", module("v0.1.1"))

		// Assert
		Expect(v).To(Equal("0.1.1"))
	})

	It("shows a build from an untagged commit by its pseudo-version", func() {
		// Act
		v := buildinfo.Version("dev", module("v0.1.1-0.20260930031500-188b14ede61c"))

		// Assert
		Expect(v).To(Equal("0.1.1-0.20260930031500-188b14ede61c"))
	})

	DescribeTable("says dev when there's no version to report",
		func(read func() (*debug.BuildInfo, bool)) {
			// Act
			v := buildinfo.Version("dev", read)

			// Assert
			Expect(v).To(Equal("dev"))
		},
		Entry("a build from a local checkout", module("(devel)")),
		Entry("no module version recorded", module("")),
		Entry("no build information at all", unreadable),
	)
})
