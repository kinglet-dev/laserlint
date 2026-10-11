package spacing

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// open reports whether a line has two free ends that could join others.
func open(l geom.Polyline) bool { return !l.Closed && len(l.Points) > 1 }

// pairEnds pairs ends of open lines that meet, within meet. Ends are found
// with a uniform grid of cells meet wide, so each end only compares with
// ends in the cells around it. A pair needs both ends to pick each other,
// so every pairing is two-way.
func pairEnds(lines []geom.Polyline) map[end]end {
	cells := map[[2]int64][]end{}
	at := func(e end) geom.Point {
		pts := lines[e.line].Points
		if e.last {
			return pts[len(pts)-1]
		}
		return pts[0]
	}
	var ends []end
	for i, l := range lines {
		if open(l) {
			for _, e := range []end{{i, false}, {i, true}} {
				p := at(e)
				k := [2]int64{cell(p.X), cell(p.Y)}
				cells[k] = append(cells[k], e)
				ends = append(ends, e)
			}
		}
	}
	pick := map[end]end{}
	for _, e := range ends {
		p := at(e)
		var near []end
		for dx := int64(-1); dx <= 1; dx++ {
			for dy := int64(-1); dy <= 1; dy++ {
				for _, o := range cells[[2]int64{cell(p.X) + dx, cell(p.Y) + dy}] {
					if q := at(o); o != e && math.Hypot(q.X-p.X, q.Y-p.Y) < meet {
						near = append(near, o)
					}
				}
			}
		}
		if len(near) > 0 {
			pick[e] = near[0]
		}
	}
	partner := map[end]end{}
	for e, o := range pick {
		if back, ok := pick[o]; ok && back == e {
			partner[e] = o
		}
	}
	return partner
}

func cell(v float64) int64 { return int64(math.Floor(v / meet)) }
