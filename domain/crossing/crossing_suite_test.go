package crossing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrossing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crossing Suite")
}
