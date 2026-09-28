// Package units parses and represents physical lengths.
package units

import (
	"math"
	"strconv"
	"strings"
)

// Length is a physical length.
type Length struct {
	mm float64
}

// Millimetres returns the length in millimetres.
func (l Length) Millimetres() float64 { return l.mm }

// unitLength is the length of every supported unit name.
const unitLength = 2

// maxMillimetres (10 km) rejects absurd values before they reach geometry code.
const maxMillimetres = 10_000_000

// millimetresPer lists the supported units. CSS defines 96 px and 72 pt per inch.
var millimetresPer = map[string]float64{
	"mm": 1,
	"cm": 10,
	"in": 25.4,
	"pt": 25.4 / 72,
	"px": 25.4 / 96,
}

// ParseLength parses a length with a unit, such as "0.25mm" or "0.1in".
// A bare number is refused so that a unit is never guessed.
func ParseLength(s string) (Length, error) {
	letters := len(s) - len(strings.TrimRight(s, "abcdefghijklmnopqrstuvwxyz"))
	if letters == 0 {
		return Length{}, ErrMissingUnit
	}
	// Every supported unit is two letters; anything before them is the number.
	split := len(s) - min(letters, unitLength)
	number, unit := s[:split], s[split:]
	perUnit, known := millimetresPer[unit]
	if !known {
		return Length{}, ErrUnknownUnit
	}
	value, err := strconv.ParseFloat(number, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return Length{}, ErrInvalidNumber
	}
	if value < 0 {
		return Length{}, ErrNegative
	}
	mm := value * perUnit
	if mm > maxMillimetres {
		return Length{}, ErrTooLarge
	}
	return Length{mm: mm}, nil
}
