package app_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/design"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// fake is a check (a test double at the check port) that records its input
// and returns canned findings.
type fake struct {
	id       string
	findings []check.Finding
	err      error
	got      *check.Input
}

func (f fake) ID() string   { return f.id }
func (f fake) Name() string { return "Fake " + f.id }
func (f fake) Run(in check.Input, _ check.Settings) ([]check.Finding, error) {
	if f.got != nil {
		*f.got = in
	}
	return f.findings, f.err
}

func found(s check.Severity) check.Finding { return check.Finding{Severity: s} }

// square is a design with one 10 mm hairline square.
var square = design.Design{Width: 20, Height: 20, Shapes: []design.Shape{{Path: geom.Path{Subpaths: []geom.Subpath{{
	Start: geom.Point{}, Closed: true,
	Segments: []geom.Segment{geom.LineTo(geom.Point{X: 10}), geom.LineTo(geom.Point{X: 10, Y: 10}), geom.LineTo(geom.Point{Y: 10})},
}}}}}}

var _ = Describe("Checker runs every check on a design", func() {
	It("reports each check with its name and findings, in order", func() {
		// Arrange
		c := app.NewChecker(fake{id: "a"}, fake{id: "b", findings: []check.Finding{found(check.Warning)}})

		// Act
		r, err := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Checks).To(HaveLen(2))
		Expect(r.Checks[0].ID).To(Equal("a"))
		Expect(r.Checks[1].Name).To(Equal("Fake b"))
		Expect(r.Checks[1].Findings).To(HaveLen(1))
	})

	It("gives every check the design and its score lines", func() {
		// Arrange
		var got check.Input
		c := app.NewChecker(fake{id: "a", got: &got})

		// Act
		_, err := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Design.Width).To(Equal(20.0))
		Expect(got.Lines).To(HaveLen(1))
		Expect(got.Lines[0].Closed).To(BeTrue())
	})

	It("keeps the design's size and the settings used", func() {
		// Act
		r, _ := app.NewChecker().Check(square, check.DefaultSettings)

		// Assert
		Expect(r.Width).To(Equal(20.0))
		Expect(r.Height).To(Equal(20.0))
		Expect(r.Settings).To(Equal(check.DefaultSettings))
	})

	It("stops at a check that can't run, naming it", func() {
		// Arrange
		broken := errors.New("too big")
		c := app.NewChecker(fake{id: "a", err: broken})

		// Act
		_, err := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(err).To(MatchError(broken))
		Expect(err).To(MatchError(ContainSubstring("Fake a")))
	})
})

var _ = Describe("A report's verdict", func() {
	It("counts findings by severity across checks", func() {
		// Arrange
		c := app.NewChecker(
			fake{id: "a", findings: []check.Finding{found(check.Problem), found(check.Info)}},
			fake{id: "b", findings: []check.Finding{found(check.Problem), found(check.Warning)}})

		// Act
		r, _ := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(r.Count(check.Problem)).To(Equal(2))
		Expect(r.Count(check.Warning)).To(Equal(1))
		Expect(r.Count(check.Info)).To(Equal(1))
	})

	It("is ready to burn with warnings but no problems", func() {
		// Arrange
		c := app.NewChecker(fake{id: "a", findings: []check.Finding{found(check.Warning), found(check.Info)}})

		// Act
		r, _ := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(r.Ready()).To(BeTrue())
	})

	It("is not ready to burn with any problem", func() {
		// Arrange
		c := app.NewChecker(fake{id: "a", findings: []check.Finding{found(check.Problem)}})

		// Act
		r, _ := c.Check(square, check.DefaultSettings)

		// Assert
		Expect(r.Ready()).To(BeFalse())
	})
})
