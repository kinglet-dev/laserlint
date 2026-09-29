// Package dense is the check for areas packed with too much scored line:
// the heat builds up and chars the material, even where every pair of
// lines keeps its gap.
package dense

import (
	"fmt"
	"math"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/density"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// The window and step of the measurement, mm, and a limit on grid cells
// (a 1 m × 1 m design at 0.5 mm).
const (
	window   = 10.0
	step     = 0.5
	maxCells = 4_000_000
)

// Check is the density check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "density" }

// Name is the check's name for people.
func (Check) Name() string { return "Density" }

// Run reports the densest area: as information up to the warning level, a
// warning above it, and a problem above the limit. Densities are compared
// as shown, to 0.01, so the report never disagrees with its own numbers.
func (c Check) Run(in check.Input, s check.Settings) ([]check.Finding, error) {
	r, err := density.Densest(in.Lines, density.Params{Window: window, Step: step, MaxCells: maxCells})
	if err != nil || r.Density == 0 {
		return nil, err
	}
	shown := math.Round(100*r.Density) / 100
	f := check.Finding{
		Check:    c.ID(),
		Severity: check.Info,
		Message: fmt.Sprintf("densest %g mm area has %.2f mm of line per mm² (lines about %.2f mm apart); warning above %g, problem above %g",
			window, shown, 1/r.Density, s.WarnDensity, s.MaxDensity),
		Fix:       "Simplify or spread out the busiest area, or remove fine detail there.",
		Locations: []geom.Point{r.At},
	}
	if shown > s.WarnDensity {
		f.Severity = check.Warning
	}
	if shown > s.MaxDensity {
		f.Severity = check.Problem
	}
	return []check.Finding{f}, nil
}
