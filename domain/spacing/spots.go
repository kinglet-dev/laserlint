package spacing

import (
	"math"
	"sort"
)

// spots groups consecutive close pieces of the same line into stretches,
// each located where the other line is nearest, and sorts them nearest first.
func spots(pieces []piece, near []float64, minDistance float64) []Spot {
	var out []Spot
	open := false
	for i, pc := range pieces {
		if near[i] >= minDistance {
			open = false
			continue
		}
		if !open || pieces[i-1].line != pc.line {
			out = append(out, Spot{Distance: math.Inf(1)})
			open = true
		}
		s := &out[len(out)-1]
		s.Length += pc.length
		if near[i] < s.Distance {
			s.Distance, s.At = near[i], lerp(pc.a, pc.b, 0.5)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Distance < out[j].Distance })
	return out
}
