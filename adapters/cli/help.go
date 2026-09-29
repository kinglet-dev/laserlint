package cli

// help is laserlint's usage text.
const help = `laserlint checks an SVG for laser score lines that are too close together.

Usage:
  laserlint [flags] FILE
  laserlint [flags] -          read the SVG from standard input

The SVG is checked at the size saved in the file, so scale it to the size
you will burn before checking.

Flags:
  --line LENGTH        width of a scored line (default 0.10mm)
  --gap LENGTH         clear gap needed between the edges of two lines (default 0.25mm)
  --warn-close PCT     warn when this share of line is too close (default 1%)
  --max-close PCT      a problem when more than this share is too close (default 10%)
  --json               write the report as JSON (kinglet.report/v1)
  --version            show the version
  -h, --help           show this help

Lengths need a unit: mm, cm, in, pt or px.

Exit status: 0 ready to burn (warnings allowed), 1 problems found,
2 the file couldn't be checked.

Examples:
  laserlint coaster.svg
  laserlint --line 0.08mm --gap 0.3mm coaster.svg
  laserlint --json coaster.svg > report.json
  cat coaster.svg | laserlint -
`
