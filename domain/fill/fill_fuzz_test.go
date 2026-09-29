package fill_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/fill"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzDark checks that any triangle is queried without a crash, and that a
// point outside the triangle's bounding box is never dark.
// Run with: go test ./domain/fill -run '^$' -fuzz=FuzzDark
func FuzzDark(f *testing.F) {
	f.Add(0.0, 0.0, 10.0, 0.0, 5.0, 8.0, 5.0, 3.0)
	f.Add(0.0, 0.0, 10.0, 0.0, 5.0, 0.0, 5.0, 0.0)
	f.Add(1e6, 1e6, 1e6+0.3, 1e6, 1e6+0.1, 1e6+0.2, 1e6+0.1, 1e6+0.1)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by, cx, cy, px, py float64) {
		// Arrange
		for _, v := range []float64{ax, ay, bx, by, cx, cy, px, py} {
			if math.IsNaN(v) || math.Abs(v) > 1e7 {
				return
			}
		}
		tri := design.Shape{Filled: true, Path: geom.Path{Subpaths: []geom.Subpath{{
			Start:    geom.Point{X: ax, Y: ay},
			Segments: []geom.Segment{geom.LineTo(geom.Point{X: bx, Y: by}), geom.LineTo(geom.Point{X: cx, Y: cy})},
		}}}}
		fl, err := fill.New([]design.Shape{tri}, 0.01, 100_000)
		if err != nil {
			return
		}

		// Act
		dark := fl.Dark(geom.Point{X: px, Y: py})

		// Assert
		outside := px < math.Min(ax, math.Min(bx, cx)) || px > math.Max(ax, math.Max(bx, cx)) ||
			py < math.Min(ay, math.Min(by, cy)) || py > math.Max(ay, math.Max(by, cy))
		if dark && outside {
			t.Fatalf("point (%v, %v) outside the triangle is dark", px, py)
		}
	})
}
