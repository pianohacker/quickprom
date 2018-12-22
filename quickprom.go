package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/pianohacker/quickprom/output"
)

func main() {
	fs := flag.NewFlagSetWithEnvPrefix("quickprom", "QUICKPROM", flag.ExitOnError)

	var (
		cfAuth = fs.Bool("cf-auth", false, "Automatically use current oAuth token from `cf` (QUICKPROM_CF_AUTH)")
		target = fs.String("target", "", "URL of Prometheus-compatible target (QUICKPROM_TARGET)")
	)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--cf-auth] [--target TARGET_URL] <query>\n", os.Args[0])
		fs.PrintDefaults()
	}

	err := fs.Parse(os.Args[1:])
	failIfErr("%s", err)

	if fs.Arg(0) == "" {
		fs.Usage()
		os.Exit(2)
	}

	if *target == "" {
		fail("Error: Must specify target URL with --target or QUICKPROM_TARGET")
	}

	query := fs.Arg(0)

	var roundTripper http.RoundTripper = http.DefaultTransport

	if *cfAuth {
		token := getCfOauthToken()
		roundTripper = &oauthRoundTripper{
			token: token,
		}
	}

	promClient := getPromClient(*target, roundTripper)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := promClient.Query(ctx, query, time.Now())
	failIfErr("Failed to run query: %s", err)

	output.OutputValue(value)
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
