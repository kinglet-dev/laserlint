package cli_test

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/kinglet-dev/laserlint/adapters/cli"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/design"
)

// fakeChecker is a test double for the checker port.
type fakeChecker struct {
	report   app.Report
	err      error
	settings check.Settings
	design   design.Design
}

func (f *fakeChecker) Check(d design.Design, s check.Settings) (app.Report, error) {
	f.design, f.settings = d, s
	return f.report, f.err
}

// file is an opened test file that records whether it was closed.
type file struct {
	io.Reader
	closed bool
}

func (f *file) Close() error { f.closed = true; return nil }

// world is everything a run touches, all in memory.
type world struct {
	stdout, stderr bytes.Buffer
	stdin          string
	files          map[string]*file
	opened         []string
	readErr        error
	checker        fakeChecker
}

func newWorld() *world {
	return &world{files: map[string]*file{"coaster.svg": {Reader: strings.NewReader("from file")}}}
}

var errNotFound = errors.New("file not found")

func (w *world) deps() cli.Deps {
	return cli.Deps{
		Stdin: strings.NewReader(w.stdin), Stdout: &w.stdout, Stderr: &w.stderr, Version: "0.1.0",
		Open: func(name string) (io.ReadCloser, error) {
			w.opened = append(w.opened, name)
			if f, ok := w.files[name]; ok {
				return f, nil
			}
			return nil, errNotFound
		},
		Read: func(r io.Reader) (design.Design, error) {
			b, _ := io.ReadAll(r)
			return design.Design{Width: float64(len(b))}, w.readErr
		},
		Checker: &w.checker,
	}
}

func (w *world) run(args ...string) int { return cli.Run(args, w.deps()) }

var problem = app.Report{Checks: []app.Result{{Name: "Lines too close", Findings: []check.Finding{{Severity: check.Problem}}}}}
