package crossing

import (
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// merge treats points nearer than r as one, keeping the first of each group
// in a fixed order, so a crossing found by several pieces counts once.
func merge(points []geom.Point, r float64) []geom.Point {
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y || (points[i].Y == points[j].Y && points[i].X < points[j].X)
	})
	cells := map[[2]int][]geom.Point{}
	key := func(v float64) int { return int(math.Floor(v / r)) }
	var out []geom.Point
	for _, p := range points {
		if !near(cells, key(p.X), key(p.Y), p, r) {
			out = append(out, p)
			k := [2]int{key(p.X), key(p.Y)}
			cells[k] = append(cells[k], p)
		}
	}
	return out
}

// near reports a kept point within r of p, looking in the 3 × 3 cells around.
func near(cells map[[2]int][]geom.Point, cx, cy int, p geom.Point, r float64) bool {
	for x := cx - 1; x <= cx+1; x++ {
		for y := cy - 1; y <= cy+1; y++ {
			for _, q := range cells[[2]int{x, y}] {
				if math.Hypot(q.X-p.X, q.Y-p.Y) < r {
					return true
				}
			}
		}
	}
	return false
}
