/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package gitignore

import (
	. "gopkg.in/check.v1"
)

func (s *MatcherSuite) TestMatcher_Match(c *C) {
	ps := []Pattern{
		ParsePattern("**/middle/v[uo]l?ano", nil),
		ParsePattern("!volcano", nil),
	}

	m := NewMatcher(ps)
	c.Assert(m.Match([]string{"head", "middle", "vulkano"}, false), Equals, true)
	c.Assert(m.Match([]string{"head", "middle", "volcano"}, false), Equals, false)
}
