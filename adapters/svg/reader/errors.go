package reader

import "errors"

// Errors returned by Read. Each says what is wrong and how to fix it.
var (
	ErrNoSize          = errors.New("the SVG has no usable size: give the <svg> element a width and height in mm or in, or a viewBox")
	ErrBadViewBox      = errors.New("the viewBox must be four numbers with a positive width and height")
	ErrUnsupported     = errors.New("this SVG feature isn't supported yet")
	ErrNotSVG          = errors.New("the file is not an SVG")
	ErrTooLarge        = errors.New("the file is too large")
	ErrTooDeep         = errors.New("the file is nested too deeply")
	ErrTooManyElements = errors.New("the file has too many elements")
	ErrEntities        = errors.New("the file's DOCTYPE defines entities, which aren't allowed: remove the DOCTYPE or re-save the file in Inkscape")
	ErrBadXML          = errors.New("the file is not well-formed XML: open and re-save it in Inkscape")
)
