package pathdata_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPathData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Path Data Suite")
}
