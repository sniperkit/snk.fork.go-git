/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package storer

import "github.com/sniperkit/snk.fork.go-git.v4/plumbing"

// ShallowStorer is a storage of references to shallow commits by hash,
// meaning that these commits have missing parents because of a shallow fetch.
type ShallowStorer interface {
    SetShallow([]plumbing.Hash) error
    Shallow() ([]plumbing.Hash, error)
}
