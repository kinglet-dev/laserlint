// Package background is the check for white-filled shapes. Laser software
// scores every filled shape's outline whatever its colour, so a white
// background rectangle burns a frame around the design.
package background

import (
	"fmt"
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// cover is the share of the design's area a white shape's bounding box must
// cover to be taken as a background; tolerance is for flattening curves, mm.
const (
	cover     = 0.9
	tolerance = 0.01
)

// Check is the background check.
type Check struct{}

// New returns the check.
func New() Check { return Check{} }

// ID is the check's stable id.
func (Check) ID() string { return "background" }

// Name is the check's name for people.
func (Check) Name() string { return "Background shape" }

// box is a white shape's bounding box.
type box struct{ x0, y0, x1, y1 float64 }

func (b box) area() float64      { return (b.x1 - b.x0) * (b.y1 - b.y0) }
func (b box) centre() geom.Point { return geom.Point{X: (b.x0 + b.x1) / 2, Y: (b.y0 + b.y1) / 2} }

// Run warns about white backgrounds and notes other white shapes.
func (c Check) Run(in check.Input, _ check.Settings) ([]check.Finding, error) {
	var backgrounds, others []box
	for _, s := range in.Design.Shapes {
		if !s.WhiteFill {
			continue
		}
		b := bounds(s.Path)
		if b.area() >= cover*in.Design.Width*in.Design.Height {
			backgrounds = append(backgrounds, b)
		} else {
			others = append(others, b)
		}
	}
	var out []check.Finding
	if len(backgrounds) > 0 {
		out = append(out, check.Finding{
			Check: c.ID(), Severity: check.Warning, Message: backgroundMessage(backgrounds),
			Fix:       "Delete it before burning; hiding it isn't enough, as laser software such as xTool Creative Space burns hidden layers too.",
			Locations: centres(backgrounds),
		})
	}
	if len(others) > 0 {
		out = append(out, check.Finding{
			Check: c.ID(), Severity: check.Info,
			Message:   fmt.Sprintf("%s scored like any other shape: white doesn't mean not burned", count(len(others))),
			Fix:       "Delete any white shapes that aren't meant to burn.",
			Locations: centres(others),
		})
	}
	return out, nil
}

func backgroundMessage(b []box) string {
	if len(b) == 1 {
		return fmt.Sprintf("a white background shape (%.1f × %.1f mm) is scored around its edge like any other shape", b[0].x1-b[0].x0, b[0].y1-b[0].y0)
	}
	sortLargest(b)
	return fmt.Sprintf("%d white background shapes (largest %.1f × %.1f mm) are scored around their edges like any other shape",
		len(b), b[0].x1-b[0].x0, b[0].y1-b[0].y0)
}

func count(n int) string {
	if n == 1 {
		return "1 white-filled shape is"
	}
	return fmt.Sprintf("%d white-filled shapes are", n)
}

func sortLargest(b []box) {
	sort.SliceStable(b, func(i, j int) bool { return b[i].area() > b[j].area() })
}

// centres lists the boxes' centres, largest first.
func centres(b []box) []geom.Point {
	sortLargest(b)
	points := make([]geom.Point, len(b))
	for i, x := range b {
		points[i] = x.centre()
	}
	return check.Places(points)
}

func bounds(p geom.Path) box {
	b := box{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	for _, l := range p.Flatten(tolerance) {
		for _, pt := range l.Points {
			b = box{math.Min(b.x0, pt.X), math.Min(b.y0, pt.Y), math.Max(b.x1, pt.X), math.Max(b.y1, pt.Y)}
		}
	}
	return b
}
