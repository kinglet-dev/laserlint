package crossing

import (
	"fmt"
	"math"
	"sort"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// segment is one straight piece of a line.
type segment struct{ a, b geom.Point }

// collect splits the lines into their non-empty straight pieces.
func collect(lines []geom.Polyline) []segment {
	var out []segment
	for _, l := range lines {
		pts := l.Points
		if l.Closed && len(pts) > 1 {
			pts = append(pts[:len(pts):len(pts)], pts[0])
		}
		for i := 1; i < len(pts); i++ {
			if pts[i] != pts[i-1] {
				out = append(out, segment{a: pts[i-1], b: pts[i]})
			}
		}
	}
	return out
}

// grid buckets segments by the cells their bounding boxes touch.
type grid struct {
	cell  float64
	segs  []segment
	cells map[[2]int][]int32
}

func newGrid(segs []segment, cell float64, maxEntries int) (*grid, error) {
	g := &grid{cell: cell, segs: segs, cells: map[[2]int][]int32{}}
	entries := 0
	for i, s := range segs {
		x0, y0, x1, y1 := g.box(s)
		entries += (x1 - x0 + 1) * (y1 - y0 + 1)
		if entries > maxEntries {
			return nil, fmt.Errorf("%w: its lines cover more than %d grid cells", ErrTooComplex, maxEntries)
		}
		for x := x0; x <= x1; x++ {
			for y := y0; y <= y1; y++ {
				g.cells[[2]int{x, y}] = append(g.cells[[2]int{x, y}], int32(i))
			}
		}
	}
	return g, nil
}

func (g *grid) key(v float64) int { return int(math.Floor(v / g.cell)) }

// box is the range of cells a segment's bounding box covers.
func (g *grid) box(s segment) (x0, y0, x1, y1 int) {
	return g.key(math.Min(s.a.X, s.b.X)), g.key(math.Min(s.a.Y, s.b.Y)),
		g.key(math.Max(s.a.X, s.b.X)), g.key(math.Max(s.a.Y, s.b.Y))
}

// eachPair visits every pair of segments sharing a cell exactly once: in
// the cell holding the lower corner of their bounding boxes' overlap.
// Cells are visited in a fixed order, so results don't vary between runs.
func (g *grid) eachPair(visit func(a, b segment)) {
	keys := make([][2]int, 0, len(g.cells))
	for k := range g.cells {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i][0] < keys[j][0] || (keys[i][0] == keys[j][0] && keys[i][1] < keys[j][1])
	})
	for _, k := range keys {
		ids := g.cells[k]
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				a, b := g.segs[ids[i]], g.segs[ids[j]]
				ax, ay, _, _ := g.box(a)
				bx, by, _, _ := g.box(b)
				if max(ax, bx) == k[0] && max(ay, by) == k[1] {
					visit(a, b)
				}
			}
		}
	}
}
