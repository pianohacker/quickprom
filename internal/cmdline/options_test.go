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
	DescribeTable("ParseOptsAndEnv",
		func(
			args []string,
			env map[string]string,
			testFunc func(opts *cmdline.QuickPromOptions, err error),
		) {
			os.Args = args
			os.Clearenv()
			for k, v := range env {
				os.Setenv(k, v)
			}

			opts, err := cmdline.ParseOptsAndEnv(false)

			testFunc(opts, err)
		},

		Entry("can parse options and environment variables",
			[]string{"quickprom", "query"},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Target).To(Equal("target"))
				Expect(opts.Query).To(Equal("query"))
			},
		),

		Entry("can override environment variables with options",
			[]string{"quickprom", "-t", "cmdline_target", "query"},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Target).To(Equal("cmdline_target"))
				Expect(opts.Query).To(Equal("query"))
			},
		),

		Entry("can parse --skip-tls-verify from command line",
			[]string{"quickprom", "--skip-tls-verify", "-t", "target", "query"},
			nil,

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.SkipTlsVerify).To(BeTrue())
			},
		),

		Entry("can parse --skip-tls-verify from environment variable",
			[]string{"quickprom", "-t", "target", "query"},
			map[string]string{
				"QUICKPROM_SKIP_TLS_VERIFY": "true",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.SkipTlsVerify).To(BeTrue())
			},
		),

		Entry("can parse --basic-auth from command line",
			[]string{"quickprom", "-t", "target", "--basic-auth", "username:password", "query"},
			nil,

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.BasicAuth).To(Equal("username:password"))
			},
		),

		Entry("can parse --basic-auth from environment variable",
			[]string{"quickprom", "-t", "target", "query"},
			map[string]string{
				"QUICKPROM_BASIC_AUTH": "env_username:env_password",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.BasicAuth).To(Equal("env_username:env_password"))
			},
		),

		Entry("can parse --cf-auth from command line",
			[]string{"quickprom", "-t", "target", "--cf-auth", "query"},
			nil,

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.CfAuth).To(BeTrue())
			},
		),

		Entry("can parse --cf-auth from environment variable",
			[]string{"quickprom", "-t", "target", "query"},
			map[string]string{
				"QUICKPROM_CF_AUTH": "true",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.CfAuth).To(BeTrue())
			},
		),

		Entry("can parse --json from command line",
			[]string{"quickprom", "-t", "target", "--json", "query"},
			nil,

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Json).To(BeTrue())
			},
		),

		Entry("can parse --json from environment variable",
			[]string{"quickprom", "-t", "target", "query"},
			map[string]string{
				"QUICKPROM_JSON": "true",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Json).To(BeTrue())
			},
		),

		Entry("can parse --range-table from command line",
			[]string{
				"quickprom",
				"--range-table",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.RangeTable).To(BeTrue())
			},
		),

		Entry("can parse --range-table from short option",
			[]string{
				"quickprom",
				"-b",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.RangeTable).To(BeTrue())
			},
		),

		Entry("can parse --range-table from environment variable",
			[]string{
				"quickprom",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET":      "env_target",
				"QUICKPROM_RANGE_TABLE": "true",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.RangeTable).To(BeTrue())
			},
		),

		Entry("can parse a timestamp when --time is given",
			[]string{
				"quickprom",
				"--time",
				"2018-01-02 00:12:45.000 UTC",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Time).To(BeTemporally("~", time.Date(
					2018, 1, 2,
					0, 12, 45,
					0,
					time.UTC,
				)))
			},
		),

		Entry("supports a short option for --time",
			[]string{
				"quickprom",
				"-i",
				"2018-01-02 00:12:45.000",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())
			},
		),

		Entry("defaults to now when --time is not given",
			[]string{
				"quickprom",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Time).To(BeTemporally("~", time.Now()))
			},
		),

		Entry("can parse timestamps when `range` is given",
			[]string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--end",
				"4:05 PM",
				"--step",
				"1d",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

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
			},
		),

		Entry("supports short options to `range`",
			[]string{
				"quickprom",
				"range",
				"-s",
				"2018-01-02 00:12:45.000 UTC",
				"-e",
				"4:05 PM",
				"-p",
				"1d",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())
			},
		),

		Entry("defaults to a range end of now",
			[]string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--step",
				"1d",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "env_target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.RangeEnd).To(BeTemporally("~", time.Now()))
			},
		),

		Entry("can parse --json from environment variable",
			[]string{"quickprom", "-t", "target", "query"},
			map[string]string{
				"QUICKPROM_JSON": "true",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).ToNot(HaveOccurred())

				Expect(opts.Json).To(BeTrue())
			},
		),

		Entry("returns an error when target is unspecified",
			[]string{"quickprom", "query"},
			map[string]string{
				"QUICKPROM_TARGET": "",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),

		Entry("returns an error when basic auth not in USER:PASS format",
			[]string{"quickprom", "--basic-auth", "badstuff", "query"},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),

		Entry("returns an error when range start is omitted",
			[]string{
				"quickprom",
				"range",
				"--step",
				"1d",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),

		Entry("returns an error when range end is before start",
			[]string{
				"quickprom",
				"range",
				"--start",
				"2018-01-01",
				"--end",
				"2017-01-01",
				"--step",
				"1d",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},

			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),

		Entry("returns an error when range step is omitted",
			[]string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},
			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),

		Entry("returns an error when range step is invalid",
			[]string{
				"quickprom",
				"range",
				"--start",
				"2018-01-02 00:12:45.000 UTC",
				"--step",
				"potato",
				"query",
			},
			map[string]string{
				"QUICKPROM_TARGET": "target",
			},
			func(opts *cmdline.QuickPromOptions, err error) {
				Expect(err).To(HaveOccurred())
			},
		),
	)

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
