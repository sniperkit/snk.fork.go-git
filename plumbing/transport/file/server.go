/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package file

import (
	"fmt"
	"os"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/internal/common"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/server"
	"github.com/sniperkit/snk.fork.go-git.v4/utils/ioutil"
)

// ServeUploadPack serves a git-upload-pack request using standard output, input
// and error. This is meant to be used when implementing a git-upload-pack
// command.
func ServeUploadPack(path string) error {
	ep, err := transport.NewEndpoint(path)
	if err != nil {
		return err
	}

	// TODO: define and implement a server-side AuthMethod
	s, err := server.DefaultServer.NewUploadPackSession(ep, nil)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	return common.ServeUploadPack(srvCmd, s)
}

// ServeReceivePack serves a git-receive-pack request using standard output,
// input and error. This is meant to be used when implementing a
// git-receive-pack command.
func ServeReceivePack(path string) error {
	ep, err := transport.NewEndpoint(path)
	if err != nil {
		return err
	}

	// TODO: define and implement a server-side AuthMethod
	s, err := server.DefaultServer.NewReceivePackSession(ep, nil)
	if err != nil {
		return fmt.Errorf("error creating session: %s", err)
	}

	return common.ServeReceivePack(srvCmd, s)
}

var srvCmd = common.ServerCommand{
	Stdin:  os.Stdin,
	Stdout: ioutil.WriteNopCloser(os.Stdout),
	Stderr: os.Stderr,
}
