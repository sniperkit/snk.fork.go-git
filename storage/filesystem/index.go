/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package filesystem

import (
	"os"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/format/index"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem/dotgit"
	"github.com/sniperkit/snk.fork.go-git.v4/utils/ioutil"
)

type IndexStorage struct {
	dir *dotgit.DotGit
}

func (s *IndexStorage) SetIndex(idx *index.Index) (err error) {
	f, err := s.dir.IndexWriter()
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(f, &err)

	e := index.NewEncoder(f)
	err = e.Encode(idx)
	return err
}

func (s *IndexStorage) Index() (i *index.Index, err error) {
	idx := &index.Index{
		Version: 2,
	}

	f, err := s.dir.Index()
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}

		return nil, err
	}

	defer ioutil.CheckClose(f, &err)

	d := index.NewDecoder(f)
	err = d.Decode(idx)
	return idx, err
}
