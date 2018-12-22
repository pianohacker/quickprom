package output

import (
	"fmt"
	"time"

	"github.com/prometheus/common/model"
)

func OutputValue(value model.Value) {
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
	var commonLabels model.LabelSet
	if len(vector) > 1 {
		commonLabels = VectorInfo(vector).GetCommonLabels()

		if len(commonLabels) > 0 {
			fmt.Printf("Common labels: %s\n", commonLabels)
		}
	}

	for _, sample := range vector {
		filterCommonLabels(sample.Metric, commonLabels)
		fmt.Printf("%v: %f\n", sample.Metric, sample.Value)
	}
}

func filterCommonLabels(metric model.Metric, commonLabels model.LabelSet) {
	for labelName, _ := range commonLabels {
		delete(metric, labelName)
	}
}
