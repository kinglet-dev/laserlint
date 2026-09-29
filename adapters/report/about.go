// Package report holds what every report format shows besides the checks.
package report

// About says which tool made a report and what it checked.
type About struct {
	Version string // laserlint's version
	Input   string // the file checked, or "standard input"
}
