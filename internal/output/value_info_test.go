package output_test

import (
	"time"

	"github.com/prometheus/common/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pianohacker/quickprom/internal/output"
)

var _ = Describe("Value Info", func() {
	Describe("CommonLabels/VaryingLabels", func() {
		It("does not consider labels common for a 1-element value", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"b": "b",
						"a": "a",
					},
				},
			})

			Expect(info.CommonLabels()).To(BeEmpty())
			Expect(info.VaryingLabels()).To(Equal([]string{"a", "b"}))
		})

		It("includes common labels", func() {
			info := output.InstantVectorInfo(model.Vector{
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
			})

			Expect(info.CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
				"shared-b": "b",
			}))
			Expect(info.VaryingLabels()).To(BeEmpty())
		})

		It("ignores varying labels", func() {
			info := output.InstantVectorInfo(model.Vector{
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
			})

			Expect(info.CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
			}))
			Expect(info.VaryingLabels()).To(Equal([]string{"varying-b"}))
		})

		It("ignores non-shared labels", func() {
			info := output.InstantVectorInfo(model.Vector{
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
			})

			Expect(info.CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
			}))
			Expect(info.VaryingLabels()).To(Equal([]string{"b"}))
		})

		It("supports range vectors", func() {
			info := output.RangeVectorInfo(model.Matrix{
				{
					Metric: model.Metric{
						"shared-a": "a",
						"varying-c": "c",
						"varying-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"varying-c": "cee",
						"varying-b": "bee",
					},
				},
			})

			Expect(info.CommonLabels()).To(Equal(map[string]string{
				"shared-a": "a",
			}))
			Expect(info.VaryingLabels()).To(Equal([]string{
				"varying-b",
				"varying-c",
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
