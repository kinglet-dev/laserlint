package tooclose_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTooClose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lines Too Close Suite")
}
