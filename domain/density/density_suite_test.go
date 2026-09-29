package density_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDensity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Density Suite")
}
