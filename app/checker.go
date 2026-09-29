// Package app is laserlint's use case: check a design against the
// physical limits and say whether it is ready to burn.
package app

import (
	"fmt"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/design"
)

// tolerance is how closely score lines follow curves, in mm: far below
// the width of a laser line.
const tolerance = 0.005

// Result is one check's findings.
type Result struct {
	ID, Name string
	Findings []check.Finding
}

// Report is the outcome of checking one design.
type Report struct {
	Width, Height float64 // mm
	Settings      check.Settings
	Checks        []Result // in the order the checks ran
}

// Count is the number of findings of a severity, over all checks.
func (r Report) Count(s check.Severity) int {
	n := 0
	for _, c := range r.Checks {
		for _, f := range c.Findings {
			if f.Severity == s {
				n++
			}
		}
	}
	return n
}

// Ready says whether the design is ready to burn: no problems, though
// warnings and information are allowed.
func (r Report) Ready() bool { return r.Count(check.Problem) == 0 }

// Checker runs a fixed list of checks, given at construction (the
// composition root chooses them).
type Checker struct{ checks []check.Check }

// NewChecker returns a checker that runs the given checks in order.
func NewChecker(checks ...check.Check) *Checker { return &Checker{checks: checks} }

// Check runs every check on the design.
func (c *Checker) Check(d design.Design, s check.Settings) (Report, error) {
	in := check.Input{Design: d, Lines: d.ScoredLines(tolerance)}
	r := Report{Width: d.Width, Height: d.Height, Settings: s}
	for _, ch := range c.checks {
		f, err := ch.Run(in, s)
		if err != nil {
			return Report{}, fmt.Errorf("%s: %w", ch.Name(), err)
		}
		r.Checks = append(r.Checks, Result{ID: ch.ID(), Name: ch.Name(), Findings: f})
	}
	return r, nil
}
