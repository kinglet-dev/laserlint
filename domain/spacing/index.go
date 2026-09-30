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
	mid    geom.Point
	arc    float64 // distance along the line to a, mm
	length float64
}

// lineInfo is what the same-line rule needs to know about a line.
type lineInfo struct {
	closed    bool
	perimeter float64
}

// index finds line near a point with a uniform grid (Spatial Hash). The
// samples are short pieces; the line they are measured against is split
// into longer targets, each filed in every cell within reach of it, so a
// query only looks in its own cell.
type index struct {
	pieces  []piece // samples, no longer than the step
	targets []piece // the same line in stretches no longer than reach
	lines   []lineInfo
	cells   map[[2]int][]int32 // target ids by cell
	reach   float64            // the minimum distance: how far a query looks
	cell    float64            // cell size, equal to reach
	window  float64            // same-line length ignored either side of a point
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

func newIndex(lines []geom.Polyline, samples int, step, minDistance float64) *index {
	idx := &index{pieces: make([]piece, 0, samples), lines: make([]lineInfo, len(lines)), cells: map[[2]int][]int32{},
		reach: minDistance, cell: minDistance, window: math.Pi / 2 * minDistance}
	for i, l := range lines {
		idx.lines[i].closed = l.Closed
	}
	segments(lines, func(line int, a, b geom.Point, length float64) {
		arc := idx.lines[line].perimeter
		idx.pieces = split(idx.pieces, line, a, b, arc, length, step)
		for _, t := range split(nil, line, a, b, arc, length, minDistance) {
			idx.add(t)
		}
		idx.lines[line].perimeter += length
	})
	return idx
}

// split appends the side ab, which starts arc along its line, cut into
// equal pieces no longer than most.
func split(out []piece, line int, a, b geom.Point, arc, length, most float64) []piece {
	n := int(math.Ceil(length / most))
	for k := 0; k < n; k++ {
		p := piece{line: line, a: lerp(a, b, float64(k)/float64(n)), b: lerp(a, b, float64(k+1)/float64(n)),
			arc: arc + length*float64(k)/float64(n), length: length / float64(n)}
		p.mid = lerp(p.a, p.b, 0.5)
		out = append(out, p)
	}
	return out
}

// add files a target in every cell its box, grown by reach, touches: any
// point within reach of the target lies in one of them.
func (idx *index) add(t piece) {
	id := int32(len(idx.targets))
	idx.targets = append(idx.targets, t)
	x0, x1 := idx.key(math.Min(t.a.X, t.b.X)-idx.reach), idx.key(math.Max(t.a.X, t.b.X)+idx.reach)
	y0, y1 := idx.key(math.Min(t.a.Y, t.b.Y)-idx.reach), idx.key(math.Max(t.a.Y, t.b.Y)+idx.reach)
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
