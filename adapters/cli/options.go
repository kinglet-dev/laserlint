package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kinglet-dev/laserlint/domain/check"
	"github.com/kinglet-dev/laserlint/domain/units"
)

// options are the parsed and validated arguments.
type options struct {
	file     string
	json     bool
	version  bool
	settings check.Settings
}

// parse reads the arguments. Values are taken as text and validated here,
// so every message names the flag and says how to fix it.
func parse(args []string) (options, error) {
	fs := flag.NewFlagSet("laserlint", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	line := fs.String("line", "0.10mm", "")
	gap := fs.String("gap", "0.25mm", "")
	warnClose := fs.String("warn-close", "1%", "")
	maxClose := fs.String("max-close", "10%", "")
	var o options
	fs.BoolVar(&o.json, "json", false, "")
	fs.BoolVar(&o.version, "version", false, "")
	if err := fs.Parse(args); err != nil || o.version {
		return o, err
	}
	s := check.DefaultSettings
	var err error
	if s.Line, err = length("--line", *line); err == nil && s.Line == 0 {
		err = errors.New("--line must be more than zero")
	}
	if err == nil {
		s.Gap, err = length("--gap", *gap)
	}
	if err == nil {
		s.WarnClose, err = percent("--warn-close", *warnClose)
	}
	if err == nil {
		s.MaxClose, err = percent("--max-close", *maxClose)
	}
	if err == nil && s.WarnClose > s.MaxClose {
		err = errors.New("--warn-close must not be more than --max-close")
	}
	if err == nil && fs.NArg() != 1 {
		err = errors.New("give one SVG file to check, or - to read standard input")
	}
	o.file, o.settings = fs.Arg(0), s
	return o, err
}

// length reads a length that must have a unit, in millimetres.
func length(flag, value string) (float64, error) {
	l, err := units.ParseLength(value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", flag, err)
	}
	return l.Millimetres(), nil
}

// percent reads a percentage such as "10%" as a share from 0 to 1.
func percent(flag, value string) (float64, error) {
	number, ok := strings.CutSuffix(value, "%")
	v, err := strconv.ParseFloat(number, 64)
	if !ok || err != nil {
		return 0, fmt.Errorf("%s: write a percentage such as 10%%", flag)
	}
	if !(v >= 0 && v <= 100) {
		return 0, fmt.Errorf("%s: write a percentage from 0%% to 100%%", flag)
	}
	return v / 100, nil
}
