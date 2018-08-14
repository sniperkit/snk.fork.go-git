/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package idxfile_test

import (
	"bytes"
	"io/ioutil"

	"github.com/sniperkit/snk.fork.go-git-fixtures.v3"
	. "github.com/sniperkit/snk.fork.go-git.v4/plumbing/format/idxfile"

	. "gopkg.in/check.v1"
)

func (s *IdxfileSuite) TestDecodeEncode(c *C) {
	fixtures.ByTag("packfile").Test(c, func(f *fixtures.Fixture) {
		expected, err := ioutil.ReadAll(f.Idx())
		c.Assert(err, IsNil)

		idx := new(MemoryIndex)
		d := NewDecoder(bytes.NewBuffer(expected))
		err = d.Decode(idx)
		c.Assert(err, IsNil)

		result := bytes.NewBuffer(nil)
		e := NewEncoder(result)
		size, err := e.Encode(idx)
		c.Assert(err, IsNil)

		c.Assert(size, Equals, len(expected))
		c.Assert(result.Bytes(), DeepEquals, expected)
	})
}
