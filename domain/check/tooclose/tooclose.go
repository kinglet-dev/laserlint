// Package tooclose is the check for scored lines too close together:
// lines nearer than their width plus a clear gap burn into one smudge.
package tooclose

import (
	"fmt"
	"math"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

// maxSamples bounds the work on huge designs (threat model: 5,000,000).
const maxSamples = 5_000_000

// maxLocations is how many of the nearest spots a finding lists.
const maxLocations = 5

// Check is the lines-too-close check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "lines-too-close" }

// Run reports scored line too close to other scored line: as information
// below the warning level, a warning from it, and a problem over the limit.
// The share is compared as shown, to 0.1%, so the report never disagrees
// with its own numbers at a boundary.
func (c Check) Run(in check.Input, s check.Settings) ([]check.Finding, error) {
	minDist := s.Line + s.Gap
	r, err := spacing.Measure(in.Lines, spacing.Params{MinDistance: minDist, Step: minDist / 8, MaxSamples: maxSamples})
	if err != nil || r.Close == 0 {
		return nil, err
	}
	percent := math.Round(1000*r.Close/r.Scored) / 10
	f := check.Finding{
		Check:    c.ID(),
		Severity: check.Info,
		Message: fmt.Sprintf("%.1f%% of scored line is within %.2f mm of other line, centre to centre (limit %g%%); nearest %.2f mm",
			percent, minDist, 100*s.MaxClose, r.Spots[0].Distance),
		Fix: fmt.Sprintf("Space lines at least %.2f mm apart, centre to centre, or remove fine detail; blunt or shorten very sharp points.", minDist),
	}
	if percent >= 100*s.WarnClose {
		f.Severity = check.Warning
	}
	if percent > 100*s.MaxClose {
		f.Severity = check.Problem
	}
	for i := 0; i < len(r.Spots) && i < maxLocations; i++ {
		f.Locations = append(f.Locations, r.Spots[i].At)
	}
	return []check.Finding{f}, nil
}
