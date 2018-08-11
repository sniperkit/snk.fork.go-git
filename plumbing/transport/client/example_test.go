/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package client_test

import (
	"crypto/tls"
	"net/http"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/client"
	githttp "github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/http"
)

func ExampleInstallProtocol() {
	// Create custom net/http client that.
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Install it as default client for https URLs.
	client.InstallProtocol("https", githttp.NewClient(httpClient))
}
