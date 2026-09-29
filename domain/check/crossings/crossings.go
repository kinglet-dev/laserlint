// Package crossings is the check for score lines that cross or lie on each
// other: the laser burns those places twice, leaving dark spots and grooves.
package crossings

import (
	"fmt"
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Limits on the search (threat model: 5,000,000 points), the grid cell, mm,
// how many places a finding lists, and how near crossings are to be listed
// as one place, mm.
const (
	maxSegments  = 5_000_000
	maxEntries   = 10_000_000
	cell         = 1.0
	maxLocations = 5
	samePlace    = 1.0
)

// Check is the crossings check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "crossings" }

// Name is the check's name for people.
func (Check) Name() string { return "Crossings" }

// Run reports crossings and stacked line, each as a warning up to its limit
// and a problem above it. Crossings nearer than a line's width burn as one
// spot, so count as one. Stacked length is compared as shown, to 0.1 mm.
func (c Check) Run(in check.Input, s check.Settings) ([]check.Finding, error) {
	r, err := crossing.Find(in.Lines, crossing.Params{Merge: s.Line, Cell: cell, MaxSegments: maxSegments, MaxEntries: maxEntries})
	if err != nil {
		return nil, err
	}
	var out []check.Finding
	if n := len(r.Crossings); n > 0 {
		out = append(out, check.Finding{
			Check:     c.ID(),
			Severity:  level(n > s.MaxCrossings),
			Message:   fmt.Sprintf("%s where score lines cross or one ends on another, burning twice (problem above %d)", places(n), s.MaxCrossings),
			Fix:       "Combine overlapping shapes into one outline (Inkscape: Path → Union), or trim lines to stop where they meet.",
			Locations: busiest(r.Crossings),
		})
	}
	if shown := math.Round(10*r.Stacked) / 10; shown > 0 {
		out = append(out, check.Finding{
			Check:     c.ID(),
			Severity:  level(shown > s.MaxStacked),
			Message:   fmt.Sprintf("%.1f mm of score line lies on other line and burns twice (problem above %g mm)", shown, s.MaxStacked),
			Fix:       "Delete duplicate lines and shapes, or combine shapes that share an edge into one outline (Inkscape: Path → Union).",
			Locations: longest(r.Stacks),
		})
	}
	return out, nil
}

func level(problem bool) check.Severity {
	if problem {
		return check.Problem
	}
	return check.Warning
}

func places(n int) string {
	if n == 1 {
		return "1 place"
	}
	return fmt.Sprintf("%d places", n)
}

// longest lists the middles of the longest stacked stretches.
func longest(stacks []crossing.Stack) []geom.Point {
	sort.SliceStable(stacks, func(i, j int) bool { return stacks[i].Length > stacks[j].Length })
	var out []geom.Point
	for _, st := range stacks[:min(len(stacks), maxLocations)] {
		out = append(out, st.At)
	}
	return out
}

// busiest lists the places with the most crossings first.
func busiest(points []geom.Point) []geom.Point {
	places := crossing.Group(points, samePlace)
	sort.SliceStable(places, func(i, j int) bool { return places[i].Count > places[j].Count })
	var out []geom.Point
	for _, pl := range places[:min(len(places), maxLocations)] {
		out = append(out, pl.At)
	}
	return out
}
