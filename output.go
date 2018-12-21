package main

import (
	"fmt"
	"time"

	"github.com/prometheus/common/model"
)

func outputValue(value model.Value) {
	switch value.Type() {
	case model.ValVector:
		outputVector(value.(model.Vector))
	}
}

func outputVector(vector model.Vector) {
	fmt.Print("Instant vector ")
	if len(vector) == 0 {
		fmt.Println("(empty result)")
		return
	}

	fmt.Printf("@ %s:\n", vector[0].Timestamp.Time().Format(time.RFC1123))
	var commonTags model.LabelSet
	if len(vector) > 1 {
		commonTags = getCommonTags(vector)

		if len(commonTags) > 0 {
			fmt.Printf("Common labels: %s\n", commonTags)
		}
	}

	for _, sample := range vector {
		filterCommonTags(sample.Metric, commonTags)
		fmt.Printf("%v: %f\n", sample.Metric, sample.Value)
	}
}

func getCommonTags(vector model.Vector) (unvaryingTags model.LabelSet) {
	unvaryingTags = make(model.LabelSet)
	allLabelInfo := getLabelInfo(vector)

	for labelName, info := range allLabelInfo {
		if len(info.valueSet) == 1 && info.occurrences == len(vector) {
			for labelValue, _ := range info.valueSet {
				unvaryingTags[labelName] = labelValue
			}
		}
	}

	return
}

type labelInfoMap map[model.LabelName]*labelInfo
type labelInfo struct {
	valueSet    map[model.LabelValue]struct{}
	occurrences int
}

func getLabelInfo(vector model.Vector) labelInfoMap {
	allLabelInfo := make(labelInfoMap)

	for _, sample := range vector {
		for labelName, labelValue := range sample.Metric {
			l, existed := allLabelInfo[labelName]

			if existed {
				l.valueSet[labelValue] = struct{}{}
				l.occurrences++
			} else {
				allLabelInfo[labelName] = &labelInfo{
					occurrences: 1,
					valueSet: map[model.LabelValue]struct{}{
						labelValue: struct{}{},
					},
				}
			}
		}
	}

	return allLabelInfo
}

func filterCommonTags(metric model.Metric, commonTags model.LabelSet) {
	for labelName, _ := range commonTags {
		delete(metric, labelName)
	}
}
