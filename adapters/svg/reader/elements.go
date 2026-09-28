package reader

import (
	"encoding/xml"
	"fmt"
)

// svgNamespace is the SVG namespace. Elements outside it (an editor's own
// elements, RDF metadata) are skipped with everything inside them.
const svgNamespace = "http://www.w3.org/2000/svg"

// skipped elements are never drawn directly, so their contents are skipped.
var skipped = map[string]bool{
	"defs": true, "symbol": true, "clipPath": true, "mask": true, "pattern": true,
	"marker": true, "linearGradient": true, "radialGradient": true, "filter": true,
	"title": true, "desc": true, "metadata": true,
}

// refused elements draw something laserlint can't measure; each has a fix.
var refused = map[string]string{
	"use":   "clones can't be measured: in Inkscape, Edit → Clone → Unlink Clone",
	"text":  "text can't be measured: in Inkscape, Path → Object to Path",
	"image": "images can't be measured: remove the image",
	"style": "CSS style sheets aren't supported: save with presentation attributes " +
		"(Illustrator: SVG Options → CSS Properties → Presentation Attributes)",
	"svg": "nested <svg> elements aren't supported: move their contents into a group",
}

// unsupported is an ErrUnsupported with its own message and fix.
type unsupported string

func (u unsupported) Error() string        { return string(u) }
func (u unsupported) Is(target error) bool { return target == ErrUnsupported }

// classify decides whether an element is read, skipped with its contents,
// or refused.
func (w *walker) classify(el xml.StartElement) (skip bool, err error) {
	name := el.Name.Local
	switch {
	case len(w.stack) == 0 && name != "svg":
		return false, fmt.Errorf("%w: its root element is %s", ErrNotSVG, describe(el))
	case len(w.stack) == 0:
		return false, nil
	case el.Name.Space != svgNamespace && el.Name.Space != "":
		return true, nil
	case name == "g" || shapeKinds[name].build != nil:
		return false, nil
	case skipped[name]:
		return true, nil
	case refused[name] != "":
		return false, fmt.Errorf("%s: %w", describe(el), unsupported(refused[name]))
	}
	return false, fmt.Errorf("%s: %w", describe(el),
		unsupported("this element isn't supported: convert it to paths in your editor"))
}
