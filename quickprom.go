package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	docopt "github.com/docopt/docopt-go"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/pianohacker/quickprom/output"
)

const USAGE = `QuickProm.

Usage:
  quickprom [-t TARGET] [--cf-auth] QUERY

Options:
  -t, --target TARGET  URL of Prometheus-compatible target (QUICKPROM_TARGET)
  --cf-auth            Automatically use current oAuth token from ` + "`cf`" + ` (QUICKPROM_CF_AUTH)
`

type QuickPromOptions struct {
	string `docopt:"--target" env:"QUICKPROM_TARGET"`
	CfAuth bool `docopt:"--cf-auth" env:"QUICKPROM_CF_AUTH"`

	Query string `docopt:"QUERY"`
}

func main() {
	opts := parseOptsAndEnv()

	var roundTripper http.RoundTripper = http.DefaultTransport

	if opts.CfAuth {
		token := getCfOauthToken()
		roundTripper = &oauthRoundTripper{
			token: token,
		}
	}

	promClient := getPromClient(opts.Target, roundTripper)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := promClient.Query(ctx, opts.Query, time.Now())
	failIfErr("Failed to run query: %s", err)

	output.OutputValue(value)
}

func parseOptsAndEnv() *QuickPromOptions {
	var opts QuickPromOptions

	err := envstruct.Load(&opts)
	failIfErr("%s", err)

	parsedOpts, err := docopt.ParseDoc(USAGE)
	failIfErr("%s", err)

	var cmdLineOpts QuickPromOptions
	err = parsedOpts.Bind(&cmdLineOpts)
	failIfErr("%s", err)

	mergeOpts(&opts, &cmdLineOpts)

	if opts.Target == "" {
		fail("Error: Must specify target URL with --target or QUICKPROM_TARGET")
	}

	return &opts
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

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func failIfErr(msg string, args ...interface{}) {
	lastArg := args[len(args)-1]
	if lastArg == nil {
		return
	}

	fail(msg, args...)
}

func getCfOauthToken() string {
	getTokenCommand := exec.Command("cf", "oauth-token")
	getTokenOutput, err := getTokenCommand.StdoutPipe()
	err = getTokenCommand.Start()
	failIfErr("Failed to launch `cf oauth-token`: %s", err)

	tokenBytes, err := ioutil.ReadAll(getTokenOutput)
	failIfErr("Failed to read from `cf oauth-token`: %s", err)

	err = getTokenCommand.Wait()
	failIfErr("Failed to run `cf oauth-token`: %s", err)

	return strings.TrimRight(string(tokenBytes), "\r\n")
}

type oauthRoundTripper struct {
	token string
}

func (o *oauthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", o.token)

	return http.DefaultTransport.RoundTrip(req)
}

func getPromClient(targetAddress string, roundTripper http.RoundTripper) v1.API {
	apiClient, err := api.NewClient(api.Config{
		Address:      targetAddress,
		RoundTripper: roundTripper,
	})
	failIfErr("Failed to initialize Prometheus API: %s", err)

	return v1.NewAPI(apiClient)
}
