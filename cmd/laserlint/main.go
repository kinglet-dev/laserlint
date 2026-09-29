// Command laserlint checks an SVG for laser score lines that are too close
// together or packed too densely. This is the composition root: it only
// wires the parts together.
package main

import (
	"os"

	"github.com/kinglet-dev/laserlint/adapters/cli"
	"github.com/kinglet-dev/laserlint/adapters/files"
	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check/dense"
	"github.com/kinglet-dev/laserlint/domain/check/tooclose"
)

// version is set at release time with -ldflags "-X main.version=…".
var version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.Deps{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Open:    files.Open,
		Read:    reader.Read,
		Checker: app.NewChecker(tooclose.New(), dense.New()),
		Version: version,
	}))
}
