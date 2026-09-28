package pathdata_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/adapters/svg/pathdata"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// FuzzParse checks that any path data either fails cleanly or yields only
// finite coordinates. Run with: go test ./adapters/svg/pathdata -fuzz=FuzzParse
func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"M 10 20 L 30 40 Z", "m1 1 h5 v5 z", "M0 0 C 1 2 3 4 5 6 S 7 8 9 10",
		"M0 0 Q 5 10 10 0 T 20 0", "M0 0 a5 5 0 0110 0", "M0 0 A 0 5 0 1 1 1e7 0",
		"M-1-2.5.5e1", "L 1 1", "M 1e999 0", "M0 0 A 5 5 0 2 1 10 0", "",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, d string) {
		// Act
		path, err := pathdata.ParseLimited(d, 10_000)

		// Assert
		if err != nil {
			return
		}
		finite := func(p geom.Point) bool {
			return !math.IsNaN(p.X) && !math.IsNaN(p.Y) && !math.IsInf(p.X, 0) && !math.IsInf(p.Y, 0)
		}
		for _, sub := range path.Subpaths {
			if !finite(sub.Start) {
				t.Fatalf("Parse(%q): non-finite start %v", d, sub.Start)
			}
			for _, seg := range sub.Segments {
				if !finite(seg.C1) || !finite(seg.C2) || !finite(seg.To) {
					t.Fatalf("Parse(%q): non-finite segment %+v", d, seg)
				}
			}
		}
	})
}
