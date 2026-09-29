// Package cli is laserlint's command line: it parses the arguments, reads
// the SVG, runs the checks and writes the report. Everything it touches is
// passed in (Deps), so the composition root does the wiring and tests run
// in memory.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/kinglet-dev/laserlint/adapters/report"
	"github.com/kinglet-dev/laserlint/adapters/report/jsonreport"
	"github.com/kinglet-dev/laserlint/adapters/report/text"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/design"
)

// Exit codes shared by every Kinglet tool.
const (
	ExitReady    = 0 // ready to burn (warnings allowed)
	ExitProblems = 1 // problems found
	ExitError    = 2 // the file couldn't be checked
)

// Checker is the use case the command line drives.
type Checker interface {
	Check(d design.Design, s check.Settings) (app.Report, error)
}

// Formatter writes a report in one format (Strategy pattern).
type Formatter interface {
	Write(w io.Writer, a report.About, r app.Report) error
}

// Deps is everything a run touches.
type Deps struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Open           func(name string) (io.ReadCloser, error)
	Read           func(r io.Reader) (design.Design, error)
	Checker        Checker
	Version        string
}

// Run runs laserlint with the given arguments (without the program name)
// and returns the exit code.
func Run(args []string, d Deps) int {
	o, err := parse(args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(d.Stdout, help)
		return ExitReady
	}
	if err != nil {
		fmt.Fprintf(d.Stderr, "laserlint: %v\nRun laserlint --help for usage.\n", err)
		return ExitError
	}
	if o.version {
		fmt.Fprintf(d.Stdout, "laserlint %s\n", d.Version)
		return ExitReady
	}
	r, input, err := checkFile(o, d)
	if err != nil {
		fmt.Fprintf(d.Stderr, "laserlint: %v\n", err)
		return ExitError
	}
	var f Formatter = text.New()
	if o.json {
		f = jsonreport.New()
	}
	if err := f.Write(d.Stdout, report.About{Version: d.Version, Input: input}, r); err != nil {
		fmt.Fprintf(d.Stderr, "laserlint: can't write the report: %v\n", err)
		return ExitError
	}
	if !r.Ready() {
		return ExitProblems
	}
	return ExitReady
}

// checkFile reads the file (or standard input) and checks it.
func checkFile(o options, d Deps) (app.Report, string, error) {
	name, in := o.file, d.Stdin
	if name == "-" {
		name = "standard input"
	} else {
		f, err := d.Open(name)
		if err != nil {
			return app.Report{}, "", fmt.Errorf("can't open %s: %w", name, err)
		}
		defer f.Close()
		in = f
	}
	dsn, err := d.Read(in)
	if err != nil {
		return app.Report{}, "", fmt.Errorf("%s: %w", name, err)
	}
	r, err := d.Checker.Check(dsn, o.settings)
	if err != nil {
		return app.Report{}, "", fmt.Errorf("%s: %w", name, err)
	}
	return r, name, nil
}
