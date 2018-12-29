package output

import (
	"time"

	"github.com/prometheus/common/model"
)

type ValueInfo struct {
	labelInfo    labelInfoMap
	length       int
	minTimestamp model.Time
	maxTimestamp model.Time
}

type labelInfoMap map[model.LabelName]*labelInfo
type labelInfo struct {
	valueSet    map[model.LabelValue]struct{}
	occurrences int
}

func VectorInfo(vector model.Vector) *ValueInfo {
	v := &ValueInfo{
		labelInfo: make(labelInfoMap),
	}

	for _, sample := range vector {
		v.addMetric(sample.Metric)
	}
	v.length = len(vector)

	if v.length > 0 {
		v.minTimestamp = vector[0].Timestamp
		v.maxTimestamp = vector[0].Timestamp
	}

	return v
}

func MatrixInfo(matrix model.Matrix) *ValueInfo {
	v := &ValueInfo{
		labelInfo: make(labelInfoMap),
		minTimestamp: model.Latest,
		maxTimestamp: model.Earliest,
	}

	for _, series := range matrix {
		v.addMetric(series.Metric)
		for _, sample := range series.Values {
			v.addTimestamp(sample.Timestamp)
		}
	}
	v.length = len(matrix)

	return v
}

func (v *ValueInfo) addMetric(metric model.Metric) {
	for labelName, labelValue := range metric {
		li, existed := v.labelInfo[labelName]

		if existed {
			li.valueSet[labelValue] = struct{}{}
			li.occurrences++
		} else {
			v.labelInfo[labelName] = &labelInfo{
				occurrences: 1,
				valueSet: map[model.LabelValue]struct{}{
					labelValue: struct{}{},
				},
			}
		}
	}
}

func (v *ValueInfo) addTimestamp(timestamp model.Time) {
	if timestamp > v.maxTimestamp {
		v.maxTimestamp = timestamp
	}

	if timestamp < v.minTimestamp {
		v.minTimestamp = timestamp
	}
}

func (v *ValueInfo) GetCommonLabels() (unvaryingTags model.LabelSet) {
	unvaryingTags = make(model.LabelSet)

	for labelName, info := range v.labelInfo {
		if len(info.valueSet) == 1 && info.occurrences == v.length {
			for labelValue, _ := range info.valueSet {
				unvaryingTags[labelName] = labelValue
			}
		}
	}

	return
}

func (v *ValueInfo) GetTimeRange() (time.Time, time.Time) {
	return v.minTimestamp.Time(), v.maxTimestamp.Time()
}
