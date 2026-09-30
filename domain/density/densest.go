// Package density finds where a design packs the most scored line into one
// area. Dense line heats and chars the material even when every pair of
// lines keeps its gap.
//
// The measure follows the density rules of chip layout: the most line in any
// square window stepped across the design, in mm of line per mm². Its
// inverse is the average spacing between lines in that window.
package density

import (
	"errors"
	"fmt"
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

// Params sets the measurement.
type Params struct {
	Window   float64 // side of the square window, mm
	Step     float64 // how far the window moves each time (the cell size), mm
	MaxCells int     // limit on grid cells, so huge designs fail clearly
}

// Result is the densest window.
type Result struct {
	Density float64    // mm of line per mm²
	At      geom.Point // centre of the window
}

// ErrTooLarge means the design's lines spread over too many cells to measure.
var ErrTooLarge = errors.New("the design is too large to measure its density")

// Densest finds the window with the most scored line.
func Densest(lines []geom.Polyline, p Params) (Result, error) {
	minP, maxP, ok := bounds(lines)
	if !ok {
		return Result{}, nil
	}
	g := newGrid(minP, maxP, p.Step)
	if g.nx*g.ny > p.MaxCells {
		return Result{}, fmt.Errorf("%w: its lines spread over %.0f × %.0f mm", ErrTooLarge, maxP.X-minP.X, maxP.Y-minP.Y)
	}
	g.cells = make([]float64, g.nx*g.ny)
	for _, l := range lines {
		pts := l.Points
		if l.Closed {
			pts = append(pts[:len(pts):len(pts)], pts[0])
		}
		for i := 1; i < len(pts); i++ {
			g.deposit(pts[i-1], pts[i])
		}
	}
	return g.densest(int(math.Round(p.Window/p.Step)), p.Window), nil
}

// densest slides an m × m-cell window over the grid and keeps the fullest.
// A window must beat the best so far by more than rounding error, so of
// windows with the same line the first, top to bottom and left to right,
// is reported on every processor (sums can round differently on each).
func (g *grid) densest(m int, window float64) Result {
	s, w := g.summed(), g.nx+1
	var best Result
	bestSum := -1.0
	for y := 0; y <= max(g.ny-m, 0); y++ {
		for x := 0; x <= max(g.nx-m, 0); x++ {
			x1, y1 := min(x+m, g.nx), min(y+m, g.ny)
			sum := s[y1*w+x1] - s[y*w+x1] - s[y1*w+x] + s[y*w+x]
			if sum > bestSum+1e-9*(1+bestSum) {
				bestSum = sum
				best = Result{Density: sum / (window * window),
					At: geom.Point{X: g.x0 + (float64(x)+float64(m)/2)*g.step, Y: g.y0 + (float64(y)+float64(m)/2)*g.step}}
			}
		}
	}
	return best
}

// bounds is the box around every point, or ok = false when there are none.
func bounds(lines []geom.Polyline) (minP, maxP geom.Point, ok bool) {
	for _, l := range lines {
		for _, p := range l.Points {
			if !ok {
				minP, maxP, ok = p, p, true
			}
			minP = geom.Point{X: math.Min(minP.X, p.X), Y: math.Min(minP.Y, p.Y)}
			maxP = geom.Point{X: math.Max(maxP.X, p.X), Y: math.Max(maxP.Y, p.Y)}
		}
	}
	return minP, maxP, ok
}
