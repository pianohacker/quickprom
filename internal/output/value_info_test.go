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
						"shared-a":  "a",
						"varying-c": "c",
						"varying-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a":  "a",
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

	Describe("SeenTimes", func() {
		It("supports instant vectors", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Timestamp: 4,
				},
				{
					Timestamp: 4,
				},
			})

			Expect(info.SeenTimes()).To(Equal([]time.Time{
				time.Unix(0, 4e6),
			}))
		})

		It("supports range vectors", func() {
			info := output.RangeVectorInfo(model.Matrix{
				{
					Values: []model.SamplePair{
						{
							Timestamp: 4,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Timestamp: 6,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Timestamp: 5,
						},
					},
				},
			})

			Expect(info.SeenTimes()).To(Equal([]time.Time{
				time.Unix(0, 4e6),
				time.Unix(0, 5e6),
				time.Unix(0, 6e6),
			}))
		})

		It("sorts ands ignores duplicates", func() {
			info := output.RangeVectorInfo(model.Matrix{
				{
					Values: []model.SamplePair{
						{
							Timestamp: 5,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Timestamp: 4,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Timestamp: 4,
						},
					},
				},
			})

			Expect(info.SeenTimes()).To(Equal([]time.Time{
				time.Unix(0, 4e6),
				time.Unix(0, 5e6),
			}))
		})
	})

	Describe("ValueInfo", func() {
		It("supports instant vectors", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Value: 60,
				},
				{
					Value: 0.125,
				},
			})

			Expect(info.MaxValueExp).To(Equal(1))
			Expect(info.MinValueExp).To(Equal(-1))
			Expect(info.MaxValueFracLength).To(Equal(3))
		})

		It("supports fraction length for huge values", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Value: 1.125e10,
				},
				{
					Value: 1.5e10,
				},
			})

			Expect(info.MaxValueExp).To(Equal(10))
			Expect(info.MinValueExp).To(Equal(10))
			Expect(info.MaxValueFracLength).To(Equal(3))
		})

		It("supports fraction length for tiny values", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Value: 1.125e-10,
				},
				{
					Value: 1.5e-10,
				},
			})

			Expect(info.MaxValueExp).To(Equal(-10))
			Expect(info.MinValueExp).To(Equal(-10))
			Expect(info.MaxValueFracLength).To(Equal(3))
		})

		It("ignores 0", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Value: 1.125e10,
				},
				{
					Value: 0,
				},
			})

			Expect(info.MaxValueExp).To(Equal(10))
			Expect(info.MinValueExp).To(Equal(10))
		})

		It("returns 0 if no nonzero values are seen", func() {
			info := output.InstantVectorInfo(model.Vector{
				{
					Value: 0,
				},
				{
					Value: 0,
				},
			})

			Expect(info.MaxValueExp).To(Equal(0))
			Expect(info.MinValueExp).To(Equal(0))
		})

		It("supports range vectors", func() {
			info := output.RangeVectorInfo(model.Matrix{
				{
					Values: []model.SamplePair{
						{
							Value: 46,
						},
						{
							Value: 6,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Value: 622,
						},
					},
				},
				{
					Values: []model.SamplePair{
						{
							Value: 4666,
						},
					},
				},
			})

			Expect(info.MaxValueExp).To(Equal(3))
			Expect(info.MinValueExp).To(Equal(0))
			Expect(info.MaxValueFracLength).To(Equal(0))
		})
	})
})
