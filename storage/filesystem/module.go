/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package filesystem

import (
    "github.com/sniperkit/snk.fork.go-git.v4/storage"
    "github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem/dotgit"
)

type ModuleStorage struct {
    dir *dotgit.DotGit
}

func (s *ModuleStorage) Module(name string) (storage.Storer, error) {
    fs, err := s.dir.Module(name)
    if err != nil {
        return nil, err
    }

    return NewStorage(fs)
}
