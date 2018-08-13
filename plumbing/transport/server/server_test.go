/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package server_test

import (
	"testing"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/client"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/server"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/test"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/memory"

	"github.com/sniperkit/snk.fork.go-git-fixtures.v3"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	fixtures.Suite
	test.ReceivePackSuite

	loader       server.MapLoader
	client       transport.Transport
	clientBackup transport.Transport
	asClient     bool
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)
	s.loader = server.MapLoader{}
	if s.asClient {
		s.client = server.NewClient(s.loader)
	} else {
		s.client = server.NewServer(s.loader)
	}

	s.clientBackup = client.Protocols["file"]
	client.Protocols["file"] = s.client
}

func (s *BaseSuite) TearDownSuite(c *C) {
	if s.clientBackup == nil {
		delete(client.Protocols, "file")
	} else {
		client.Protocols["file"] = s.clientBackup
	}
}

func (s *BaseSuite) prepareRepositories(c *C) {
	var err error

	fs := fixtures.Basic().One().DotGit()
	s.Endpoint, err = transport.NewEndpoint(fs.Root())
	c.Assert(err, IsNil)
	s.loader[s.Endpoint.String()], err = filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	s.EmptyEndpoint, err = transport.NewEndpoint("/empty.git")
	c.Assert(err, IsNil)
	s.loader[s.EmptyEndpoint.String()] = memory.NewStorage()

	s.NonExistentEndpoint, err = transport.NewEndpoint("/non-existent.git")
	c.Assert(err, IsNil)
}
