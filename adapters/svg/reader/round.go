package reader

import (
	"math"

	"github.com/kinglet-dev/laserlint/domain/geom"
)

func rectShape(a attributes, _ int) (geom.Path, bool, error) {
	v, err := a.numbers("x", "y", "width", "height")
	if err != nil {
		return geom.Path{}, false, err
	}
	x, y, w, h := v[0], v[1], v[2], v[3]
	rx, ry, err := a.radii()
	if err != nil || w <= 0 || h <= 0 {
		return geom.Path{}, false, err
	}
	rx, ry = math.Min(rx, w/2), math.Min(ry, h/2)
	if rx == 0 || ry == 0 {
		return polyPath([]geom.Point{{X: x, Y: y}, {X: x + w, Y: y}, {X: x + w, Y: y + h}, {X: x, Y: y + h}}, true), true, nil
	}
	o := startAt(geom.Point{X: x + rx, Y: y})
	o.line(geom.Point{X: x + w - rx, Y: y})
	o.arc(rx, ry, geom.Point{X: x + w, Y: y + ry})
	o.line(geom.Point{X: x + w, Y: y + h - ry})
	o.arc(rx, ry, geom.Point{X: x + w - rx, Y: y + h})
	o.line(geom.Point{X: x + rx, Y: y + h})
	o.arc(rx, ry, geom.Point{X: x, Y: y + h - ry})
	o.line(geom.Point{X: x, Y: y + ry})
	o.arc(rx, ry, geom.Point{X: x + rx, Y: y})
	return o.closed(), true, nil
}

func circleShape(a attributes, _ int) (geom.Path, bool, error) {
	v, err := a.numbers("cx", "cy", "r")
	if err != nil || v[2] <= 0 {
		return geom.Path{}, false, err
	}
	return ellipsePath(v[0], v[1], v[2], v[2]), true, nil
}

func ellipseShape(a attributes, _ int) (geom.Path, bool, error) {
	v, err := a.numbers("cx", "cy")
	if err != nil {
		return geom.Path{}, false, err
	}
	rx, ry, err := a.radii()
	if err != nil || rx <= 0 || ry <= 0 {
		return geom.Path{}, false, err
	}
	return ellipsePath(v[0], v[1], rx, ry), true, nil
}

// ellipsePath is four quarter arcs, clockwise on screen from the right.
func ellipsePath(cx, cy, rx, ry float64) geom.Path {
	o := startAt(geom.Point{X: cx + rx, Y: cy})
	o.arc(rx, ry, geom.Point{X: cx, Y: cy + ry})
	o.arc(rx, ry, geom.Point{X: cx - rx, Y: cy})
	o.arc(rx, ry, geom.Point{X: cx, Y: cy - ry})
	o.arc(rx, ry, geom.Point{X: cx + rx, Y: cy})
	return o.closed()
}

// outline builds one closed subpath from lines and arcs (Builder pattern).
type outline struct {
	sub geom.Subpath
	at  geom.Point
}

func startAt(p geom.Point) *outline { return &outline{sub: geom.Subpath{Start: p}, at: p} }

// line adds a straight side, leaving it out when it has no length.
func (o *outline) line(p geom.Point) {
	if p != o.at {
		o.sub.Segments = append(o.sub.Segments, geom.LineTo(p))
	}
	o.at = p
}

// arc adds a quarter of an ellipse, clockwise on screen.
func (o *outline) arc(rx, ry float64, p geom.Point) {
	o.sub.Segments = append(o.sub.Segments, geom.ArcTo(o.at, rx, ry, 0, false, true, p)...)
	o.at = p
}

func (o *outline) closed() geom.Path {
	o.sub.Closed = true
	return geom.Path{Subpaths: []geom.Subpath{o.sub}}
}
