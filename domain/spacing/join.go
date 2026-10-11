package spacing

import (
	"slices"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// meet is how near, in mm, two ends must be to count as joined: the same
// tolerance the crossings check uses for a point on a line.
const meet = 1e-6

// end is one end of an open line: its first point, or its last.
type end struct {
	line int
	last bool
}

// join merges open lines that meet end to end into longer lines, so a
// split path is measured like the one line it draws (the line merging of
// JTS's LineMerger). A chain that comes back to its start is closed. Where
// three or more ends meet, which two join doesn't change the measurement:
// every line there passes through the meeting point, so the stretch near
// it is close to another line either way.
func join(lines []geom.Polyline) []geom.Polyline {
	partner := pairEnds(lines)
	seen := make([]bool, len(lines))
	var out []geom.Polyline
	for i, l := range lines {
		if seen[i] {
			continue
		}
		if l.Closed {
			seen[i] = true
			out = append(out, l)
			continue
		}
		out = append(out, walk(lines, partner, head(partner, i), seen))
	}
	return out
}

// head walks back from line i's first end to the start of its chain, or
// to i itself when the chain is a loop. It returns the end the chain
// enters by.
func head(partner map[end]end, i int) end {
	in := end{i, false}
	for {
		p, ok := partner[in]
		if !ok || p.line == i {
			return in
		}
		in = end{p.line, !p.last}
	}
}

// walk follows the chain from the end it enters by, joining its lines.
// Each joint keeps both copies of its point; a repeated point adds no line.
func walk(lines []geom.Polyline, partner map[end]end, in end, seen []bool) geom.Polyline {
	first := in
	var out geom.Polyline
	for {
		seen[in.line] = true
		pts := lines[in.line].Points
		if in.last {
			pts = slices.Clone(pts)
			slices.Reverse(pts)
		}
		out.Points = append(out.Points, pts...)
		p, ok := partner[end{in.line, !in.last}]
		if !ok {
			return out
		}
		if p == first {
			out.Closed = true
			return out
		}
		in = p
	}
}
