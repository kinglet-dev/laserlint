package fill_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFill(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fill Suite")
}
