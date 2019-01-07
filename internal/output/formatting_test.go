package output_test

import (
	"time"

	"github.com/prometheus/common/model"

	"github.com/pianohacker/quickprom/internal/output"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatting", func() {
	Describe("FormatInstantVector()", func() {
		It("can handle an empty instant vector", func() {
			formatted := output.FormatInstantVector(model.Vector{})

			Expect(formatted.Empty).To(BeTrue())
		})

		It("can handle a single-sample instant vector", func() {
			formatted := output.FormatInstantVector(model.Vector{
				{
					Timestamp: 4,
					Metric: model.Metric{
						"label": "value",
					},
					Value: 123,
				},
			})

			Expect(formatted.Empty).To(BeFalse())
			Expect(formatted.Time).To(BeTemporally("~", time.Unix(0, 4e6)))
			Expect(formatted.CommonLabels).To(Equal(map[string]string{
				"label": "value",
			}))
			Expect(formatted.VaryingLabels).To(BeEmpty())
			Expect(formatted.Samples).To(ContainElement(output.FormattedSample{
				Value:       123,
			}))
		})

		It("can handle a multi-sample instant vector", func() {
			formatted := output.FormatInstantVector(model.Vector{
				{
					Timestamp: 4,
					Metric: model.Metric{
						"varying-label-a": "varying-value-1",
						"varying-label-b": "varying-value-2",
						"shared-label":    "shared-value",
					},
					Value: 123,
				},
				{
					Timestamp: 4,
					Metric: model.Metric{
						"varying-label-c": "varying-value-4",
						"varying-label-a": "varying-value-3",
						"shared-label":    "shared-value",
					},
					Value: 321,
				},
			})

			Expect(formatted.Empty).To(BeFalse())
			Expect(formatted.Time).To(BeTemporally("~", time.Unix(0, 4e6)))
			Expect(formatted.CommonLabels).To(Equal(map[string]string{
				"shared-label": "shared-value",
			}))
			Expect(formatted.VaryingLabels).To(Equal([]string{
				"varying-label-a",
				"varying-label-b",
				"varying-label-c",
			}))
			Expect(formatted.Samples).To(Equal([]output.FormattedSample{
				{
					LabelValues: []string{
						"varying-value-1",
						"varying-value-2",
						"",
					},
					Value: 123,
				},
				{
					LabelValues: []string{
						"varying-value-3",
						"",
						"varying-value-4",
					},
					Value: 321,
				},
			}))
		})
	})

	Describe("FormatRangeVector()", func() {
		It("can handle an empty range vector", func() {
			formatted := output.FormatRangeVector(model.Matrix{})

			Expect(formatted.Empty).To(BeTrue())
		})

		It("can handle a single-series range vector", func() {
			formatted := output.FormatRangeVector(model.Matrix{
				{
					Metric: model.Metric{
						"label": "value",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 1,
							Value:     11,
						},
						{
							Timestamp: 2,
							Value:     12,
						},
					},
				},
			})

			Expect(formatted.Empty).To(BeFalse())
			Expect(formatted.CommonLabels).To(Equal(map[string]string{
				"label": "value",
			}))
			Expect(formatted.VaryingLabels).To(BeEmpty())
			Expect(formatted.MinTime).To(BeTemporally("~", time.Unix(0, 1e6)))
			Expect(formatted.MaxTime).To(BeTemporally("~", time.Unix(0, 2e6)))
			Expect(formatted.Series).To(Equal([]output.FormattedSeries{
				{
					Values: []output.FormattedSamplePair{
						{
							Time:  time.Unix(0, 1e6),
							Value: 11,
						},
						{
							Time:  time.Unix(0, 2e6),
							Value: 12,
						},
					},
				},
			}))
		})

		It("can handle a multi-series range vector", func() {
			formatted := output.FormatRangeVector(model.Matrix{
				{
					Metric: model.Metric{
						"varying-label-a": "varying-value-1",
						"varying-label-b": "varying-value-2",
						"shared-label":    "shared-value",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 1,
							Value:     11,
						},
						{
							Timestamp: 3,
							Value:     13,
						},
					},
				},
				{
					Metric: model.Metric{
						"varying-label-c": "varying-value-4",
						"varying-label-a": "varying-value-3",
						"shared-label":    "shared-value",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 2,
							Value:     12,
						},
						{
							Timestamp: 4,
							Value:     14,
						},
					},
				},
			})

			Expect(formatted.Empty).To(BeFalse())
			Expect(formatted.CommonLabels).To(Equal(map[string]string{
				"shared-label": "shared-value",
			}))
			Expect(formatted.VaryingLabels).To(Equal([]string{
				"varying-label-a",
				"varying-label-b",
				"varying-label-c",
			}))
			Expect(formatted.MinTime).To(BeTemporally("~", time.Unix(0, 1e6)))
			Expect(formatted.MaxTime).To(BeTemporally("~", time.Unix(0, 4e6)))
			Expect(formatted.Series).To(Equal([]output.FormattedSeries{
				{
					LabelValues: []string{
						"varying-value-1",
						"varying-value-2",
						"",
					},
					Values: []output.FormattedSamplePair{
						{
							Time:  time.Unix(0, 1e6),
							Value: 11,
						},
						{
							Time:  time.Unix(0, 3e6),
							Value: 13,
						},
					},
				},
				{
					LabelValues: []string{
						"varying-value-3",
						"",
						"varying-value-4",
					},
					Values: []output.FormattedSamplePair{
						{
							Time:  time.Unix(0, 2e6),
							Value: 12,
						},
						{
							Time:  time.Unix(0, 4e6),
							Value: 14,
						},
					},
				},
			}))
		})
	})

	DescribeTable("SharedDateParts",
		func(
			expected *output.DateParts,
			times ...time.Time,
		) {
			Expect(output.SharedDateParts(times)).To(Equal(expected))
		},

		Entry("no input",
			&output.DateParts{
				Date:            false,
				ZeroSecond:      false,
				ZeroMillisecond: false,
			},
		),

		Entry("one input",
			&output.DateParts{
				Date:            false,
				ZeroSecond:      false,
				ZeroMillisecond: false,
			},
			mktime(2006, 1, 2, 3, 4, 5, 6),
		),

		Entry("shared date",
			&output.DateParts{
				Date:            true,
				ZeroSecond:      false,
				ZeroMillisecond: false,
			},
			mktime(2006, 1, 2, 3, 4, 5, 6),
			mktime(2006, 1, 2, 4, 4, 5, 6),
			mktime(2006, 1, 2, 5, 4, 5, 6),
		),

		Entry("not shared date",
			&output.DateParts{
				Date:            false,
				ZeroSecond:      false,
				ZeroMillisecond: false,
			},
			mktime(2006, 1, 2, 3, 4, 5, 6),
			mktime(2006, 1, 2, 3, 4, 5, 6),
			mktime(2006, 1, 3, 3, 4, 5, 6),
		),

		Entry("all zero millisecond",
			&output.DateParts{
				Date:            true,
				ZeroSecond:      false,
				ZeroMillisecond: true,
			},
			mktime(2006, 1, 2, 3, 4, 5, 0),
			mktime(2006, 1, 2, 4, 4, 5, 0),
			mktime(2006, 1, 2, 5, 4, 5, 0),
		),

		Entry("all zero second",
			&output.DateParts{
				Date:            true,
				ZeroSecond:      true,
				ZeroMillisecond: true,
			},
			mktime(2006, 1, 2, 3, 4, 0, 0),
			mktime(2006, 1, 2, 4, 4, 0, 0),
			mktime(2006, 1, 2, 5, 4, 0, 0),
		),

		Entry("not all zero second when millisecond nonzero",
			&output.DateParts{
				Date:            true,
				ZeroSecond:      false,
				ZeroMillisecond: false,
			},
			mktime(2006, 1, 2, 3, 4, 0, 1),
			mktime(2006, 1, 2, 4, 4, 0, 1),
			mktime(2006, 1, 2, 5, 4, 0, 1),
		),
	)

	Describe("CollateSeriesValuesByTime", func() {
		It("fills gaps with nil", func() {
			rangeVector := model.Matrix{
				{
					Values: []model.SamplePair{
						{
							Timestamp: 1,
							Value:     11,
						},
						{
							Timestamp: 3,
							Value:     13,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Timestamp: 2,
							Value:     12,
						},
						{
							Timestamp: 4,
							Value:     14,
						},
					},
				},
			}
			formatted := output.FormatRangeVector(rangeVector)

			Expect(formatted.CollateSeriesValuesByTime()).To(Equal([][]*float64{
				{fptr(11), nil, fptr(13), nil},
				{nil, fptr(12), nil, fptr(14)},
			}))
		})
	})

	DescribeTable("BestFloatFormat",
		func(f *output.FormattedValue, expected string) {
			Expect(f.BestFloatFormat()).To(Equal(expected))
		},

		Entry(
			"small integers",
			&output.FormattedValue{
				MinValueExp:        0,
				MaxValueExp:        1,
				MaxValueFracLength: 0,
			},
			"%.0f",
		),

		Entry(
			"small fractions",
			&output.FormattedValue{
				MinValueExp:        0,
				MaxValueExp:        1,
				MaxValueFracLength: 3,
			},
			"%.3f",
		),

		Entry(
			"over-precise fractions",
			&output.FormattedValue{
				MinValueExp:        1,
				MaxValueExp:        1,
				MaxValueFracLength: 7,
			},
			"%.6f",
		),

		Entry(
			"tiny fractions",
			&output.FormattedValue{
				MinValueExp:        -6,
				MaxValueExp:        1,
				MaxValueFracLength: 3,
			},
			"%.3e",
		),

		Entry(
			"small fractions",
			&output.FormattedValue{
				MinValueExp:        0,
				MaxValueExp:        10,
				MaxValueFracLength: 5,
			},
			"%.5e",
		),
	)
})

func mktime(y int, mo time.Month, d int, h, m, s, ms int) time.Time {
	return time.Date(
		y, mo, d,
		h, m, s,
		ms*1e6,
		time.Local,
	)
}

func fptr(f float64) *float64 {
	return &f
}
