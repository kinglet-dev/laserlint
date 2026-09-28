package reader_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// only reads a document in mm and returns its single shape.
func only(body string) design.Shape {
	d, err := read(svg(mm96, body))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, d.Shapes).To(HaveLen(1))
	return d.Shapes[0]
}

// ends lists where each segment of a subpath ends, rounded to 1e-9.
func ends(s geom.Subpath) []geom.Point {
	var out []geom.Point
	for _, seg := range s.Segments {
		out = append(out, pt(math.Round(seg.To.X*1e9)/1e9, math.Round(seg.To.Y*1e9)/1e9))
	}
	return out
}

var _ = Describe("Read turns basic shapes into closed outlines", func() {
	It("reads a rect as four sides", func() {
		// Act
		s := only(`<rect x="1" y="2" width="10" height="5"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(1, 2)))
		Expect(ends(s)).To(Equal([]geom.Point{pt(11, 2), pt(11, 7), pt(1, 7)}))
		Expect(s.Closed).To(BeTrue())
	})

	It("rounds a rect's corners, using rx for ry when only rx is given", func() {
		// Act
		s := only(`<rect x="1" y="2" width="10" height="5" rx="2"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(3, 2)))
		Expect(ends(s)).To(Equal([]geom.Point{
			pt(9, 2), pt(11, 4), pt(11, 5), pt(9, 7), pt(3, 7), pt(1, 5), pt(1, 4), pt(3, 2)}))
		Expect(s.Closed).To(BeTrue())
	})

	It("uses ry for rx when only ry is given", func() {
		// Act
		s := only(`<rect width="10" height="10" ry="1"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(1, 0)))
	})

	It("limits corner radii to half the rect and leaves out zero-length sides", func() {
		// Act
		s := only(`<rect x="1" y="2" width="10" height="5" rx="20"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(6, 2)))
		Expect(ends(s)).To(Equal([]geom.Point{pt(11, 4.5), pt(6, 7), pt(1, 4.5), pt(6, 2)}))
	})

	It("ignores a negative corner radius and uses the other one", func() {
		// Act
		s := only(`<rect width="10" height="10" rx="-1" ry="2"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(2, 0)))
	})

	It("reads a circle as four quarter arcs", func() {
		// Act
		s := only(`<circle cx="5" cy="5" r="2"/>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(7, 5)))
		Expect(ends(s)).To(Equal([]geom.Point{pt(5, 7), pt(3, 5), pt(5, 3), pt(7, 5)}))
		for _, p := range (geom.Path{Subpaths: []geom.Subpath{s}}).Flatten(0.001)[0].Points {
			Expect(math.Hypot(p.X-5, p.Y-5)).To(BeNumerically("~", 2, 0.002))
		}
	})

	It("reads an ellipse, using rx for a missing ry", func() {
		// Act
		a := only(`<ellipse cx="5" cy="5" rx="3" ry="1"/>`).Path.Subpaths[0]
		b := only(`<ellipse cx="5" cy="5" rx="3"/>`).Path.Subpaths[0]

		// Assert
		Expect(ends(a)).To(Equal([]geom.Point{pt(5, 6), pt(2, 5), pt(5, 4), pt(8, 5)}))
		Expect(ends(b)).To(Equal([]geom.Point{pt(5, 8), pt(2, 5), pt(5, 2), pt(8, 5)}))
	})

	It("reads a polygon as a closed outline", func() {
		// Act
		sh := only(`<polygon points="0,0 10,0 10,10"/>`)

		// Assert
		Expect(sh.Filled).To(BeTrue())
		Expect(sh.Path.Subpaths[0].Closed).To(BeTrue())
		Expect(ends(sh.Path.Subpaths[0])).To(Equal([]geom.Point{pt(10, 0), pt(10, 10)}))
	})

	It("reads a polyline as an open line", func() {
		// Act
		sh := only(`<polyline points="0,0 10,0 10,10" fill="none" stroke="#000"/>`)

		// Assert
		Expect(sh.Path.Subpaths[0].Closed).To(BeFalse())
		Expect(ends(sh.Path.Subpaths[0])).To(Equal([]geom.Point{pt(10, 0), pt(10, 10)}))
	})

	It("reads a stroked line as an unfilled hairline", func() {
		// Act
		sh := only(`<line x1="1" y1="2" x2="3" y2="4" stroke="#000"/>`)

		// Assert
		Expect(sh.Filled).To(BeFalse(), "a line has no area to fill")
		Expect(sh.Path.Subpaths[0].Start).To(Equal(pt(1, 2)))
		Expect(ends(sh.Path.Subpaths[0])).To(Equal([]geom.Point{pt(3, 4)}))
	})

	It("applies transforms to shapes like any other element", func() {
		// Act
		s := only(`<g transform="translate(10 0)"><circle r="1"/></g>`).Path.Subpaths[0]

		// Assert
		Expect(s.Start).To(Equal(pt(11, 0)))
	})
})

var _ = DescribeTable("Read leaves out shapes that SVG doesn't draw",
	func(body string) {
		// Act
		d, err := read(svg(mm96, body))

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(d.Shapes).To(BeEmpty())
	},
	Entry("a rect with no width", `<rect height="5"/>`),
	Entry("a rect with zero height", `<rect width="5" height="0"/>`),
	Entry("a rect with a negative height", `<rect width="5" height="-5"/>`),
	Entry("a circle with no radius", `<circle cx="5" cy="5"/>`),
	Entry("an ellipse with no radii", `<ellipse cx="5" cy="5"/>`),
	Entry("a line without a stroke", `<line x2="5" y2="5"/>`),
	Entry("a polygon with no points", `<polygon points=""/>`),
)

var _ = DescribeTable("Read reports bad shape attributes with the element",
	func(body, want string) {
		// Act
		_, err := read(svg(mm96, body))

		// Assert
		Expect(err).To(MatchError(ContainSubstring(want)))
	},
	Entry("a unit in a width", `<rect id="r1" width="10mm" height="5"/>`, `<rect id="r1">: width: units aren't supported`),
	Entry("a bad radius", `<circle r="x"/>`, `<circle>: r: expected a number`),
	Entry("a bad corner radius", `<rect width="5" height="5" rx="1px"/>`, `<rect>: rx: units aren't supported`),
	Entry("a bad ellipse centre", `<ellipse cx="," rx="1"/>`, `<ellipse>: cx: expected a number`),
	Entry("a bad line end", `<line x2="1 2" stroke="#000"/>`, `<line>: x2: expected a number`),
	Entry("odd polygon points", `<polygon points="0 0 1"/>`, `<polygon>: points: points must come in x,y pairs`),
)
