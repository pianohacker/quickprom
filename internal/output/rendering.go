package output

import (
	"fmt"
	"os"
	"sort"
	"strings"

	isatty "github.com/mattn/go-isatty"
	"github.com/olekukonko/tablewriter"
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
		return FormatVector(value.(model.Vector))
	case model.ValMatrix:
		return FormatMatrix(value.(model.Matrix))
	}

	return nil
}

type tableOutput interface {
	Append([]string)
	Render()
}

func (f *FormattedVector) RenderText() {
	fmt.Print("Instant vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println("")

	fmt.Printf("  At: %s\n", f.Time.Format(TimeFormatWithTZ))

	outputCommonLabels(f.CommonLabels)

	fmt.Println("")

	// Value column
	header := append(f.VaryingLabels, "")

	tw := getTableWriter(header)

	for _, sample := range f.Samples {
		row := sample.LabelValues

		row = append(row, fmt.Sprintf("%f", sample.Value))

		tw.Append(row)
	}

	tw.Render()
}

func (f *FormattedMatrix) RenderText() {
	fmt.Print("Range vector:")
	if f.Empty {
		fmt.Println(" (empty result)")
		return
	}
	fmt.Println("")

	outputCommonLabels(f.CommonLabels)

	var timestampFormat = TimeFormatWithDate
	minDate := f.MinTime.Format(TimeFormatDateOnly)
	maxDate := f.MaxTime.Format(TimeFormatDateOnly)

	if minDate == maxDate {
		fmt.Printf("  All on date: %s\n", minDate)
		timestampFormat = TimeFormat
	}

	fmt.Println("")

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
	fmt.Println("")
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

func bold(s string) string {
	if outputIsATty {
		return "\x1b[1m" + s + "\x1b[0m"
	}

	return s
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
