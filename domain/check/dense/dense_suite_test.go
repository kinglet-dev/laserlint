package dense_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDense(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Density Check Suite")
}
