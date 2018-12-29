package cmdline

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	fuzzytime "github.com/bcampbell/fuzzytime"
	docopt "github.com/docopt/docopt-go"
	"github.com/prometheus/common/model"
)

const USAGE = `QuickProm.

Usage:
  quickprom [-t TARGET] [--cf-auth] QUERY [--time TIME]
  quickprom [-t TARGET] [--cf-auth] range QUERY --start START [--end END] --step STEP

Options:
  -t, --target TARGET  URL of Prometheus-compatible target (QUICKPROM_TARGET)
  --cf-auth            Automatically use current oAuth token from ` + "`cf`" + ` (QUICKPROM_CF_AUTH)
  --time TIME          Evaluation time of instant query (defaults to now)
  --start START        Start time of range query
  --end END            End time of range query (inclusive, defaults to now)
  --step STEP          Step of range query
`

type QuickPromOptions struct {
	Target string `docopt:"--target" env:"QUICKPROM_TARGET"`
	CfAuth bool   `docopt:"--cf-auth" env:"QUICKPROM_CF_AUTH"`

	TimeInput string `docopt:"--time"`
	Time time.Time

	RangeEnabled    bool   `docopt:"range"`
	RangeStartInput string `docopt:"--start"`
	RangeStart      time.Time
	RangeEndInput   string `docopt:"--end"`
	RangeEnd        time.Time
	RangeStepInput  string `docopt:"--step"`
	RangeStep       time.Duration

	Query string `docopt:"QUERY"`
}

func ParseOptsAndEnv(exitOnError bool) (*QuickPromOptions, error) {
	var opts QuickPromOptions

	err := envstruct.Load(&opts)
	if err != nil {
		return nil, err
	}

	var helpHandler func(error, string)
	if exitOnError {
		helpHandler = docopt.PrintHelpAndExit
	} else {
		helpHandler = docopt.NoHelpHandler
	}

	parser := &docopt.Parser{
		HelpHandler: helpHandler,
	}

	parsedOpts, err := parser.ParseArgs(USAGE, os.Args[1:], "")
	if err != nil {
		return nil, err
	}

	var cmdLineOpts QuickPromOptions
	err = parsedOpts.Bind(&cmdLineOpts)
	if err != nil {
		return nil, err
	}

	mergeOpts(&opts, &cmdLineOpts)

	if opts.Target == "" {
		return nil, errors.New("must specify target URL with --target or QUICKPROM_TARGET")
	}

	if opts.RangeEnabled {
		opts.RangeStart, err = ParseTime(opts.RangeStartInput)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --start: %s", err)
		}

		if opts.RangeEndInput == "" {
			opts.RangeEnd = time.Now()
		} else {
			opts.RangeEnd, err = ParseTime(opts.RangeEndInput)
			if err != nil {
				return nil, fmt.Errorf("failed to parse --end: %s", err)
			}
		}

		if opts.RangeEnd.Before(opts.RangeStart) {
			return nil, errors.New("--end before --start")
		}

		parsedStep, err := model.ParseDuration(opts.RangeStepInput)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --step: %s", err)
		}

		opts.RangeStep = time.Duration(parsedStep)
	} else {
		if opts.TimeInput == "" {
			opts.Time = time.Now()
		} else {
			opts.Time, err = ParseTime(opts.TimeInput)
			if err != nil {
				return nil, fmt.Errorf("failed to parse --time: %s", err)
			}
		}
	}

	return &opts, nil
}

func mergeOpts(destOpts, srcOpts *QuickPromOptions) {
	destOptsVal := reflect.ValueOf(destOpts).Elem()
	srcOptsVal := reflect.ValueOf(srcOpts).Elem()

	for i := 0; i < destOptsVal.NumField(); i++ {
		destFieldVal := destOptsVal.Field(i)
		srcFieldVal := srcOptsVal.Field(i)

		zeroVal := reflect.Zero(destFieldVal.Type()).Interface()

		if !reflect.DeepEqual(srcFieldVal.Interface(), zeroVal) {
			destFieldVal.Set(srcFieldVal)
		}
	}
}

func ParseTime(s string) (time.Time, error) {
	dateTime, _, err := fuzzytime.WesternContext.Extract(s)

	if dateTime.Empty() {
		return time.Time{}, errors.New("no valid time found")
	}

	if err != nil {
		return time.Time{}, err
	}

	loc := time.Local

	if dateTime.HasTZOffset() {
		loc = time.FixedZone("", dateTime.TZOffset())
	}

	now := time.Now()

	if dateTime.Date.Empty() {
		dateTime.Date.SetYear(now.Year())
		dateTime.Date.SetMonth(int(now.Month()))
		dateTime.Date.SetDay(now.Day())
	}

	return time.Date(
		maybeInt(dateTime.HasYear, dateTime.Year, now.Year()),
		time.Month(dateTime.Month()),
		maybeInt(dateTime.HasDay, dateTime.Day, 1),
		maybeInt(dateTime.HasHour, dateTime.Hour, 0),
		maybeInt(dateTime.HasMinute, dateTime.Minute, 0),
		maybeInt(dateTime.HasSecond, dateTime.Second, 0),
		0,
		loc,
	), nil
}

func maybeInt(isSet func() bool, getter func() int, def int) int {
	if isSet() {
		return getter()
	} else {
		return def
	}
}
