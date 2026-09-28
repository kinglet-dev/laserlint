package syntax_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/adapters/svg/syntax"
)

// FuzzParseTransform checks that any transform list either fails cleanly or
// yields a finite matrix. Run with: go test ./adapters/svg/syntax -fuzz=FuzzParseTransform
func FuzzParseTransform(f *testing.F) {
	for _, seed := range []string{
		"translate(10 20)", "scale(2) rotate(45 5 5)", "matrix(1 2 3 4 5 6)", "skewX(89.9999)",
		"skewY(-90)", "rotate(1e7)", "translate(1e999)", "spin(1)", "scale(10000000) scale(10000000)", "",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, list string) {
		// Act
		m, err := syntax.ParseTransform(list)

		// Assert
		if err != nil {
			return
		}
		for _, v := range [...]float64{m.A, m.B, m.C, m.D, m.E, m.F} {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				t.Fatalf("ParseTransform(%q) = %+v, want only finite entries", list, m)
			}
		}
	})
}
