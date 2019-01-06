package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apcera/termtables"
	isatty "github.com/mattn/go-isatty"
	"github.com/prometheus/common/model"
)

const TimeFormatWithTZ = "2006-01-02 15:04:05.000 MST"
const TimeFormatWithDate = "2006-01-02 15:04:05.000"
const TimeFormat = "15:04:05.000"
const TimeFormatDateOnly = "2006-01-02"

type Renderable interface {
	RenderText()
}

func FormatValue(value model.Value) Renderable {
	switch value.Type() {
	case model.ValVector:
		return FormatInstantVector(value.(model.Vector))
	case model.ValMatrix:
		return FormatRangeVector(value.(model.Matrix))
	}

	return nil
}

type tableOutput interface {
	AddRow(...interface{}) *termtables.Row
	Render() string
}

func (f *FormattedInstantVector) RenderText() {
	fmt.Print("Instant vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println()

	fmt.Printf("  At: %s\n", f.Time.Format(TimeFormatWithTZ))

	outputCommonLabels(f.CommonLabels)

	// Value column
	header := append(f.VaryingLabels, "")

	tw := getTableWriter(header)

	for _, sample := range f.Samples {
		var row []interface{}

		for _, labelValue := range sample.LabelValues {
			row = append(row, labelValue)
		}

		row = append(row, fmt.Sprintf("%f", sample.Value))

		tw.AddRow(row...)
	}

	fmt.Print(tw.Render())
}

func (f *FormattedRangeVector) RenderText() {
	fmt.Print("Range vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println()

	outputCommonLabels(f.CommonLabels)

	var timestampFormat = TimeFormatWithDate
	minDate := f.MinTime.Format(TimeFormatDateOnly)
	maxDate := f.MaxTime.Format(TimeFormatDateOnly)

	if minDate == maxDate {
		fmt.Printf("  All on date: %s\n", minDate)
		timestampFormat = TimeFormat
	}

	fmt.Println()

	for _, series := range f.Series {
		for i, labelName := range f.VaryingLabels {
			if i != 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s %s", bold(labelName+":"), series.LabelValues[i])
		}
		fmt.Println(":")

		for _, sample := range series.Values {
			fmt.Printf("    %s: %f\n", sample.Time.Format(timestampFormat), sample.Value)
		}
	}
}

func outputCommonLabels(commonLabels map[string]string) {
	var labels []string
	for labelName, _ := range commonLabels {
		labels = append(labels, labelName)
	}
	sort.Sort(sort.StringSlice(labels))

	fmt.Print("  All have labels: ")
	for i, labelName := range labels {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s %s", bold(labelName+":"), commonLabels[labelName])
	}
	fmt.Println()
}

func getTableWriter(headers []string) tableOutput {
	var tw tableOutput
	if outputIsATty {
		tt := termtables.CreateTable()
		tt.Style.SkipBorder = true
		tt.Style.BorderX = ""
		tt.Style.BorderY = ""
		tt.Style.BorderI = ""

		var headerVals []interface{}
		for _, header := range headers {
			headerVals = append(headerVals, bold(header))
		}

		tt.AddHeaders(headerVals...)

		tw = tt
	} else {
		fmt.Println()
		fmt.Println(strings.Join(headers, "\t"))
		tw = &dumbTableWriter{}
	}

	return tw
}

func bold(s string) string {
	if outputIsATty {
		return "\x1b[1m" + s + "\x1b[0m"
	}

	return s
}

var outputIsATty = isatty.IsTerminal(os.Stdout.Fd())

type dumbTableWriter struct{}

func (d *dumbTableWriter) AddRow(cells ...interface{}) *termtables.Row {
	for i, v := range cells {
		if i != 0 {
			fmt.Print("\t")
		}
		fmt.Print(v)
	}

	fmt.Println()

	return nil
}

func (*dumbTableWriter) Render() string {
	return ""
}

type jsonValue struct {
	ResultType model.ValueType `json:"resultType"`
	Result     model.Value     `json:"result"`
}

func RenderJson(value model.Value) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(&jsonValue{
		ResultType: value.Type(),
		Result:     value,
	})
}
