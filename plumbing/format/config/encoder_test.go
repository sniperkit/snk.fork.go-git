/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package config

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type EncoderSuite struct{}

var _ = Suite(&EncoderSuite{})

func (s *EncoderSuite) TestEncode(c *C) {
	for idx, fixture := range fixtures {
		buf := &bytes.Buffer{}
		e := NewEncoder(buf)
		err := e.Encode(fixture.Config)
		c.Assert(err, IsNil, Commentf("encoder error for fixture: %d", idx))
		c.Assert(buf.String(), Equals, fixture.Text, Commentf("bad result for fixture: %d", idx))
	}
}
