/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package storer

import "github.com/sniperkit/snk.fork.go-git.v4/plumbing/format/index"

// IndexStorer generic storage of index.Index
type IndexStorer interface {
    SetIndex(*index.Index) error
    Index() (*index.Index, error)
}
