package output_test

import (
	"time"

	"github.com/prometheus/common/model"

	"github.com/pianohacker/quickprom/internal/output"

	. "github.com/onsi/ginkgo"
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
			Expect(formatted.CommonLabels).To(BeEmpty())
			Expect(formatted.VaryingLabels).To(ConsistOf("label"))
			Expect(formatted.Samples).To(ContainElement(output.FormattedSample{
				LabelValues: []string{"value"},
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
			Expect(formatted.CommonLabels).To(BeEmpty())
			Expect(formatted.VaryingLabels).To(ConsistOf("label"))
			Expect(formatted.MinTime).To(BeTemporally("~", time.Unix(0, 1e6)))
			Expect(formatted.MaxTime).To(BeTemporally("~", time.Unix(0, 2e6)))
			Expect(formatted.Series).To(Equal([]output.FormattedSeries{
				{
					LabelValues: []string{"value"},
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
})
