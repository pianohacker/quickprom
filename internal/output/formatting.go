package output

import (
	"time"

	"github.com/prometheus/common/model"
)

type FormattedInstantVector struct {
	Empty         bool
	Time          time.Time
	CommonLabels  map[string]string
	VaryingLabels []string
	Samples       []FormattedSample
}

type FormattedSample struct {
	LabelValues []string
	Value       float64
}

type FormattedRangeVector struct {
	Empty         bool
	MinTime       time.Time
	MaxTime       time.Time
	CommonLabels  map[string]string
	VaryingLabels []string
	Series        []FormattedSeries
}

type FormattedSeries struct {
	LabelValues []string
	Values      []FormattedSamplePair
}

type FormattedSamplePair struct {
	Time  time.Time
	Value float64
}

func FormatInstantVector(v model.Vector) *FormattedInstantVector {
	if len(v) == 0 {
		return &FormattedInstantVector{Empty: true}
	}

	result := &FormattedInstantVector{}

	result.Time = v[0].Timestamp.Time()

	info := InstantVectorInfo(v)
	result.CommonLabels = info.CommonLabels()
	result.VaryingLabels = info.VaryingLabels()

	for _, s := range v {
		var labelValues []string
		for _, varyingLabelName := range result.VaryingLabels {
			labelValues = append(labelValues, string(s.Metric[model.LabelName(varyingLabelName)]))
		}

		result.Samples = append(result.Samples, FormattedSample{
			LabelValues: labelValues,
			Value:       float64(s.Value),
		})
	}

	return result
}

func FormatRangeVector(m model.Matrix) *FormattedRangeVector {
	if len(m) == 0 {
		return &FormattedRangeVector{Empty: true}
	}

	result := &FormattedRangeVector{}

	info := RangeVectorInfo(m)
	result.CommonLabels = info.CommonLabels()
	result.VaryingLabels = info.VaryingLabels()
	result.MinTime, result.MaxTime = info.TimeRange()

	for _, s := range m {
		var values []FormattedSamplePair

		for _, p := range s.Values {
			values = append(values, FormattedSamplePair{
				Time:  p.Timestamp.Time(),
				Value: float64(p.Value),
			})
		}

		result.Series = append(result.Series, FormattedSeries{
			LabelValues: getLabelValues(result.VaryingLabels, s.Metric),
			Values:      values,
		})
	}

	return result
}

func getLabelValues(labelNames []string, metric model.Metric) (labelValues []string) {
	for _, labelName := range labelNames {
		labelValues = append(labelValues, string(metric[model.LabelName(labelName)]))
	}

	return
}
