package crossing

import (
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Place is a group of points near each other.
type Place struct {
	At    geom.Point // the group's first point
	Count int        // how many points it holds
}

// Group gathers points nearer than r to a place's first point into that
// place. Points are taken in a fixed order (top to bottom, then left to
// right), so the result doesn't depend on the order they were found in.
// Cells of size r are searched 3 × 3 around each point (a spatial hash).
func Group(points []geom.Point, r float64) []Place {
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y || (points[i].Y == points[j].Y && points[i].X < points[j].X)
	})
	cells := map[[2]int][]int{}
	key := func(v float64) int { return int(math.Floor(v / r)) }
	var out []Place
	for _, p := range points {
		i := nearest(cells, out, key(p.X), key(p.Y), p, r)
		if i < 0 {
			i = len(out)
			out = append(out, Place{At: p})
			k := [2]int{key(p.X), key(p.Y)}
			cells[k] = append(cells[k], i)
		}
		out[i].Count++
	}
	return out
}

// merge keeps one point of each place, so a crossing found by several
// pieces counts once.
func merge(points []geom.Point, r float64) []geom.Point {
	var out []geom.Point
	for _, pl := range Group(points, r) {
		out = append(out, pl.At)
	}
	return out
}

// nearest finds a place within r of p in the 3 × 3 cells around, or -1.
func nearest(cells map[[2]int][]int, places []Place, cx, cy int, p geom.Point, r float64) int {
	for x := cx - 1; x <= cx+1; x++ {
		for y := cy - 1; y <= cy+1; y++ {
			for _, i := range cells[[2]int{x, y}] {
				if math.Hypot(places[i].At.X-p.X, places[i].At.Y-p.Y) < r {
					return i
				}
			}
		}
	}
	return -1
}
