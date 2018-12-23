package output

import (
	"fmt"

	"github.com/prometheus/common/model"
)

const timestampFormat = "2006-01-02 15:04:05.999 MST"
const shortTimestampFormat = "2006-01-02 15:04:05.999"

func OutputValue(value model.Value) {
	switch value.Type() {
	case model.ValVector:
		outputVector(value.(model.Vector))
	case model.ValMatrix:
		outputMatrix(value.(model.Matrix))
	}
}

func outputVector(vector model.Vector) {
	fmt.Print("Instant vector ")
	if len(vector) == 0 {
		fmt.Println("(empty result)")
		return
	}

	fmt.Printf("@ %s:\n", vector[0].Timestamp.Time().Format(timestampFormat))
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

func outputMatrix(matrix model.Matrix) {
	fmt.Print("Range vector ")
	if len(matrix) == 0 {
		fmt.Println("(empty result)")
		return
	}
	fmt.Println("")

	var commonLabels model.LabelSet
	if len(matrix) > 1 {
		commonLabels = MatrixInfo(matrix).GetCommonLabels()

		if len(commonLabels) > 0 {
			fmt.Printf("Common labels: %s\n", commonLabels)
		}
	}

	for _, series := range matrix {
		filterCommonLabels(series.Metric, commonLabels)
		fmt.Printf("%v:\n", series.Metric)

		for _, sample := range series.Values {
			fmt.Printf("    %s: %f\n", sample.Timestamp.Time().Format(shortTimestampFormat), sample.Value)
		}
	}
}

func filterCommonLabels(metric model.Metric, commonLabels model.LabelSet) {
	for labelName, _ := range commonLabels {
		delete(metric, labelName)
	}
}
