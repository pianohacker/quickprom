package output_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/model"

	"github.com/pianohacker/quickprom/output"
)

var _ = Describe("Labels", func() {
	Describe("GetCommonLabels", func() {
		It("includes common labels", func() {
			Expect(output.VectorInfo(model.Vector{
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
			}).GetCommonLabels()).To(Equal(model.LabelSet{
				"shared-a": "a",
				"shared-b": "b",
			}))
		})

		It("ignores varying labels", func() {
			Expect(output.VectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"shared-a": "a",
						"varying-b": "b",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"varying-b": "bee",
					},
				},
			}).GetCommonLabels()).To(Equal(model.LabelSet{
				"shared-a": "a",
			}))
		})

		It("ignores non-shared labels", func() {
			Expect(output.VectorInfo(model.Vector{
				{
					Metric: model.Metric{
						"shared-a": "a",
					},
				},
				{
					Metric: model.Metric{
						"shared-a": "a",
						"b": "bee",
					},
				},
			}).GetCommonLabels()).To(Equal(model.LabelSet{
				"shared-a": "a",
			}))
		})
	})
})
