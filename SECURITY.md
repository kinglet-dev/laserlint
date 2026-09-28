# Security policy

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub:
**[Report a vulnerability](https://github.com/kinglet-dev/laserlint/security/advisories/new)**
(Security tab → "Report a vulnerability").

Please don't open a public issue for security problems.

Include what you found, how to reproduce it (a sample SVG helps), and the
impact you expect.

## What laserlint trusts

Nothing in the SVG it checks. laserlint reads one file or standard input,
writes only to the terminal, uses no network, and never follows references
inside the file. Inputs over its size and complexity limits are refused.

The machine-readable contact is published at
https://kinglet.dev/.well-known/security.txt (RFC 9116).
