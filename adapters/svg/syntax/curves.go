package syntax

import "github.com/kinglet-dev/laserlint/domain/geom"

func cubicTo(st *state, s *scanner, rel bool) error {
	c1, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	return finishCubic(st, s, rel, c1)
}

func smoothCubicTo(st *state, s *scanner, rel bool) error {
	return finishCubic(st, s, rel, st.mirror('C'))
}

// finishCubic reads the second control and end points and adds the cubic.
func finishCubic(st *state, s *scanner, rel bool, c1 geom.Point) error {
	c2, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	to, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	st.add(geom.CubicTo(c1, c2, to))
	st.lastKind, st.lastCtrl = 'C', c2
	return nil
}

func quadTo(st *state, s *scanner, rel bool) error {
	q, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	return finishQuad(st, s, rel, q)
}

func smoothQuadTo(st *state, s *scanner, rel bool) error {
	return finishQuad(st, s, rel, st.mirror('Q'))
}

// finishQuad reads the end point and adds the quadratic with control q.
func finishQuad(st *state, s *scanner, rel bool, q geom.Point) error {
	to, err := readPoint(st, s, rel)
	if err != nil {
		return err
	}
	st.add(geom.QuadTo(st.current, q, to))
	st.lastKind, st.lastCtrl = 'Q', q
	return nil
}

// mirror returns the previous control point reflected about the pen when the
// previous command was of the given kind ('C' for C/S, 'Q' for Q/T), and the
// pen itself otherwise, as the SVG specification defines for S and T.
func (st *state) mirror(kind byte) geom.Point {
	if st.prevKind != kind {
		return st.current
	}
	return geom.Point{X: 2*st.current.X - st.prevCtrl.X, Y: 2*st.current.Y - st.prevCtrl.Y}
}
