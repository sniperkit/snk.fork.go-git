/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package filesystem

import (
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/storer"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem/dotgit"
)

type ReferenceStorage struct {
	dir *dotgit.DotGit
}

func (r *ReferenceStorage) SetReference(ref *plumbing.Reference) error {
	return r.dir.SetRef(ref, nil)
}

func (r *ReferenceStorage) CheckAndSetReference(ref, old *plumbing.Reference) error {
	return r.dir.SetRef(ref, old)
}

func (r *ReferenceStorage) Reference(n plumbing.ReferenceName) (*plumbing.Reference, error) {
	return r.dir.Ref(n)
}

func (r *ReferenceStorage) IterReferences() (storer.ReferenceIter, error) {
	refs, err := r.dir.Refs()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceSliceIter(refs), nil
}

func (r *ReferenceStorage) RemoveReference(n plumbing.ReferenceName) error {
	return r.dir.RemoveRef(n)
}

func (r *ReferenceStorage) CountLooseRefs() (int, error) {
	return r.dir.CountLooseRefs()
}

func (r *ReferenceStorage) PackRefs() error {
	return r.dir.PackRefs()
}
