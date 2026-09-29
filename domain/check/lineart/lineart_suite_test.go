package lineart_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLineArt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Line Art Check Suite")
}
