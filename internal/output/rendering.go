package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/apcera/termtables"
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
	case model.ValVector:
		return FormatInstantVector(value.(model.Vector))
	case model.ValMatrix:
		return FormatRangeVector(value.(model.Matrix))
	}

	return nil
}

func (f *FormattedInstantVector) RenderText(_ *RenderOptions) {
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

		row = append(row, termtables.CreateCell(
			fmt.Sprintf("%f", sample.Value),
			&termtables.CellStyle{
				Alignment: termtables.AlignRight,
			},
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

	outputCommonLabels(f.CommonLabels)

	sharedDateParts := SharedDateParts(f.SeenTimes)

	if sharedDateParts.Date {
		fmt.Printf("  All on date: %s\n", f.SeenTimes[0].Format(TimeFormatDateOnly))
	}

	if sharedDateParts.ZeroSecond {
		fmt.Println("  All timestamps end with: 00.000")
	} else if sharedDateParts.ZeroMillisecond {
		fmt.Println("  All timestamps end with: .000")
	}

	timestampFormat := getTimestampFormat(sharedDateParts)

	if opts.RangeVectorAsTable {
		// Value column
		header := f.VaryingLabels

		for _, seenTime := range f.SeenTimes {
			header = append(header, seenTime.Format(timestampFormat))
		}

		tw := getTableWriter(header)

		for _, series := range f.Series {
			var row []interface{}

			for _, labelValue := range series.LabelValues {
				row = append(row, labelValue)
			}

			samplePos := 0
			for _, seenTime := range f.SeenTimes {
				for samplePos < len(series.Values) && series.Values[samplePos].Time != seenTime {
					samplePos++
				}

				if samplePos < len(series.Values) {
					row = append(row, termtables.CreateCell(
						fmt.Sprintf("%f", series.Values[samplePos].Value),
						&termtables.CellStyle{
							Alignment: termtables.AlignRight,
						},
					))
				}
			}

			tw.AddRow(row...)
		}

		for i := len(f.VaryingLabels); i < len(header); i++ {
			tw.SetAlign(termtables.AlignRight, i+1)
		}

		fmt.Print(tw.Render())
	} else {
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

func getTableWriter(headers []string) *termtables.Table {
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

	return tt
}

func bold(s string) string {
	if outputIsATty {
		return "\x1b[1m" + s + "\x1b[0m"
	}

	return s
}

var outputIsATty = isatty.IsTerminal(os.Stdout.Fd())

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
