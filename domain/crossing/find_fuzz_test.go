package crossing_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/domain/crossing"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzFind checks that any two segments are compared without a crash, that
// stacked length never exceeds the shorter segment, and that every crossing
// lies within both segments' bounding boxes.
// Run with: go test ./domain/crossing -run '^$' -fuzz=FuzzFind
func FuzzFind(f *testing.F) {
	f.Add(0.0, 0.0, 10.0, 10.0, 0.0, 10.0, 10.0, 0.0)
	f.Add(0.0, 0.0, 10.0, 0.0, 6.0, 0.0, 20.0, 0.0)
	f.Add(0.0, 0.0, 10.0, 0.0, 5.0, 5.0, 5.0, 0.0)
	f.Add(1e6, 1e6, 1e6+0.3, 1e6, 1e6+0.1, 1e6-0.2, 1e6+0.1, 1e6+0.2)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by, cx, cy, dx, dy float64) {
		// Arrange
		for _, v := range []float64{ax, ay, bx, by, cx, cy, dx, dy} {
			if math.IsNaN(v) || math.Abs(v) > 1e7 {
				return
			}
		}
		p := geom.Polyline{Points: []geom.Point{{X: ax, Y: ay}, {X: bx, Y: by}}}
		q := geom.Polyline{Points: []geom.Point{{X: cx, Y: cy}, {X: dx, Y: dy}}}
		params := crossing.Params{Merge: 0.1, Cell: 1, MaxSegments: 10, MaxEntries: 100_000}

		// Act
		r, err := crossing.Find([]geom.Polyline{p, q}, params)

		// Assert
		if err != nil {
			return
		}
		shorter := math.Min(math.Hypot(bx-ax, by-ay), math.Hypot(dx-cx, dy-cy))
		if r.Stacked < 0 || r.Stacked > shorter*(1+1e-9)+1e-9 {
			t.Fatalf("stacked %v outside [0, %v]", r.Stacked, shorter)
		}
		for _, c := range r.Crossings {
			if !within(c, ax, ay, bx, by) || !within(c, cx, cy, dx, dy) {
				t.Fatalf("crossing %v outside the segments", c)
			}
		}
	})
}

// within reports whether c lies in the box of a segment, allowing for rounding.
func within(c geom.Point, ax, ay, bx, by float64) bool {
	tol := 1e-6 * (1 + math.Max(math.Abs(ax)+math.Abs(bx), math.Abs(ay)+math.Abs(by)))
	return c.X >= math.Min(ax, bx)-tol && c.X <= math.Max(ax, bx)+tol &&
		c.Y >= math.Min(ay, by)-tol && c.Y <= math.Max(ay, by)+tol
}
