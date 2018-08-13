/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package server

import (
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/storer"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/transport"
	"github.com/sniperkit/snk.fork.go-git.v4/storage/filesystem"

	"github.com/sniperkit/snk.fork.go-billy.v4"
	"github.com/sniperkit/snk.fork.go-billy.v4/osfs"
)

// DefaultLoader is a filesystem loader ignoring host and resolving paths to /.
var DefaultLoader = NewFilesystemLoader(osfs.New(""))

// Loader loads repository's storer.Storer based on an optional host and a path.
type Loader interface {
	// Load loads a storer.Storer given a transport.Endpoint.
	// Returns transport.ErrRepositoryNotFound if the repository does not
	// exist.
	Load(ep *transport.Endpoint) (storer.Storer, error)
}

type fsLoader struct {
	base billy.Filesystem
}

// NewFilesystemLoader creates a Loader that ignores host and resolves paths
// with a given base filesystem.
func NewFilesystemLoader(base billy.Filesystem) Loader {
	return &fsLoader{base}
}

// Load looks up the endpoint's path in the base file system and returns a
// storer for it. Returns transport.ErrRepositoryNotFound if a repository does
// not exist in the given path.
func (l *fsLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	fs, err := l.base.Chroot(ep.Path)
	if err != nil {
		return nil, err
	}

	if _, err := fs.Stat("config"); err != nil {
		return nil, transport.ErrRepositoryNotFound
	}

	return filesystem.NewStorage(fs)
}

// MapLoader is a Loader that uses a lookup map of storer.Storer by
// transport.Endpoint.
type MapLoader map[string]storer.Storer

// Load returns a storer.Storer for given a transport.Endpoint by looking it up
// in the map. Returns transport.ErrRepositoryNotFound if the endpoint does not
// exist.
func (l MapLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	s, ok := l[ep.String()]
	if !ok {
		return nil, transport.ErrRepositoryNotFound
	}

	return s, nil
}
