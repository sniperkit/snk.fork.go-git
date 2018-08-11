/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package server_test

import (
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport"

	. "gopkg.in/check.v1"
)

type ReceivePackSuite struct {
	BaseSuite
}

var _ = Suite(&ReceivePackSuite{})

func (s *ReceivePackSuite) SetUpSuite(c *C) {
	s.BaseSuite.SetUpSuite(c)
	s.ReceivePackSuite.Client = s.client
}

func (s *ReceivePackSuite) SetUpTest(c *C) {
	s.prepareRepositories(c)
}

func (s *ReceivePackSuite) TearDownTest(c *C) {
	s.Suite.TearDownSuite(c)
}

// Overwritten, server returns error earlier.
func (s *ReceivePackSuite) TestAdvertisedReferencesNotExists(c *C) {
	r, err := s.Client.NewReceivePackSession(s.NonExistentEndpoint, s.EmptyAuth)
	c.Assert(err, Equals, transport.ErrRepositoryNotFound)
	c.Assert(r, IsNil)
}
