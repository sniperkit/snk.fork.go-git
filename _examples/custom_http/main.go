/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sniperkit/snk.fork.go-git.v4"
	. "github.com/sniperkit/snk.fork.go-git.v4/_examples"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/client"
	githttp "github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/http"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/memory"
)

// Here is an example to configure http client according to our own needs.
func main() {
	CheckArgs("<url>")
	url := os.Args[1]

	// Create a custom http(s) client with your config
	customClient := &http.Client{
		// accept any certificate (might be useful for testing)
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},

		// 15 second timeout
		Timeout: 15 * time.Second,

		// don't follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(customClient))

	// Clone repository using the new client if the protocol is https://
	Info("git clone %s", url)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	CheckIfError(err)

	// Retrieve the branch pointed by HEAD
	Info("git rev-parse HEAD")

	head, err := r.Head()
	CheckIfError(err)
	fmt.Println(head.Hash())
}
