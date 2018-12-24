package output

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/common/model"
)

const timestampFormat = "2006-01-02 15:04:05.000 MST"
const shortTimestampFormat = "2006-01-02 15:04:05.000"

func OutputValue(value model.Value) {
	switch value.Type() {
	case model.ValVector:
		outputVector(value.(model.Vector))
	case model.ValMatrix:
		outputMatrix(value.(model.Matrix))
	}
}

type tableOutput interface {
	Append([]string)
	Render()
}

func outputVector(vector model.Vector) {
	fmt.Print("Instant vector ")
	if len(vector) == 0 {
		fmt.Println("(empty result)")
		return
	}

	fmt.Printf("@ %s:\n", vector[0].Timestamp.Time().Format(timestampFormat))
	var commonLabels model.LabelSet
	info := VectorInfo(vector)

	if len(vector) > 1 {
		commonLabels = VectorInfo(vector).GetCommonLabels()

		if len(commonLabels) > 0 {
			fmt.Printf("Common labels: %s\n", commonLabels)
		}
	}

	var uncommonLabelSet []string

	for labelName, _ := range info.labelInfo {
		if _, ok := commonLabels[labelName]; ok {
			continue
		}

		uncommonLabelSet = append(uncommonLabelSet, string(labelName))
	}
	sort.Sort(sort.StringSlice(uncommonLabelSet))

	// Value column
	header := append(uncommonLabelSet, "")

	tw := getTableWriter(header)

	for _, sample := range vector {
		var row []string

		for _, uncommonLabel := range uncommonLabelSet {
			row = append(row, string(sample.Metric[model.LabelName(uncommonLabel)]))
		}

		row = append(row, fmt.Sprintf("%f", sample.Value))

		tw.Append(row)
	}

	tw.Render()
}

func outputMatrix(matrix model.Matrix) {
	fmt.Print("Range vector ")
	if len(matrix) == 0 {
		fmt.Println("(empty result)")
		return
	}

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

func getTableWriter(header []string) tableOutput {
	if outputIsATty {
		tw := tablewriter.NewWriter(os.Stdout)
		tw.SetHeader(header)
		tw.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		tw.SetHeaderLine(false)
		tw.SetAutoFormatHeaders(false)
		tw.SetBorder(false)
		tw.SetCenterSeparator("")
		tw.SetColumnSeparator("")
		tw.SetRowSeparator("")

		var headerColors []tablewriter.Colors

		for range header {
			headerColors = append(headerColors, tablewriter.Colors{tablewriter.Bold})
		}

		tw.SetHeaderColor(headerColors...)

		return tw
	} else {
		return &dumbTableWriter{
			header: header,
		}
	}
}

var outputIsATty = isatty.IsTerminal(os.Stdout.Fd())

type dumbTableWriter struct {
	header []string
}

func (d *dumbTableWriter) Append(row []string) {
	if d.header != nil {
		fmt.Println(strings.Join(d.header, "\t"))
		d.header = nil
	}

	fmt.Println(strings.Join(row, "\t"))
}

func (*dumbTableWriter) Render() {}

func filterCommonLabels(metric model.Metric, commonLabels model.LabelSet) {
	for labelName, _ := range commonLabels {
		delete(metric, labelName)
	}
}
