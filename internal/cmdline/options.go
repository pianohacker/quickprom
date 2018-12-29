package cmdline

import (
	"errors"
	"os"
	"reflect"

	envstruct "code.cloudfoundry.org/go-envstruct"
	docopt "github.com/docopt/docopt-go"
)

const USAGE = `QuickProm.

Usage:
  quickprom [-t TARGET] [--cf-auth] QUERY

Options:
  -t, --target TARGET  URL of Prometheus-compatible target (QUICKPROM_TARGET)
  --cf-auth            Automatically use current oAuth token from ` + "`cf`" + ` (QUICKPROM_CF_AUTH)
`

type QuickPromOptions struct {
	Target string `docopt:"--target" env:"QUICKPROM_TARGET"`
	CfAuth bool   `docopt:"--cf-auth" env:"QUICKPROM_CF_AUTH"`

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
