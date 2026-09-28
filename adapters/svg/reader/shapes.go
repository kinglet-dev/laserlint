package reader

import (
	"fmt"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// shapeKind builds one kind of SVG shape element (Strategy pattern): the
// element's outline in user units, or drawn = false when SVG draws nothing.
// fillable is false for elements with no area, whose fill is never scored.
type shapeKind struct {
	build    func(a attributes) (outline geom.Path, drawn bool, err error)
	fillable bool
}

// shapeKinds lists the supported shape elements by name.
var shapeKinds = map[string]shapeKind{
	"path":     {pathShape, true},
	"rect":     {rectShape, true},
	"circle":   {circleShape, true},
	"ellipse":  {ellipseShape, true},
	"line":     {lineShape, false},
	"polyline": {pointsShape(false), true},
	"polygon":  {pointsShape(true), true},
}

func pathShape(a attributes) (geom.Path, bool, error) {
	p, err := syntax.ParsePath(a["d"])
	return p, err == nil, err
}

func lineShape(a attributes) (geom.Path, bool, error) {
	v, err := a.numbers("x1", "y1", "x2", "y2")
	if err != nil {
		return geom.Path{}, false, err
	}
	return polyPath([]geom.Point{{X: v[0], Y: v[1]}, {X: v[2], Y: v[3]}}, false), true, nil
}

// pointsShape builds a polyline, or a polygon when closed.
func pointsShape(closed bool) func(a attributes) (geom.Path, bool, error) {
	return func(a attributes) (geom.Path, bool, error) {
		points, err := syntax.ParsePoints(a["points"])
		if err != nil {
			return geom.Path{}, false, fmt.Errorf("points: %w", err)
		}
		if len(points) == 0 {
			return geom.Path{}, false, nil
		}
		return polyPath(points, closed), true, nil
	}
}

// polyPath joins points with straight lines.
func polyPath(points []geom.Point, closed bool) geom.Path {
	sub := geom.Subpath{Start: points[0], Closed: closed}
	for _, p := range points[1:] {
		sub.Segments = append(sub.Segments, geom.LineTo(p))
	}
	return geom.Path{Subpaths: []geom.Subpath{sub}}
}
