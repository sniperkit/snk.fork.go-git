/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package packp

import (
	"bytes"
	"io"
	"testing"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/format/pktline"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// returns a byte slice with the pkt-lines for the given payloads.
func pktlines(c *C, payloads ...string) []byte {
	var buf bytes.Buffer
	e := pktline.NewEncoder(&buf)

	err := e.EncodeString(payloads...)
	c.Assert(err, IsNil, Commentf("building pktlines for %v\n", payloads))

	return buf.Bytes()
}

func toPktLines(c *C, payloads []string) io.Reader {
	var buf bytes.Buffer
	e := pktline.NewEncoder(&buf)
	err := e.EncodeString(payloads...)
	c.Assert(err, IsNil)

	return &buf
}
