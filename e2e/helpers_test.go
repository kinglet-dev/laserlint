package e2e_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type buffer = bytes.Buffer

func stringsReader(s string) io.Reader { return strings.NewReader(s) }

// svg is a 20 × 20 mm document where one user unit is one millimetre.
func svg(body string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="20mm" height="20mm" viewBox="0 0 20 20">` + body + `</svg>`
}

// file writes content to a temporary file and returns its path.
func file(name, content string) string {
	path := filepath.Join(GinkgoT().TempDir(), name)
	Expect(os.WriteFile(path, []byte(content), 0o600)).To(Succeed())
	return path
}

// clean has one hairline square, far from any other line.
var clean = svg(`<rect x="2" y="2" width="16" height="16" fill="none" stroke="#000"/>`)

// tooClose has two 10 mm hairlines 0.2 mm apart.
var tooClose = svg(`<g fill="none" stroke="#000"><path d="M2 5 H12"/><path d="M2 5.2 H12"/></g>`)
