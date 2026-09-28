package spacing

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// piece is a short stretch of one line, no longer than the sample step.
// Its midpoint is a sample, and the piece itself is a target for others.
type piece struct {
	line   int
	a, b   geom.Point
	arc    float64 // distance along the line to a, mm
	length float64
}

// lineInfo is what the same-line rule needs to know about a line.
type lineInfo struct {
	closed    bool
	perimeter float64
}

// index finds pieces near a point with a uniform grid (Spatial Hash).
type index struct {
	pieces []piece
	lines  []lineInfo
	cells  map[[2]int][]int32
	cell   float64 // cell size = minimum distance, so a query spans 3 × 3 cells
	window float64 // same-line length ignored either side of a point
}

// segments calls f for each straight side of each line. A side of zero
// length (a repeated point) splits into no pieces, so it adds nothing.
func segments(lines []geom.Polyline, f func(line int, a, b geom.Point, length float64)) {
	for i, l := range lines {
		pts := l.Points
		if l.Closed && len(pts) > 1 {
			pts = append(pts[:len(pts):len(pts)], pts[0])
		}
		for j := 1; j < len(pts); j++ {
			f(i, pts[j-1], pts[j], math.Hypot(pts[j].X-pts[j-1].X, pts[j].Y-pts[j-1].Y))
		}
	}
}

// sampleCount is how many pieces the lines split into at the given step.
func sampleCount(lines []geom.Polyline, step float64) int {
	n := 0
	segments(lines, func(_ int, _, _ geom.Point, length float64) { n += int(math.Ceil(length / step)) })
	return n
}

func newIndex(lines []geom.Polyline, step, minDistance float64) *index {
	idx := &index{lines: make([]lineInfo, len(lines)), cells: map[[2]int][]int32{},
		cell: minDistance, window: math.Pi / 2 * minDistance}
	for i, l := range lines {
		idx.lines[i].closed = l.Closed
	}
	segments(lines, func(line int, a, b geom.Point, length float64) {
		n := int(math.Ceil(length / step))
		for k := 0; k < n; k++ {
			p := piece{line: line, a: lerp(a, b, float64(k)/float64(n)), b: lerp(a, b, float64(k+1)/float64(n)),
				arc: idx.lines[line].perimeter, length: length / float64(n)}
			idx.lines[line].perimeter += p.length
			idx.add(p)
		}
	})
	return idx
}

// add stores a piece in every cell its bounding box touches.
func (idx *index) add(p piece) {
	id := int32(len(idx.pieces))
	idx.pieces = append(idx.pieces, p)
	x0, y0 := idx.key(math.Min(p.a.X, p.b.X)), idx.key(math.Min(p.a.Y, p.b.Y))
	x1, y1 := idx.key(math.Max(p.a.X, p.b.X)), idx.key(math.Max(p.a.Y, p.b.Y))
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			idx.cells[[2]int{x, y}] = append(idx.cells[[2]int{x, y}], id)
		}
	}
}

func (idx *index) key(v float64) int { return int(math.Floor(v / idx.cell)) }

func lerp(a, b geom.Point, t float64) geom.Point {
	return geom.Point{X: a.X + t*(b.X-a.X), Y: a.Y + t*(b.Y-a.Y)}
}
