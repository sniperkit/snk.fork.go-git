/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package filesystem

import (
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing"
)

type deltaObject struct {
	plumbing.EncodedObject
	base plumbing.Hash
	hash plumbing.Hash
	size int64
}

func newDeltaObject(
	obj plumbing.EncodedObject,
	hash plumbing.Hash,
	base plumbing.Hash,
	size int64) plumbing.DeltaObject {
	return &deltaObject{
		EncodedObject: obj,
		hash:          hash,
		base:          base,
		size:          size,
	}
}

func (o *deltaObject) BaseHash() plumbing.Hash {
	return o.base
}

func (o *deltaObject) ActualSize() int64 {
	return o.size
}

func (o *deltaObject) ActualHash() plumbing.Hash {
	return o.hash
}
