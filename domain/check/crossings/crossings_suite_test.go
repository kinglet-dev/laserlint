package crossings_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrossings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crossings Check Suite")
}
