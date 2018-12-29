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

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/pianohacker/quickprom/internal/cmdline"
	"github.com/pianohacker/quickprom/internal/output"
)

func main() {
	opts, err := cmdline.ParseOptsAndEnv(true)
	failIfErr("Error: %s", err)

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
