/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package object

import (
	fixtures "github.com/sniperkit/snk.fork.go-git-fixtures.v3"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem"
	. "gopkg.in/check.v1"
)

type PatchSuite struct {
	BaseObjectsSuite
}

var _ = Suite(&PatchSuite{})

func (s *PatchSuite) TestStatsWithSubmodules(c *C) {
	storer, err := filesystem.NewStorage(
		fixtures.ByURL("https://github.com/git-fixtures/submodule.git").One().DotGit())

	commit, err := GetCommit(storer, plumbing.NewHash("b685400c1f9316f350965a5993d350bc746b0bf4"))

	tree, err := commit.Tree()
	c.Assert(err, IsNil)

	e, err := tree.entry("basic")
	c.Assert(err, IsNil)

	ch := &Change{
		From: ChangeEntry{
			Name:      "basic",
			Tree:      tree,
			TreeEntry: *e,
		},
		To: ChangeEntry{
			Name:      "basic",
			Tree:      tree,
			TreeEntry: *e,
		},
	}

	p, err := getPatch("", ch)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
}
