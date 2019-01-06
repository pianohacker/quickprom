package output_test

import (
	"time"

	"github.com/prometheus/common/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pianohacker/quickprom/internal/output"
)

var _ = Describe("Value Info", func() {
	Describe("CommonLabels", func() {
		It("includes common labels", func() {
			Expect(output.InstantVectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
			}).CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
				"shared-b": "b",
			}))
		})

		It("ignores varying labels", func() {
			Expect(output.InstantVectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"shared-a":  "a",
						"varying-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a":  "a",
						"varying-b": "bee",
					},
				},
			}).CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
			}))
		})

		It("ignores non-shared labels", func() {
			Expect(output.InstantVectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"shared-a": "a",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"b":        "bee",
					},
				},
			}).CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
			}))
		})

		It("supports range vectors", func() {
			Expect(output.RangeVectorInfo(model.Matrix{
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
			}).CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
				"shared-b": "b",
			}))
		})
	})

	Describe("TimeRange", func() {
		It("supports instant vectors", func() {
			min, max := output.InstantVectorInfo(model.Vector{
				{
					Timestamp: 4,
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
				{
					Timestamp: 4,
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
				},
			}).TimeRange()

			Expect(min).To(Equal(time.Unix(0, 4e6)))
			Expect(max).To(Equal(time.Unix(0, 4e6)))
		})

		It("supports range vectors", func() {
			min, max := output.RangeVectorInfo(model.Matrix{
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 4,
						},
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 6,
						},
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"shared-b": "b",
					},
					Values: []model.SamplePair{
						{
							Timestamp: 5,
						},
					},
				},
			}).TimeRange()

			Expect(min).To(Equal(time.Unix(0, 4e6)))
			Expect(max).To(Equal(time.Unix(0, 6e6)))
		})
	})
})
