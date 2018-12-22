package output

import "github.com/prometheus/common/model"

type ValueInfo struct {
	labelInfo labelInfoMap
	length int
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

	return v
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
