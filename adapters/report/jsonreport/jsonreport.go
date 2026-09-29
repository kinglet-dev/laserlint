// Package jsonreport writes laserlint's report as kinglet.report/v1 JSON,
// the machine-readable format shared by every Kinglet tool. It is public
// surface: fields are only ever added within a major version.
package jsonreport

import (
	"encoding/json"
	"io"
	"math"

	"github.com/kinglet-dev/laserlint/adapters/report"
	"github.com/kinglet-dev/laserlint/app"
	"github.com/kinglet-dev/laserlint/domain/check"
)

// Schema names the format and its major version.
const Schema = "kinglet.report/v1"

type document struct {
	Schema   string    `json:"schema"`
	Tool     tool      `json:"tool"`
	Input    input     `json:"input"`
	Settings settings  `json:"settings"`
	Checks   []ran     `json:"checks"`
	Findings []finding `json:"findings"`
	Summary  summary   `json:"summary"`
}

type tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type input struct {
	Name   string  `json:"name"`
	Width  float64 `json:"width_mm"`
	Height float64 `json:"height_mm"`
}

type settings struct {
	Line      float64 `json:"line_mm"`
	Gap       float64 `json:"gap_mm"`
	WarnClose float64 `json:"warn_close"`
	MaxClose  float64 `json:"max_close"`
}

type ran struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type finding struct {
	Check     string     `json:"check"`
	Severity  string     `json:"severity"`
	Message   string     `json:"message"`
	Fix       string     `json:"fix"`
	Locations []location `json:"locations"`
}

type location struct {
	X float64 `json:"x_mm"`
	Y float64 `json:"y_mm"`
}

type summary struct {
	Ready    bool `json:"ready"`
	Problems int  `json:"problems"`
	Warnings int  `json:"warnings"`
	Info     int  `json:"info"`
}

var severities = map[check.Severity]string{check.Info: "info", check.Warning: "warning", check.Problem: "problem"}

// Format is the JSON report format.
type Format struct{}

// New returns the JSON format.
func New() Format { return Format{} }

// Write writes the report as indented JSON.
func (Format) Write(w io.Writer, a report.About, r app.Report) error {
	d := document{
		Schema:   Schema,
		Tool:     tool{Name: "laserlint", Version: a.Version},
		Input:    input{Name: a.Input, Width: mm(r.Width), Height: mm(r.Height)},
		Settings: settings{Line: r.Settings.Line, Gap: r.Settings.Gap, WarnClose: r.Settings.WarnClose, MaxClose: r.Settings.MaxClose},
		Checks:   []ran{},
		Findings: []finding{},
		Summary: summary{Ready: r.Ready(), Problems: r.Count(check.Problem),
			Warnings: r.Count(check.Warning), Info: r.Count(check.Info)},
	}
	for _, c := range r.Checks {
		d.Checks = append(d.Checks, ran{ID: c.ID, Name: c.Name})
		for _, f := range c.Findings {
			out := finding{Check: f.Check, Severity: severities[f.Severity], Message: f.Message, Fix: f.Fix, Locations: []location{}}
			for _, p := range f.Locations {
				out.Locations = append(out.Locations, location{X: mm(p.X), Y: mm(p.Y)})
			}
			d.Findings = append(d.Findings, out)
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(d)
}

// mm rounds to a thousandth of a millimetre, so output is the same on
// every CPU.
func mm(v float64) float64 { return math.Round(v*1000) / 1000 }
