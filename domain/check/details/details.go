// Package details is the check for closed shapes too small to survive:
// the laser scores their outline so tightly that they burn as a dot.
package details

import (
	"fmt"
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/enclose"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Check is the small-details check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "small-details" }

// Name is the check's name for people.
func (Check) Name() string { return "Small details" }

// Run counts closed shapes smaller than the detail size across, measured
// by the smallest circle that holds the shape: a warning up to the limit
// and a problem above it. Sizes are compared as shown, to 0.01 mm. Thin
// but long shapes are left to the lines-too-close check.
func (c Check) Run(in check.Input, s check.Settings) ([]check.Finding, error) {
	var small []enclose.Circle
	for _, l := range in.Lines {
		if !l.Closed {
			continue
		}
		e := enclose.Smallest(l.Points)
		if shown := math.Round(100*e.Diameter) / 100; e.Diameter > 0 && shown < s.Detail {
			small = append(small, e)
		}
	}
	if len(small) == 0 {
		return nil, nil
	}
	sort.SliceStable(small, func(i, j int) bool { return small[i].Diameter < small[j].Diameter })
	f := check.Finding{
		Check:    c.ID(),
		Severity: check.Warning,
		Message: fmt.Sprintf("%s under %.2f mm across burns as a dot (problem above %d); smallest %.2f mm",
			shapes(len(small)), s.Detail, s.MaxDetails, small[0].Diameter),
		Fix:       fmt.Sprintf("Delete specks and slivers, or enlarge details to at least %.2f mm across; where overlapping shapes leave slivers, combine them (Inkscape: Path → Union).", s.Detail),
		Locations: centres(small),
	}
	if len(small) > s.MaxDetails {
		f.Severity = check.Problem
	}
	return []check.Finding{f}, nil
}

func shapes(n int) string {
	if n == 1 {
		return "1 closed shape"
	}
	return fmt.Sprintf("%d closed shapes", n)
}

func centres(circles []enclose.Circle) []geom.Point {
	points := make([]geom.Point, len(circles))
	for i, c := range circles {
		points[i] = c.Centre
	}
	return check.Places(points)
}
