package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Quickprom", func() {
	It("can build and run", func() {
		compiledPath, err := gexec.Build("main.go")
		Expect(err).ToNot(HaveOccurred())

		cmd := exec.Command(compiledPath, "--help")
		quickpromSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		quickpromSession.Wait()
		Expect(quickpromSession.ExitCode()).To(Equal(0))

		gexec.CleanupBuildArtifacts()
	})
})
