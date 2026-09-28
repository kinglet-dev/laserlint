package spacing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSpacing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spacing Suite")
}
