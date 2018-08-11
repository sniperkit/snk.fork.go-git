/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package git

import (
	. "gopkg.in/check.v1"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/object"
)

type OptionsSuite struct {
	BaseSuite
}

var _ = Suite(&OptionsSuite{})

func (s *OptionsSuite) TestCommitOptionsParentsFromHEAD(c *C) {
	o := CommitOptions{Author: &object.Signature{}}
	err := o.Validate(s.Repository)
	c.Assert(err, IsNil)
	c.Assert(o.Parents, HasLen, 1)
}

func (s *OptionsSuite) TestCommitOptionsMissingAuthor(c *C) {
	o := CommitOptions{}
	err := o.Validate(s.Repository)
	c.Assert(err, Equals, ErrMissingAuthor)
}

func (s *OptionsSuite) TestCommitOptionsCommitter(c *C) {
	sig := &object.Signature{}

	o := CommitOptions{Author: sig}
	err := o.Validate(s.Repository)
	c.Assert(err, IsNil)

	c.Assert(o.Committer, Equals, o.Author)
}
