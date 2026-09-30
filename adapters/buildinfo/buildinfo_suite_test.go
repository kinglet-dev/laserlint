package buildinfo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBuildInfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Info Suite")
}
