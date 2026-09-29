package background_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBackground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Background Check Suite")
}
