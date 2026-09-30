// Package buildinfo works out which version of laserlint is running.
// Release builds have it written in at build time; a copy built with
// go install has none, but Go records the module version it was built from.
package buildinfo

import (
	"runtime/debug"
	"strings"
)

// Version is the stamped version, or else the module version Go recorded
// (without its leading v), or else "dev". read is debug.ReadBuildInfo,
// passed in so tests can supply their own.
func Version(stamped string, read func() (*debug.BuildInfo, bool)) string {
	if stamped != "dev" {
		return stamped
	}
	info, ok := read()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}
	return strings.TrimPrefix(info.Main.Version, "v")
}
