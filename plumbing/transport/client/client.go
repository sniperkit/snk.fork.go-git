/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

// Package client contains helper function to deal with the different client
// protocols.
package client

import (
	"fmt"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/file"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/git"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/http"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/ssh"
)

// Protocols are the protocols supported by default.
var Protocols = map[string]transport.Transport{
	"http":  http.DefaultClient,
	"https": http.DefaultClient,
	"ssh":   ssh.DefaultClient,
	"git":   git.DefaultClient,
	"file":  file.DefaultClient,
}

// InstallProtocol adds or modifies an existing protocol.
func InstallProtocol(scheme string, c transport.Transport) {
	if c == nil {
		delete(Protocols, scheme)
		return
	}

	Protocols[scheme] = c
}

// NewClient returns the appropriate client among of the set of known protocols:
// http://, https://, ssh:// and file://.
// See `InstallProtocol` to add or modify protocols.
func NewClient(endpoint *transport.Endpoint) (transport.Transport, error) {
	f, ok := Protocols[endpoint.Protocol]
	if !ok {
		return nil, fmt.Errorf("unsupported scheme %q", endpoint.Protocol)
	}

	if f == nil {
		return nil, fmt.Errorf("malformed client for scheme %q, client is defined as nil", endpoint.Protocol)
	}

	return f, nil
}
