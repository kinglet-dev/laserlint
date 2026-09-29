package text_test

import (
	"bytes"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/laserlint/adapters/report"
	"github.com/kinglet-dev/laserlint/adapters/report/text"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

var about = report.About{Version: "0.1.0", Input: "coaster.svg"}

func result(name string, findings ...check.Finding) app.Result {
	return app.Result{ID: "id", Name: name, Findings: findings}
}

func reportOf(results ...app.Result) app.Report {
	return app.Report{Width: 80, Height: 74.5, Settings: check.DefaultSettings, Checks: results}
}

func write(r app.Report) string {
	var b bytes.Buffer
	Expect(text.New().Write(&b, about, r)).To(Succeed())
	return b.String()
}

var closeFinding = check.Finding{
	Severity:  check.Problem,
	Message:   "38.0% of scored line is too close",
	Fix:       "Space lines apart.",
	Locations: []geom.Point{{X: 12.34, Y: 5}, {X: 60, Y: 70.06}},
}

var _ = Describe("The text report", func() {
	It("shows the tool, file, size and limits, each check, and the verdict", func() {
		// Act
		out := write(reportOf(result("Lines too close", closeFinding)))

		// Assert
		Expect(out).To(Equal(`laserlint 0.1.0 · coaster.svg · 80.0 × 74.5 mm · line 0.10 mm · gap 0.25 mm

PROBLEM  Lines too close
         38.0% of scored line is too close
         Where: (12.3, 5.0), (60.0, 70.1) mm from the top-left
         Fix: Space lines apart.

Not ready to burn: 1 problem.
`))
	})

	It("shows a check with no findings as OK", func() {
		// Act
		out := write(reportOf(result("Small details")))

		// Assert
		Expect(out).To(ContainSubstring("\nOK       Small details\n"))
		Expect(out).To(HaveSuffix("\nReady to burn.\n"))
	})

	It("leaves out an empty Where line", func() {
		// Arrange
		f := check.Finding{Severity: check.Warning, Message: "m", Fix: "f"}

		// Act
		out := write(reportOf(result("Crossings", f)))

		// Assert
		Expect(out).NotTo(ContainSubstring("Where:"))
	})

	DescribeTable("sums up the verdict",
		func(severities []check.Severity, want string) {
			// Arrange
			var fs []check.Finding
			for _, s := range severities {
				fs = append(fs, check.Finding{Severity: s})
			}

			// Act
			out := write(reportOf(result("Any", fs...)))

			// Assert
			Expect(out).To(HaveSuffix("\n" + want + "\n"))
		},
		Entry("warnings allowed", []check.Severity{check.Warning}, "Ready to burn, with 1 warning."),
		Entry("warnings and info", []check.Severity{check.Warning, check.Warning, check.Info}, "Ready to burn, with 2 warnings and 1 info."),
		Entry("info only", []check.Severity{check.Info}, "Ready to burn, with 1 info."),
		Entry("problems and a warning", []check.Severity{check.Problem, check.Problem, check.Warning}, "Not ready to burn: 2 problems and 1 warning."),
		Entry("several info, which has no plural", []check.Severity{check.Info, check.Info}, "Ready to burn, with 2 info."),
		Entry("all three", []check.Severity{check.Info, check.Warning, check.Problem}, "Not ready to burn: 1 problem, 1 warning and 1 info."),
	)

	It("labels each severity with a word, never a symbol or colour alone", func() {
		// Act
		out := write(reportOf(
			result("A", check.Finding{Severity: check.Warning}),
			result("B", check.Finding{Severity: check.Info})))

		// Assert
		Expect(out).To(ContainSubstring("\nWARNING  A\n"))
		Expect(out).To(ContainSubstring("\nINFO     B\n"))
	})

	It("passes on a write failure", func() {
		// Act
		err := text.New().Write(failing{}, about, reportOf())

		// Assert
		Expect(err).To(MatchError(errBroken))
	})

	It("keeps the first write failure even when later writes succeed", func() {
		// Act
		err := text.New().Write(&failsOnce{}, about, reportOf())

		// Assert
		Expect(err).To(MatchError(errBroken))
	})
})

var errBroken = errors.New("broken pipe")

// failsOnce fails its first write only.
type failsOnce struct{ failed bool }

func (f *failsOnce) Write(b []byte) (int, error) {
	if !f.failed {
		f.failed = true
		return 0, errBroken
	}
	return len(b), nil
}

type failing struct{}

func (failing) Write([]byte) (int, error) { return 0, errBroken }
