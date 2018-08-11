/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package git

import (
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport/test"

	. "gopkg.in/check.v1"
	"github.com/sniperkit/snk.fork.go-git-fixtures.v3"
)

type ReceivePackSuite struct {
	test.ReceivePackSuite
	BaseSuite
}

var _ = Suite(&ReceivePackSuite{})

func (s *ReceivePackSuite) SetUpTest(c *C) {
	s.BaseSuite.SetUpTest(c)

	s.ReceivePackSuite.Client = DefaultClient
	s.ReceivePackSuite.Endpoint = s.prepareRepository(c, fixtures.Basic().One(), "basic.git")
	s.ReceivePackSuite.EmptyEndpoint = s.prepareRepository(c, fixtures.ByTag("empty").One(), "empty.git")
	s.ReceivePackSuite.NonExistentEndpoint = s.newEndpoint(c, "non-existent.git")

	s.StartDaemon(c)
}
