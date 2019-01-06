package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/pianohacker/quickprom/internal/auth"
	"github.com/pianohacker/quickprom/internal/cmdline"
	"github.com/pianohacker/quickprom/internal/output"
)

func main() {
	opts, err := cmdline.ParseOptsAndEnv(true)
	failIfErr("Error: %s", err)

	promClient := getPromClient(opts)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var value model.Value
	if opts.RangeEnabled {
		value, err = promClient.QueryRange(ctx, opts.Query, v1.Range{
			Start: opts.RangeStart,
			End:   opts.RangeEnd,
			Step:  opts.RangeStep,
		})
	} else {
		value, err = promClient.Query(ctx, opts.Query, opts.Time)
	}
	failIfErr("Failed to run query: %s", err)

	if opts.Json {
		failIfErr("Failed to marshal result to JSON: %s", output.RenderJson(value))
	} else {
		output.FormatValue(value).RenderText()
	}
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func failIfErr(msg string, err error) {
	if err == nil {
		return
	}

	fail(msg, err)
}

func getPromClient(opts *cmdline.QuickPromOptions) v1.API {
	apiClient, err := api.NewClient(api.Config{
		Address:      opts.Target,
		RoundTripper: getRoundTripper(opts),
	})
	failIfErr("Failed to initialize Prometheus API: %s", err)

	return v1.NewAPI(apiClient)
}

func getRoundTripper(opts *cmdline.QuickPromOptions) http.RoundTripper {
	var roundTripper http.RoundTripper = api.DefaultRoundTripper

	if opts.SkipTlsVerify {
		roundTripper.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if opts.CfAuth {
		var err error
		roundTripper, err = auth.CfAuthRoundTripper(roundTripper)
		failIfErr("Error: %s", err)
	} else if opts.BasicAuth != "" {
		roundTripper = auth.BasicAuthRoundTripper(opts.BasicAuth, roundTripper)
	}

	return roundTripper
}
