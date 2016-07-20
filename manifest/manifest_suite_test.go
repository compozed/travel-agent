package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"testing"
)

var manifestPath string
var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestCi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manifest Suite")
}
