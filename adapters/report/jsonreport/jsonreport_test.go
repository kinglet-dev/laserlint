package jsonreport_test

import (
	"bytes"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/report"
	"github.com/kinglet-dev/laserlint/adapters/report/jsonreport"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var about = report.About{Version: "0.1.0", Input: "coaster.svg"}

func write(r app.Report) string {
	var b bytes.Buffer
	Expect(jsonreport.New().Write(&b, about, r)).To(Succeed())
	return b.String()
}

var _ = Describe("The JSON report (kinglet.report/v1)", func() {
	It("has the stable shape shared by Kinglet tools", func() {
		// Arrange
		r := app.Report{Width: 80, Height: 74.5, Settings: check.DefaultSettings, Checks: []app.Result{
			{ID: "lines-too-close", Name: "Lines too close", Findings: []check.Finding{{
				Check: "lines-too-close", Severity: check.Problem, Message: "38.0% too close", Fix: "Space lines apart.",
				Locations: []geom.Point{{X: 12.3456789, Y: 5}},
			}}},
			{ID: "crossings", Name: "Crossings"},
		}}

		// Act
		out := write(r)

		// Assert
		Expect(out).To(MatchJSON(`{
			"schema": "kinglet.report/v1",
			"tool": {"name": "laserlint", "version": "0.1.0"},
			"input": {"name": "coaster.svg", "width_mm": 80, "height_mm": 74.5},
			"settings": {"line_mm": 0.1, "gap_mm": 0.25, "warn_close": 0.01, "max_close": 0.1, "warn_density": 0.9, "max_density": 2,
				"detail_mm": 0.5, "max_details": 20, "max_crossings": 50, "max_stacked_mm": 2},
			"checks": [
				{"id": "lines-too-close", "name": "Lines too close"},
				{"id": "crossings", "name": "Crossings"}
			],
			"findings": [{
				"check": "lines-too-close", "severity": "problem",
				"message": "38.0% too close", "fix": "Space lines apart.",
				"locations": [{"x_mm": 12.346, "y_mm": 5}]
			}],
			"summary": {"ready": false, "problems": 1, "warnings": 0, "info": 0}
		}`))
	})

	It("names every severity", func() {
		// Arrange
		r := app.Report{Checks: []app.Result{{ID: "a", Findings: []check.Finding{
			{Severity: check.Info}, {Severity: check.Warning}}}}}

		// Act
		out := write(r)

		// Assert
		Expect(out).To(ContainSubstring(`"severity": "info"`))
		Expect(out).To(ContainSubstring(`"severity": "warning"`))
	})

	It("writes empty lists as [], never null", func() {
		// Act
		out := write(app.Report{Checks: []app.Result{{ID: "a", Findings: []check.Finding{{Severity: check.Warning}}}}})
		empty := write(app.Report{})

		// Assert
		Expect(out).To(ContainSubstring(`"locations": []`))
		Expect(empty).To(ContainSubstring(`"checks": []`))
		Expect(empty).To(ContainSubstring(`"findings": []`))
	})

	It("keeps text such as <, > and & readable", func() {
		// Arrange
		r := app.Report{Checks: []app.Result{{ID: "a", Findings: []check.Finding{{Message: "<use> & more"}}}}}

		// Act
		out := write(r)

		// Assert
		Expect(out).To(ContainSubstring(`"<use> & more"`))
	})

	It("passes on a write failure", func() {
		// Act
		err := jsonreport.New().Write(failing{}, about, app.Report{})

		// Assert
		Expect(err).To(MatchError(errBroken))
	})
})

var errBroken = errors.New("broken pipe")

type failing struct{}

func (failing) Write([]byte) (int, error) { return 0, errBroken }
