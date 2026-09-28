package pathdata

import "github.com/kinglet-dev/laserlint/domain/geom"

func moveTo(st *state, s *scanner, rel bool) error {
	p, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	st.current, st.start = p, p
	st.path.Subpaths = append(st.path.Subpaths, geom.Subpath{Start: p})
	return nil
}

func lineTo(st *state, s *scanner, rel bool) error {
	p, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	st.add(geom.LineTo(p))
	return nil
}

func horizontal(st *state, s *scanner, rel bool) error {
	x, err := s.number()
	if err != nil {
		return err
	}
	if !rel {
		x -= st.current.X
	}
	p, err := st.resolve(s, x, 0, true)
	if err != nil {
		return err
	}
	st.add(geom.LineTo(p))
	return nil
}

func vertical(st *state, s *scanner, rel bool) error {
	y, err := s.number()
	if err != nil {
		return err
	}
	if !rel {
		y -= st.current.Y
	}
	p, err := st.resolve(s, 0, y, true)
	if err != nil {
		return err
	}
	st.add(geom.LineTo(p))
	return nil
}

func closePath(st *state, _ *scanner, _ bool) error {
	st.path.Subpaths[len(st.path.Subpaths)-1].Closed = true
	st.current = st.start
	return nil
}

// readPoint reads an x, y pair and resolves it against the pen.
func readPoint(st *state, s *scanner, rel bool) (geom.Point, error) {
	x, y, err := s.point()
	if err != nil {
		return geom.Point{}, err
	}
	return st.resolve(s, x, y, rel)
}
