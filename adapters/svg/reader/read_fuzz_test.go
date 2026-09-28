package reader_test

import (
	"math"
	"strings"
	"testing"

	"github.com/kinglet-dev/laserlint/adapters/svg/reader"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzRead checks that any document either fails cleanly or yields a design
// with a positive, finite size and every coordinate within ±10 km.
// Run with: go test ./adapters/svg/reader -run '^$' -fuzz=FuzzRead
func FuzzRead(f *testing.F) {
	for _, seed := range []string{
		svg(mm96, `<g transform="rotate(30) scale(2)" style="fill:none;stroke:#000"><path d="M0 0 A5 5 0 1 1 10 0z"/></g>`),
		svg(`width="4in" height="100" viewBox="0,0,10,20" preserveAspectRatio="none"`, `<rect width="5" height="5" rx="1"/><circle r="2"/>`),
		svg(`viewBox="0 0 10 10"`, `<ellipse rx="3"/><line x2="4" stroke="red"/><polyline points="0 0 1 1"/><polygon points="1,1 2,2 3,1"/>`),
		`<!DOCTYPE svg><svg width="1mm" height="1mm"><defs><path d="M0 0"/></defs><x:a xmlns:x="urn:x"/></svg>`,
		svg(mm96, `<use href="#a"/>`), `<svg width="1" height="1"><g><g><g/></g></g></svg>`, "", "<svg",
	} {
		f.Add(seed)
	}
	limits := reader.Limits{Bytes: 100_000, Depth: 64, Elements: 1_000, Pieces: 10_000}
	f.Fuzz(func(t *testing.T, doc string) {
		// Act
		d, err := reader.ReadLimited(strings.NewReader(doc), limits)

		// Assert
		if err != nil {
			return
		}
		if !(d.Width > 0 && d.Height > 0 && d.Width <= 1e7 && d.Height <= 1e7) {
			t.Fatalf("Read(%q): bad size %v × %v", doc, d.Width, d.Height)
		}
		inRange := func(p geom.Point) bool { return math.Abs(p.X) <= 1e7 && math.Abs(p.Y) <= 1e7 }
		for _, s := range d.Shapes {
			for _, sub := range s.Path.Subpaths {
				if !inRange(sub.Start) {
					t.Fatalf("Read(%q): start out of range %v", doc, sub.Start)
				}
				for _, seg := range sub.Segments {
					if !inRange(seg.C1) || !inRange(seg.C2) || !inRange(seg.To) {
						t.Fatalf("Read(%q): segment out of range %+v", doc, seg)
					}
				}
			}
		}
	})
}
