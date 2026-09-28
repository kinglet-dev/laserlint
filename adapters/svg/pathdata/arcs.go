package pathdata

import "github.com/kinglet-dev/laserlint/domain/geom"

func arcTo(st *state, s *scanner, rel bool) error {
	var radii [3]float64 // rx, ry, x-axis rotation
	for i := range radii {
		v, err := s.number()
		if err != nil {
			return err
		}
		radii[i] = v
	}
	large, err := s.flag()
	if err != nil {
		return err
	}
	sweep, err := s.flag()
	if err != nil {
		return err
	}
	to, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	for _, seg := range geom.ArcTo(st.current, radii[0], radii[1], radii[2], large, sweep, to) {
		st.add(seg)
	}
	st.current = to
	return nil
}

// flag reads an arc flag: exactly one character, 0 or 1, which may run
// straight into the next number ("0110 0" is 0, 1, 10, 0).
func (s *scanner) flag() (bool, error) {
	s.skipSeparators()
	if s.pos >= len(s.d) || (s.d[s.pos] != '0' && s.d[s.pos] != '1') {
		return false, s.fail(ErrBadFlag, s.pos)
	}
	s.pos++
	return s.d[s.pos-1] == '1', nil
}
