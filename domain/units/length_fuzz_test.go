package units_test

import (
	"math"
	"testing"

	"github.com/kinglet-dev/laserlint/domain/units"
)

// FuzzParseLength checks that any input either fails cleanly or yields a
// finite, non-negative length. Run with: go test ./domain/units -fuzz=FuzzParseLength
func FuzzParseLength(f *testing.F) {
	for _, seed := range []string{"0.25mm", "4in", "", "mm", "3ft", "-1mm", "NaNmm", "1e308in", "0x1p-2mm", "1_0mm"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		// Act
		length, err := units.ParseLength(input)

		// Assert
		if err != nil {
			return
		}
		mm := length.Millimetres()
		if math.IsNaN(mm) || math.IsInf(mm, 0) || mm < 0 {
			t.Fatalf("ParseLength(%q) = %v mm, want a finite, non-negative length", input, mm)
		}
	})
}
