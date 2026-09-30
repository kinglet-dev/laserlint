# Changelog

All notable changes to laserlint. Versions follow [Semantic Versioning](https://semver.org/): within a major version, flags, exit codes and the JSON report only change in backwards-compatible ways.

## 0.1.0 (2026-09-29)

The first release.

### Checks

- **Lines too close:** the share of scored line within line width + gap (0.35 mm by default) of other line. A warning from 1%, a problem above 10%.
- **Crossings:** places where score lines cross or one ends on another, and line drawn on top of other line. A warning up to 50 crossings and 2 mm of stacked line, a problem above.
- **Small details:** closed shapes smaller than 0.5 mm across, which burn as dots. A warning up to 20, a problem above.
- **Density:** the most scored line in any 10 mm × 10 mm area. A warning above 0.9 mm of line per mm², a problem above 2.
- **Background shape:** white shapes covering the design, which laser software still scores around. A warning; smaller white shapes are noted as information.
- **Line art:** designs drawn as thin black strokes, which burn as double lines. Information.

### Command line

- `laserlint FILE` checks an SVG at the size saved in the file; `laserlint -` reads standard input.
- Flags for every limit: `--line`, `--gap`, `--warn-close`, `--max-close`, `--detail`, `--max-details`, `--max-crossings`, `--max-stacked`, `--warn-density`, `--max-density`.
- `--json` writes the shared Kinglet report format, `kinglet.report/v1`.
- Exit codes: 0 ready to burn, 1 problems found, 2 the file couldn't be checked.
