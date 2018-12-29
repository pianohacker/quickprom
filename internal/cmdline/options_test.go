package cmdline_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pianohacker/quickprom/internal/cmdline"
)

var _ = Describe("Options", func() {
	It("can parse options and environment variables", func() {
		os.Args = []string{"", "query"}
		os.Setenv("QUICKPROM_TARGET", "target")
		opts, err := cmdline.ParseOptsAndEnv(false)

		Expect(err).ToNot(HaveOccurred())

		Expect(opts).To(Equal(&cmdline.QuickPromOptions{
			Target: "target",
			Query: "query",
		}))
	})

	It("can override environment variables with options", func() {
		os.Args = []string{"", "-t", "cmdline_target", "query"}
		os.Setenv("QUICKPROM_TARGET", "env_target")
		opts, err := cmdline.ParseOptsAndEnv(false)

		Expect(err).ToNot(HaveOccurred())

		Expect(opts).To(Equal(&cmdline.QuickPromOptions{
			Target: "cmdline_target",
			Query: "query",
		}))
	})

	It("returns an error when target is unspecified", func() {
		os.Args = []string{"", "query"}
		os.Setenv("QUICKPROM_TARGET", "")
		_, err := cmdline.ParseOptsAndEnv(false)

		Expect(err).To(HaveOccurred())
	})
})
