// Package check is the extension point for laserlint's checks: each check
// looks at the design and its score lines and reports findings. New checks
// are added as new implementations of Check (Strategy pattern).
package check

import (
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Severity says how much a finding matters. Only problems make a design
// not ready to burn.
type Severity int

// Severities, least to most serious.
const (
	Info Severity = iota
	Warning
	Problem
)

// Finding is one result of a check.
type Finding struct {
	Check     string       // the check's stable id
	Severity  Severity     // how much it matters
	Message   string       // what was found
	Fix       string       // how to fix it
	Locations []geom.Point // worst places first, mm from the top-left
}

// Settings are the physical limits the checks apply, in millimetres.
type Settings struct {
	Line      float64 // width of a scored line
	Gap       float64 // minimum clear gap between the edges of two lines
	Detail    float64 // smallest closed shape that survives, across
	MaxClose  float64 // largest share (0–1) of scored line allowed too close
	WarnClose float64 // share (0–1) of scored line too close from which to warn

	WarnDensity float64 // mm of line per mm² in the densest 10 mm area above which to warn
	MaxDensity  float64 // mm of line per mm² in the densest 10 mm area allowed before it is a problem

	MaxDetails   int     // closed shapes smaller than Detail allowed before it is a problem
	MaxCrossings int     // places where lines cross allowed before it is a problem
	MaxStacked   float64 // length of line lying on other line allowed before it is a problem
}

// DefaultSettings suit xTool D1 Pro and P2 diode lasers on plywood.
var DefaultSettings = Settings{Line: 0.10, Gap: 0.25, Detail: 0.5, MaxClose: 0.10, WarnClose: 0.01,
	WarnDensity: 0.9, MaxDensity: 2.0, MaxDetails: 20, MaxCrossings: 50, MaxStacked: 2.0}

// Input is what every check looks at.
type Input struct {
	Design design.Design
	Lines  []geom.Polyline // the score lines, from Design.ScoredLines
}

// Check is one kind of check.
type Check interface {
	ID() string   // stable id, used in JSON output
	Name() string // short name for people, e.g. "Lines too close"
	Run(in Input, s Settings) ([]Finding, error)
}
