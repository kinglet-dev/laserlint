package check

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// How many places a finding lists, and how near, in mm, two points are to
// count as one place.
const (
	maxPlaces = 5
	samePlace = 1.0
)

// Places lists up to five of the points, worst first as given, leaving out
// any within 1 mm of one already listed.
func Places(points []geom.Point) []geom.Point {
	var out []geom.Point
	for _, p := range points {
		if len(out) == maxPlaces {
			break
		}
		if !listed(out, p) {
			out = append(out, p)
		}
	}
	return out
}

func listed(points []geom.Point, p geom.Point) bool {
	for _, q := range points {
		if math.Hypot(q.X-p.X, q.Y-p.Y) < samePlace {
			return true
		}
	}
	return false
}
