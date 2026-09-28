// Package design is the laser job laserlint checks: the shapes whose
// outlines the laser will score, measured in millimetres.
package design

import "github.com/kinglet-dev/laserlint/domain/geom"

// Design is a whole file at its physical size.
type Design struct {
	Width, Height float64 // millimetres
	Shapes        []Shape
}

// Shape is one drawn element. Laser software scores its path: every outline
// of a filled shape, and the line itself of a hairline (stroke, no fill).
type Shape struct {
	Path      geom.Path // millimetres
	Filled    bool      // has a fill (its outlines are scored, not the fill)
	WhiteFill bool      // filled white, often a background rectangle
}
