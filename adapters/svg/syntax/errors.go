package syntax

import "errors"

// Errors returned by Parse, wrapped with the position of the problem.
var (
	ErrNoMoveFirst    = errors.New("path data must start with a move command (M or m)")
	ErrUnknownCommand = errors.New("unknown path command")
	ErrBadNumber      = errors.New("expected a number")
	ErrOutOfRange     = errors.New("coordinate is outside ±10,000,000")
	ErrBadFlag        = errors.New("expected an arc flag, 0 or 1")
	ErrTooComplex     = errors.New("path data has more pieces than allowed")
)

// Errors returned by ParseTransform.
var (
	ErrUnknownTransform = errors.New("unknown transform function: use matrix, translate, scale, rotate, skewX or skewY")
	ErrTransformArgs    = errors.New("wrong number of arguments for this transform function")
	ErrTransformSyntax  = errors.New("malformed transform: expected name(numbers)")
)

// Errors returned by ParseNumber and ParsePoints.
var (
	ErrUnitInAttribute = errors.New("units aren't supported in shape attributes: give plain numbers in user units (Inkscape: Path → Object to Path)")
	ErrOddCoordinates  = errors.New("points must come in x,y pairs, but there is an odd number of coordinates")
)
