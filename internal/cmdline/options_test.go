package cmdline_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/pianohacker/quickprom/internal/cmdline"
)

var _ = Describe("Options", func() {
	Context("ParseOptsAndEnv", func() {
		It("can parse options and environment variables", func() {
			os.Args = []string{"quickprom", "query"}
			os.Setenv("QUICKPROM_TARGET", "target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.Target).To(Equal("target"))
			Expect(opts.Query).To(Equal("query"))
		})

		It("can override environment variables with options", func() {
			os.Args = []string{"quickprom", "-t", "cmdline_target", "query"}
			os.Setenv("QUICKPROM_TARGET", "env_target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.Target).To(Equal("cmdline_target"))
			Expect(opts.Query).To(Equal("query"))
		})

		It("can parse a timestamp when --time is given", func() {
			os.Args = []string{
				"quickprom",
				"--time",
				"2018-01-02 00:12:45.000 UTC",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "env_target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.Time).To(BeTemporally("~", time.Date(
				2018, 1, 2,
				0, 12, 45,
				0,
				time.UTC,
			)))
		})

		It("defaults to now when --time is not given", func() {
			os.Args = []string{
				"quickprom",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "env_target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.Time).To(BeTemporally("~", time.Now()))
		})

		It("can parse timestamps when `range` is given", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--end",
				"4:05 PM",
				"--step",
				"1d",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "env_target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			now := time.Now()
			fourPM := time.Date(
				now.Year(), now.Month(), now.Day(),
				16, 5, 0,
				0,
				time.Local,
			)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.RangeStart).To(BeTemporally("~", time.Date(
				2018, 1, 2,
				0, 12, 45,
				0,
				time.UTC,
			)))
			Expect(opts.RangeEnd).To(Equal(fourPM))
			Expect(opts.RangeStep).To(Equal(24 * time.Hour))
		})

		It("defaults to a range end of now", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--step",
				"1d",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "env_target")
			opts, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).ToNot(HaveOccurred())

			Expect(opts.RangeEnd).To(BeTemporally("~", time.Now()))
		})

		It("returns an error when target is unspecified", func() {
			os.Args = []string{"quickprom", "query"}
			os.Setenv("QUICKPROM_TARGET", "")
			_, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when range start is omitted", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--step",
				"1d",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "target")
			_, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when range end is before start", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--start",
				"2018-01-01",
				"--end",
				"2017-01-01",
				"--step",
				"1d",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "target")
			_, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when range step is omitted", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "target")
			_, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).To(HaveOccurred())
		})

		It("returns an error when range step is invalid", func() {
			os.Args = []string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--step",
				"potato",
				"query",
			}
			os.Setenv("QUICKPROM_TARGET", "target")
			_, err := cmdline.ParseOptsAndEnv(false)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("ParseTime", func() {
		DescribeTable("handles partial dates",
			func(s string, expected time.Time) {
				t, err := cmdline.ParseTime(s)

				Expect(err).ToNot(HaveOccurred())

				Expect(t).To(BeTemporally("~", expected))
			},

			Entry(
				"defaults to current year if no year is given",
				"April 24th",
				time.Date(
					time.Now().Year(), 4, 24,
					0, 0, 0,
					0,
					time.Local,
				)),

			Entry(
				"defaults to start of the month if no day is given",
				"May 2018",
				time.Date(
					2018, 5, 1,
					0, 0, 0,
					0,
					time.Local,
				)),

			Entry(
				"defaults to midnight if no time is given",
				"2018-01-01",
				time.Date(
					2018, 1, 1,
					0, 0, 0,
					0,
					time.Local,
				)),

			Entry(
				"defaults to current day if no date is given",
				"4:45:46",
				time.Date(
					time.Now().Year(), time.Now().Month(), time.Now().Day(),
					4, 45, 46,
					0,
					time.Local,
				)),

			Entry(
				"defaults to start of minute if no second is given",
				"2016-01-02 4:45 PM",
				time.Date(
					2016, 1, 2,
					16, 45, 0,
					0,
					time.Local,
				)),
		)

		// This is mostly a test over the fuzzytime dependency, but is here to ensure it doesn't
		// violate any of our assumptions
		DescribeTable(
			"returns an error for invalid times", func(s string) {
				_, err := cmdline.ParseTime(s)
				Expect(err).To(HaveOccurred())
			},

			Entry("only a year", "2018"),
			Entry("only a year and month", "2018-01"),
			Entry("only a month and day", "01-01"),
		)
	})
})
