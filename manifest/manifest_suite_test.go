  package main_test

  import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/onsi/gomega/gexec"
    "testing"
  )

  var manifestPath string
  var _ = BeforeSuite(func() {
    var err error

    manifestPath, err = gexec.Build("../manifest")
    Ω(err).ShouldNot(HaveOccurred())
  })

  var _ = AfterSuite(func() {
    gexec.CleanupBuildArtifacts()
  })

  func TestCi(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Ci Suite")
  }
