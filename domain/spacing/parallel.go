package spacing

import (
	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/parallel"
)

// nearestAll finds the nearest other line for every sample, on all
// processors; each result goes to its sample's own slot.
func (idx *index) nearestAll() ([]float64, []geom.Point) {
	near, at := make([]float64, len(idx.pieces)), make([]geom.Point, len(idx.pieces))
	parallel.For(len(idx.pieces), func(i int) { near[i], at[i] = idx.nearest(idx.pieces[i]) })
	return near, at
}
