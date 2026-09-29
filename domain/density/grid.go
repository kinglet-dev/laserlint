package density

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// grid holds the length of scored line inside each square cell.
type grid struct {
	x0, y0, step float64 // top-left corner and cell size, mm
	nx, ny       int
	cells        []float64
}

// newGrid covers the points with cells aligned to multiples of step.
func newGrid(minP, maxP geom.Point, step float64) *grid {
	g := &grid{x0: math.Floor(minP.X/step) * step, y0: math.Floor(minP.Y/step) * step, step: step}
	g.nx = int((maxP.X-g.x0)/step) + 1
	g.ny = int((maxP.Y-g.y0)/step) + 1
	return g
}

// cell is the index of the cell holding v, measured from origin, clamped.
func (g *grid) cell(v, origin float64, n int) int {
	return min(max(int(math.Floor((v-origin)/g.step)), 0), n-1)
}

// deposit adds a segment's length to the cells it passes through, exactly,
// using the grid traversal of Amanatides and Woo ("A Fast Voxel Traversal
// Algorithm for Ray Tracing", 1987): step from cell to cell at each grid
// line the segment crosses, in order along the segment.
func (g *grid) deposit(a, b geom.Point) {
	length := math.Hypot(b.X-a.X, b.Y-a.Y)
	cx, cy := g.cell(a.X, g.x0, g.nx), g.cell(a.Y, g.y0, g.ny)
	sx, nextX, deltaX := g.axis(a.X, b.X, g.x0, cx)
	sy, nextY, deltaY := g.axis(a.Y, b.Y, g.y0, cy)
	for t := 0.0; t < 1; {
		next := math.Min(math.Min(nextX, nextY), 1)
		g.cells[cy*g.nx+cx] += (next - t) * length
		t = next
		if nextX < nextY {
			cx, nextX = min(max(cx+sx, 0), g.nx-1), nextX+deltaX
		} else {
			cy, nextY = min(max(cy+sy, 0), g.ny-1), nextY+deltaY
		}
	}
}

// axis sets up one axis of the traversal: the step direction, the segment
// parameter t at which it next crosses a grid line, and the t between lines.
func (g *grid) axis(from, to, origin float64, c int) (step int, next, delta float64) {
	d := to - from
	switch {
	case d > 0:
		return 1, (origin + float64(c+1)*g.step - from) / d, g.step / d
	case d < 0:
		return -1, (origin + float64(c)*g.step - from) / d, -g.step / d
	}
	return 0, math.Inf(1), math.Inf(1)
}

// summed returns the summed-area table (integral image, Crow 1984): entry
// (x, y) is the total of all cells above and left of it, so any rectangle's
// total takes four lookups.
func (g *grid) summed() []float64 {
	w := g.nx + 1
	s := make([]float64, w*(g.ny+1))
	for y := 0; y < g.ny; y++ {
		for x := 0; x < g.nx; x++ {
			s[(y+1)*w+x+1] = g.cells[y*g.nx+x] + s[y*w+x+1] + s[(y+1)*w+x] - s[y*w+x]
		}
	}
	return s
}
