package syntax_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSyntax(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SVG Syntax Suite")
}
