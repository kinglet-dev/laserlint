package units

import "errors"

// Errors returned by ParseLength. Each message says what is wrong and how to fix it.
var (
	ErrMissingUnit   = errors.New("length has no unit: add one of mm, cm, in, pt or px, for example 0.25mm")
	ErrUnknownUnit   = errors.New("unknown unit: use mm, cm, in, pt or px")
	ErrInvalidNumber = errors.New("not a valid number: write a finite number before the unit, for example 0.25mm")
	ErrNegative      = errors.New("length is negative: use zero or a positive value")
	ErrTooLarge      = errors.New("length is larger than 10 km: check the number and unit")
)
