// Package lineart is the check for designs drawn as thin black strokes.
// Laser software scores both edges of a filled stroke, so a line meant to
// burn once burns as two parallel lines. Drawn as hairlines, each stroke
// would burn once.
package lineart

import (
	"fmt"
	"math"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/fill"
	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/parallel"
	"github.com/kinglet-dev/laserlint/domain/spacing"
)

// Limits on the work (threat model: 5,000,000 samples), how closely the
// fill follows curves, mm, and the share of line, in percent, above which
// the design counts as line art.
const (
	maxSamples = 5_000_000
	maxRows    = 10_000_000
	tolerance  = 0.005
	most       = 50
)

// Check is the line-art check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "line-art" }

// Name is the check's name for people.
func (Check) Name() string { return "Line art" }

// Run notes a design where most scored line is the two edges of black
// strokes thinner than twice the minimum spacing: each line faces another
// within that width, with black fill between them. Compared as shown, to 1%.
func (c Check) Run(in check.Input, s check.Settings) ([]check.Finding, error) {
	f, err := fill.New(in.Design.Shapes, tolerance, maxRows)
	if err != nil {
		return nil, err
	}
	width := 2 * (s.Line + s.Gap)
	var between []geom.Point // halfway between each sample and the line it faces
	var lengths []float64
	total, err := spacing.Facing(in.Lines, spacing.Params{MinDistance: width, Step: width / 8, MaxSamples: maxSamples},
		func(from, to geom.Point, length float64) {
			between = append(between, geom.Point{X: (from.X + to.X) / 2, Y: (from.Y + to.Y) / 2})
			lengths = append(lengths, length)
		})
	if err != nil {
		return nil, err
	}
	dark := make([]bool, len(between))
	parallel.For(len(between), func(i int) { dark[i] = f.Dark(between[i]) })
	var strokes float64
	for i, d := range dark {
		if d {
			strokes += lengths[i]
		}
	}
	percent := math.Round(100 * strokes / total)
	if !(percent > most) { // also NaN, for a design with no line
		return nil, nil
	}
	return []check.Finding{{
		Check:    c.ID(),
		Severity: check.Info,
		Message: fmt.Sprintf("%g%% of scored line is the two edges of thin black strokes (under %.2f mm wide), so each stroke burns as a double line",
			percent, width),
		Fix: "Draw the strokes as single hairlines (a stroke and no fill), for example by tracing the original image again with centerline tracing (Inkscape: Path → Trace Bitmap), then check again.",
	}}, nil
}
