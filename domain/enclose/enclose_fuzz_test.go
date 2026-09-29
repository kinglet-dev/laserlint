package enclose_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/domain/enclose"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzSmallest checks that the circle holds every point, and is no larger
// than the points' bounding box's diagonal.
// Run with: go test ./domain/enclose -run '^$' -fuzz=FuzzSmallest
func FuzzSmallest(f *testing.F) {
	f.Add(0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0)
	f.Add(0.0, 0.0, 1.0, 2.0, 2.0, 4.0, 3.0, 6.0)
	f.Add(1e6, 1e6, 1e6+0.3, 1e6, 1e6+0.1, 1e6-0.2, 1e6+0.1, 1e6+0.2)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by, cx, cy, dx, dy float64) {
		// Arrange
		v := []float64{ax, ay, bx, by, cx, cy, dx, dy}
		for _, x := range v {
			if math.IsNaN(x) || math.Abs(x) > 1e7 {
				return
			}
		}
		points := []geom.Point{{X: ax, Y: ay}, {X: bx, Y: by}, {X: cx, Y: cy}, {X: dx, Y: dy}}

		// Act
		c := enclose.Smallest(points)

		// Assert
		x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, p := range points {
			x0, y0, x1, y1 = math.Min(x0, p.X), math.Min(y0, p.Y), math.Max(x1, p.X), math.Max(y1, p.Y)
		}
		tol := 1e-6 * (1 + math.Max(math.Abs(x0)+math.Abs(x1), math.Abs(y0)+math.Abs(y1)))
		if !(c.Diameter <= math.Hypot(x1-x0, y1-y0)+tol) {
			t.Fatalf("diameter %v larger than the bounding box's diagonal", c.Diameter)
		}
		for _, p := range points {
			if !(math.Hypot(p.X-c.Centre.X, p.Y-c.Centre.Y) <= c.Diameter/2+tol) {
				t.Fatalf("point %v outside circle %v", p, c)
			}
		}
	})
}
