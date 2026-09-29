package details_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDetails(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Small Details Check Suite")
}
