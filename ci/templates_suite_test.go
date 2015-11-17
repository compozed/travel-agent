package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"

	"os"
	"syscall"
	"testing"
)

var _ = BeforeSuite(func() {
	// build the ego -> go templates
	var waitStatus syscall.WaitStatus
	cmd := exec.Command("ego", "-package=templates")
	_, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		}
	}
})

var _ = AfterSuite(func() {
	os.Remove("ego.go")
})

func TestCi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ci Suite")
}
