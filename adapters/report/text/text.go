// Package text writes laserlint's report for people: plain text, with
// every severity spelled out so no meaning rests on symbols or colour.
package text

import (
	"fmt"
	"io"
	"strings"

	"github.com/kinglet-dev/laserlint/adapters/report"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/geom"
)

// labels are the severity words, padded so check names line up.
var labels = map[check.Severity]string{check.Info: "INFO    ", check.Warning: "WARNING ", check.Problem: "PROBLEM "}

// indent lines details up under the check name.
const indent = "         "

// Format is the text report format.
type Format struct{}

// New returns the text format.
func New() Format { return Format{} }

// Write writes the report.
func (Format) Write(w io.Writer, a report.About, r app.Report) error {
	p := &printer{w: w}
	p.printf("laserlint %s · %s · %.1f × %.1f mm · line %.2f mm · gap %.2f mm\n\n",
		a.Version, a.Input, r.Width, r.Height, r.Settings.Line, r.Settings.Gap)
	for _, c := range r.Checks {
		if len(c.Findings) == 0 {
			p.printf("OK       %s\n", c.Name)
		}
		for _, f := range c.Findings {
			p.finding(c.Name, f)
		}
	}
	p.printf("\n%s\n", verdict(r))
	return p.err
}

// finding writes one finding under its check's name.
func (p *printer) finding(name string, f check.Finding) {
	p.printf("%s %s\n%s%s\n", labels[f.Severity], name, indent, f.Message)
	if len(f.Locations) > 0 {
		p.printf("%sWhere: %s mm from the top-left\n", indent, points(f.Locations))
	}
	p.printf("%sFix: %s\n", indent, f.Fix)
}

func points(ps []geom.Point) string {
	s := make([]string, len(ps))
	for i, pt := range ps {
		s[i] = fmt.Sprintf("(%.1f, %.1f)", pt.X, pt.Y)
	}
	return strings.Join(s, ", ")
}

// verdict sums up the report in one sentence.
func verdict(r app.Report) string {
	var counts []string
	for _, c := range []struct {
		s    check.Severity
		word string
	}{{check.Problem, "problem"}, {check.Warning, "warning"}, {check.Info, "info"}} {
		if n := r.Count(c.s); n > 0 {
			counts = append(counts, plural(n, c.word))
		}
	}
	switch {
	case !r.Ready():
		return "Not ready to burn: " + list(counts) + "."
	case len(counts) > 0:
		return "Ready to burn, with " + list(counts) + "."
	}
	return "Ready to burn."
}

func plural(n int, word string) string {
	if n == 1 || word == "info" {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// list joins items as "a, b and c".
func list(items []string) string {
	if len(items) == 1 {
		return items[0]
	}
	return strings.Join(items[:len(items)-1], ", ") + " and " + items[len(items)-1]
}

// printer remembers the first write error, so Write can check once.
type printer struct {
	w   io.Writer
	err error
}

func (p *printer) printf(format string, args ...any) {
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.w, format, args...)
	}
}
