package auth

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func GetCfAuthRoundTripper(innerRoundTripper http.RoundTripper) (http.RoundTripper, error) {
	getTokenCommand := exec.Command("cf", "oauth-token")
	getTokenOutput, err := getTokenCommand.StdoutPipe()
	err = getTokenCommand.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to launch `cf oauth-token`: %s", err)
	}

	tokenBytes, err := ioutil.ReadAll(getTokenOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to read from `cf oauth-token`: %s", err)
	}

	err = getTokenCommand.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to run `cf oauth-token`: %s", err)
	}

	return &authRoundTripper{
		authorization:     strings.TrimRight(string(tokenBytes), "\r\n"),
		innerRoundTripper: innerRoundTripper,
	}, nil
}

func GetBasicAuthRoundTripper(basicAuth string, innerRoundTripper http.RoundTripper) http.RoundTripper {
	return &authRoundTripper{
		authorization:     "Basic " + base64.StdEncoding.EncodeToString([]byte(basicAuth)),
		innerRoundTripper: innerRoundTripper,
	}
}

type authRoundTripper struct {
	authorization     string
	innerRoundTripper http.RoundTripper
}

func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", a.authorization)

	return a.innerRoundTripper.RoundTrip(req)
}
