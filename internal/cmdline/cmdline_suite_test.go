package cmdline_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmdline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmdline Suite")
}
