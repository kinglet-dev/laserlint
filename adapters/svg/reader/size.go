package reader

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kinglet-dev/laserlint/domain/geom"
	"github.com/kinglet-dev/laserlint/domain/units"
)

// mmPerPx is the CSS pixel: 96 per inch.
const mmPerPx = 25.4 / 96

// box is an SVG viewBox.
type box struct{ x, y, w, h float64 }

// size works out the document's physical size and the transform from its
// user units to millimetres, following the SVG rules: numbers without units
// are CSS pixels, a missing width or height comes from the viewBox, and
// without a viewBox one user unit is one CSS pixel.
func size(a map[string]string) (width, height float64, toMM geom.Matrix, err error) {
	vb, hasBox, err := viewBox(a["viewBox"])
	if err != nil {
		return 0, 0, geom.Matrix{}, err
	}
	width, err = dimension(a["width"], vb.w, hasBox)
	if err != nil {
		return 0, 0, geom.Matrix{}, err
	}
	height, err = dimension(a["height"], vb.h, hasBox)
	if err != nil {
		return 0, 0, geom.Matrix{}, err
	}
	if !hasBox {
		return width, height, geom.Matrix{A: mmPerPx, D: mmPerPx}, nil
	}
	toMM, err = fit(vb, width, height, strings.TrimSpace(a["preserveAspectRatio"]))
	return width, height, toMM, err
}

// dimension reads width or height in mm; empty falls back to the viewBox size in CSS pixels.
func dimension(value string, fromBox float64, hasBox bool) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		if !hasBox {
			return 0, ErrNoSize
		}
		return fromBox * mmPerPx, nil
	}
	if strings.HasSuffix(value, "%") {
		return 0, ErrNoSize
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		value += "px" // SVG: a number without a unit is in CSS pixels
	}
	l, err := units.ParseLength(value)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrNoSize, err)
	}
	if l.Millimetres() <= 0 {
		return 0, ErrNoSize
	}
	return l.Millimetres(), nil
}

// viewBox parses "min-x min-y width height".
func viewBox(value string) (box, bool, error) {
	fields := strings.FieldsFunc(value, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r' })
	if len(fields) == 0 {
		return box{}, false, nil
	}
	if len(fields) != 4 {
		return box{}, false, ErrBadViewBox
	}
	var n [4]float64
	for i, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			return box{}, false, ErrBadViewBox
		}
		n[i] = v
	}
	if n[2] <= 0 || n[3] <= 0 {
		return box{}, false, ErrBadViewBox
	}
	return box{n[0], n[1], n[2], n[3]}, true, nil
}

// fit maps the viewBox onto width × height as preserveAspectRatio says.
func fit(vb box, width, height float64, par string) (geom.Matrix, error) {
	sx, sy := width/vb.w, height/vb.h
	origin := geom.Matrix{A: 1, D: 1, E: -vb.x, F: -vb.y}
	switch par {
	case "none":
		return origin.Then(geom.Matrix{A: sx, D: sy}), nil
	case "", "xMidYMid", "xMidYMid meet":
		s := math.Min(sx, sy)
		centre := geom.Matrix{A: 1, D: 1, E: (width - vb.w*s) / 2, F: (height - vb.h*s) / 2}
		return origin.Then(geom.Matrix{A: s, D: s}).Then(centre), nil
	}
	return geom.Matrix{}, fmt.Errorf("%w: preserveAspectRatio=%q (in Inkscape, set the page to fit the drawing)", ErrUnsupported, par)
}
