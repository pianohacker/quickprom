package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/xlab/termtables"
	isatty "github.com/mattn/go-isatty"
	"github.com/prometheus/common/model"
)

const TimeFormatWithTZ = "2006-01-02 15:04:05.000 MST"
const TimeFormatWithDate = "2006-01-02 15:04:05.000"
const TimeFormatDateOnly = "2006-01-02"
const TimeFormat = "15:04:05.000"
const TimeFormatZeroSecond = "15:04"
const TimeFormatZeroMillisecond = "15:04:05"

func getTimestampFormat(sharedDateParts *DateParts) string {
	var timestampFormat string

	if !sharedDateParts.Date {
		timestampFormat = TimeFormatDateOnly + " "
	}

	if sharedDateParts.ZeroSecond {
		timestampFormat += TimeFormatZeroSecond
	} else if sharedDateParts.ZeroMillisecond {
		timestampFormat += TimeFormatZeroMillisecond
	} else {
		timestampFormat += TimeFormat
	}

	return timestampFormat
}

type Renderable interface {
	RenderText(opts *RenderOptions)
}

type RenderOptions struct {
	RangeVectorAsTable bool
}

func FormatValue(value model.Value) Renderable {
	switch value.Type() {
	case model.ValScalar:
		return FormatScalar(value.(*model.Scalar))
	case model.ValVector:
		return FormatInstantVector(value.(model.Vector))
	case model.ValMatrix:
		return FormatRangeVector(value.(model.Matrix))
	}

	return nil
}

func (f *FormattedScalar) RenderText(_ *RenderOptions) {
	fmt.Print("Scalar:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println()

	fmt.Printf("  At: %s\n", f.Time.Format(TimeFormatWithTZ))

	tw := getTableWriter([]interface{}{bold("value")})
	tw.AddRow(fmt.Sprintf("%g", f.Value))
	fmt.Print(tw.Render())
}

func (f *FormattedInstantVector) RenderText(_ *RenderOptions) {
	fmt.Print("Instant vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println()

	fmt.Printf("  At: %s\n", f.Time.Format(TimeFormatWithTZ))

	outputCommonLabels("samples", f.CommonLabels)

	// Value column
	var header []interface{}

	for _, labelName := range f.VaryingLabels {
		header = append(header, bold(labelName))
	}

	header = append(header, bold("value"))

	tw := getTableWriter(header)
	floatFormat := f.BestFloatFormat()

	for _, sample := range f.Samples {
		var row []interface{}

		for _, labelValue := range sample.LabelValues {
			row = append(row, labelValue)
		}

		row = append(row, rightAlignedCell(
			fmt.Sprintf(floatFormat, sample.Value),
		))

		tw.AddRow(row...)
	}

	fmt.Print(tw.Render())
}

func (f *FormattedRangeVector) RenderText(opts *RenderOptions) {
	fmt.Print("Range vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println()

	sharedDateParts := SharedDateParts(f.SeenTimes)

	if sharedDateParts.Date {
		fmt.Printf("  All on date: %s\n", f.SeenTimes[0].Format(TimeFormatDateOnly))
	}

	if sharedDateParts.ZeroSecond {
		fmt.Println("  All timestamps end with: 00.000")
	} else if sharedDateParts.ZeroMillisecond {
		fmt.Println("  All timestamps end with: .000")
	}

	outputCommonLabels("series", f.CommonLabels)

	timestampFormat := getTimestampFormat(sharedDateParts)

	if opts.RangeVectorAsTable {
		f.renderRangeTable(timestampFormat)
	} else {
		f.renderRangeList(timestampFormat)
	}
}

func (f *FormattedRangeVector) renderRangeTable(timestampFormat string) {
	var header []interface{}

	for _, labelName := range f.VaryingLabels {
		header = append(header, bold(labelName))
	}

	for _, seenTime := range f.SeenTimes {
		header = append(header, rightAlignedCell(
			bold(seenTime.Format(timestampFormat)),
		))
	}

	tw := getTableWriter(header)

	collatedValues := f.CollateSeriesValuesByTime()
	floatFormat := f.BestFloatFormat()

	for i, series := range f.Series {
		var row []interface{}

		for _, labelValue := range series.LabelValues {
			row = append(row, labelValue)
		}

		for _, value := range collatedValues[i] {
			if value == nil {
				row = append(row, "")
			} else {
				row = append(row, rightAlignedCell(
					fmt.Sprintf(floatFormat, *value),
				))
			}
		}

		tw.AddRow(row...)
	}

	fmt.Print(tw.Render())
}

func (f *FormattedRangeVector) renderRangeList(timestampFormat string) {
	fmt.Println()
	floatFormat := f.BestFloatFormat()

	for _, series := range f.Series {
		for i, labelName := range f.VaryingLabels {
			if i != 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s %s", bold(labelName+":"), series.LabelValues[i])
		}
		fmt.Println(":")

		for _, sample := range series.Values {
			fmt.Printf("    %s: ", sample.Time.Format(timestampFormat))
			fmt.Printf(floatFormat+"\n", sample.Value)
		}
	}
}

func outputCommonLabels(subValueType string, commonLabels map[string]string) {
	if len(commonLabels) == 0 {
		return
	}

	var labels []string
	for labelName, _ := range commonLabels {
		labels = append(labels, labelName)
	}
	sort.Sort(sort.StringSlice(labels))

	fmt.Printf("  All %s are labeled: \n", subValueType)
	for _, labelName := range labels {
		fmt.Printf("    %s %s\n", bold(labelName+":"), commonLabels[labelName])
	}
}

func getTableWriter(headers []interface{}) *termtables.Table {
	tt := termtables.CreateTable()
	tt.Style.SkipBorder = true
	tt.Style.BorderX = ""
	tt.Style.BorderY = ""
	tt.Style.BorderI = ""

	tt.AddHeaders(headers...)

	return tt
}

func bold(s string) string {
	if outputIsATty {
		return "\x1b[1m" + s + "\x1b[0m"
	}

	return s
}

var outputIsATty = isatty.IsTerminal(os.Stdout.Fd())

func rightAlignedCell(s string) *termtables.Cell {
	return termtables.CreateCell(
		s,
		&termtables.CellStyle{
			Alignment: termtables.AlignRight,
		},
	)
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
