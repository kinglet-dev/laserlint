package density_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/domain/density"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzDensest checks that any segment is measured without a crash, with a
// density between zero and the segment's length over the window's area.
// Run with: go test ./domain/density -run '^$' -fuzz=FuzzDensest
func FuzzDensest(f *testing.F) {
	f.Add(0.25, 1.1, 10.25, 1.1, 0.5)
	f.Add(-3.0, -7.5, 2.0, 4.0, 0.5)
	f.Add(1e6, 1e6, 1e6+0.3, 1e6-0.1, 0.1)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by, step float64) {
		for _, v := range []float64{ax, ay, bx, by} {
			if math.IsNaN(v) || math.Abs(v) > 1e7 {
				return
			}
		}
		if !(step >= 0.01 && step <= 5) {
			return
		}
		p := density.Params{Window: 10 * step, Step: step, MaxCells: 1_000_000}
		line := geom.Polyline{Points: []geom.Point{{X: ax, Y: ay}, {X: bx, Y: by}}}

		// Act
		r, err := density.Densest([]geom.Polyline{line}, p)

		// Assert
		if err != nil {
			return
		}
		limit := math.Hypot(bx-ax, by-ay)/(p.Window*p.Window) + 1e-9
		if r.Density < 0 || r.Density > limit*(1+1e-9) {
			t.Fatalf("density %v outside [0, %v]", r.Density, limit)
		}
	})
}
