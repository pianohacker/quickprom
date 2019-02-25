package output

import (
	"fmt"
	"time"

	"github.com/prometheus/common/model"
)

type FormattedValue struct {
	MinValueExp        int
	MaxValueExp        int
	MaxValueFracLength int
	Empty              bool
	CommonLabels       map[string]string
	VaryingLabels      []string
}

type FormattedScalar struct {
	FormattedValue
	FormattedSamplePair
}

type FormattedInstantVector struct {
	FormattedValue
	Time    time.Time
	Samples []FormattedSample
}

type FormattedSample struct {
	LabelValues []string
	Value       float64
}

type FormattedRangeVector struct {
	FormattedValue
	MinTime   time.Time
	MaxTime   time.Time
	SeenTimes []time.Time
	Series    []FormattedSeries
}

type FormattedSeries struct {
	LabelValues []string
	Values      []FormattedSamplePair
}

type FormattedSamplePair struct {
	Time  time.Time
	Value float64
}

func FormatScalar(s *model.Scalar) *FormattedScalar {
	if s == nil {
		return &FormattedScalar{
			FormattedValue: FormattedValue{
				Empty: true,
			},
		}
	}

	result := &FormattedScalar{}

	info := ScalarInfo(s)
	result.MinValueExp = info.MinValueExp
	result.MaxValueExp = info.MaxValueExp
	result.MaxValueFracLength = info.MaxValueFracLength

	result.Time = s.Timestamp.Time()
	result.Value = float64(s.Value)

	return result
}

func FormatInstantVector(v model.Vector) *FormattedInstantVector {
	if len(v) == 0 {
		return &FormattedInstantVector{
			FormattedValue: FormattedValue{
				Empty: true,
			},
		}
	}

	result := &FormattedInstantVector{}

	info := InstantVectorInfo(v)
	result.CommonLabels = info.CommonLabels()
	result.VaryingLabels = info.VaryingLabels()
	result.MinValueExp = info.MinValueExp
	result.MaxValueExp = info.MaxValueExp
	result.MaxValueFracLength = info.MaxValueFracLength

	result.Time = v[0].Timestamp.Time()

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
		return &FormattedRangeVector{
			FormattedValue: FormattedValue{
				Empty: true,
			},
		}
	}

	result := &FormattedRangeVector{}

	info := RangeVectorInfo(m)
	result.CommonLabels = info.CommonLabels()
	result.VaryingLabels = info.VaryingLabels()
	result.MinValueExp = info.MinValueExp
	result.MaxValueExp = info.MaxValueExp
	result.MaxValueFracLength = info.MaxValueFracLength

	result.SeenTimes = info.SeenTimes()
	result.MinTime = result.SeenTimes[0]
	result.MaxTime = result.SeenTimes[len(result.SeenTimes)-1]

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

func (f *FormattedRangeVector) CollateSeriesValuesByTime() (result [][]*float64) {
	for _, series := range f.Series {
		var row []*float64

		samplePos := 0
		for _, seenTime := range f.SeenTimes {
			for samplePos < len(series.Values) && series.Values[samplePos].Time.Before(seenTime) {
				samplePos++
			}

			if samplePos < len(series.Values) && series.Values[samplePos].Time == seenTime {
				row = append(row, (*float64)(&(series.Values[samplePos].Value)))
			} else {
				row = append(row, nil)
			}
		}

		result = append(result, row)
	}

	return
}

func (f *FormattedValue) BestFloatFormat() string {
	prec := f.MaxValueFracLength
	if prec > 6 {
		prec = 6
	}

	if f.MinValueExp <= -4 || f.MaxValueExp >= 6 {
		return fmt.Sprintf("%%.%de", prec)
	}

	return fmt.Sprintf("%%.%df", prec)
}

type DateParts struct {
	Date            bool
	ZeroSecond      bool
	ZeroMillisecond bool
}

func SharedDateParts(times []time.Time) *DateParts {
	if len(times) <= 1 {
		return &DateParts{}
	}

	result := &DateParts{
		Date:            true,
		ZeroSecond:      true,
		ZeroMillisecond: true,
	}

	firstYear, firstMonth, firstDay := times[0].Date()

	for i, t := range times {
		if i != 0 && result.Date {
			year, month, day := t.Date()

			if year != firstYear || month != firstMonth || day != firstDay {
				result.Date = false
			}
		}

		if result.ZeroMillisecond {
			result.ZeroMillisecond = t.Nanosecond() == 0
		}

		if result.ZeroSecond {
			result.ZeroSecond = result.ZeroMillisecond && t.Second() == 0
		}
	}

	return result
}
