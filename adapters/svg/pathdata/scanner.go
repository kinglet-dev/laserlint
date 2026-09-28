package pathdata

import (
	"fmt"
	"math"
	"strconv"
)

// maxCoordinate bounds every number and every resulting point.
const maxCoordinate = 10_000_000

// scanner walks path data one token at a time.
type scanner struct {
	d   string
	pos int
}

// fail wraps err with the (1-based) character position of the problem.
func (s *scanner) fail(err error, pos int) error {
	return fmt.Errorf("%w at character %d", err, pos+1)
}

// skipSeparators skips whitespace and commas.
func (s *scanner) skipSeparators() {
	for s.pos < len(s.d) {
		switch s.d[s.pos] {
		case ' ', '\t', '\n', '\r', '\f', ',':
			s.pos++
		default:
			return
		}
	}
}

// done reports whether only separators remain.
func (s *scanner) done() bool {
	s.skipSeparators()
	return s.pos >= len(s.d)
}

// command reads a command letter; the caller checks that it is one.
func (s *scanner) command() byte {
	s.skipSeparators()
	c := s.d[s.pos]
	s.pos++
	return c
}

// numberNext reports whether the next token starts a number.
func (s *scanner) numberNext() bool {
	s.skipSeparators()
	return s.pos < len(s.d) && isNumberStart(s.d[s.pos])
}

// number reads one number in the SVG grammar: an optional sign, digits with
// an optional fraction (or a fraction alone), and an optional exponent.
func (s *scanner) number() (float64, error) {
	s.skipSeparators()
	start := s.pos
	if s.pos < len(s.d) && (s.d[s.pos] == '+' || s.d[s.pos] == '-') {
		s.pos++
	}
	digits := s.digits()
	if s.pos < len(s.d) && s.d[s.pos] == '.' {
		s.pos++
		digits += s.digits()
	}
	if digits == 0 {
		return 0, s.fail(ErrBadNumber, start)
	}
	s.exponent()
	value, err := strconv.ParseFloat(s.d[start:s.pos], 64)
	if err != nil || math.Abs(value) > maxCoordinate {
		return 0, s.fail(ErrOutOfRange, start)
	}
	return value, nil
}

// digits skips a run of digits and returns how many there were.
func (s *scanner) digits() int {
	start := s.pos
	for s.pos < len(s.d) && s.d[s.pos] >= '0' && s.d[s.pos] <= '9' {
		s.pos++
	}
	return s.pos - start
}

// exponent consumes "e", an optional sign and digits, but only if digits follow.
func (s *scanner) exponent() {
	if s.pos >= len(s.d) || (s.d[s.pos] != 'e' && s.d[s.pos] != 'E') {
		return
	}
	save := s.pos
	s.pos++
	if s.pos < len(s.d) && (s.d[s.pos] == '+' || s.d[s.pos] == '-') {
		s.pos++
	}
	if s.digits() == 0 {
		s.pos = save
	}
}

// point reads an x, y pair.
func (s *scanner) point() (float64, float64, error) {
	x, err := s.number()
	if err != nil {
		return 0, 0, err
	}
	y, err := s.number()
	return x, y, err
}

func isNumberStart(c byte) bool {
	return c == '-' || c == '+' || c == '.' || (c >= '0' && c <= '9')
}
