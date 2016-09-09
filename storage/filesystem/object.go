package filesystem

import (
	"fmt"
	"io"

	"gopkg.in/src-d/go-git.v4/core"
	"gopkg.in/src-d/go-git.v4/formats/objfile"
	"gopkg.in/src-d/go-git.v4/formats/packfile"
	"gopkg.in/src-d/go-git.v4/storage/filesystem/internal/dotgit"
	"gopkg.in/src-d/go-git.v4/storage/filesystem/internal/index"
	"gopkg.in/src-d/go-git.v4/utils/fs"
)

// ObjectStorage is an implementation of core.ObjectStorage that stores
// data on disk in the standard git format (this is, the .git directory).
//
// Zero values of this type are not safe to use, see the New function below.
//
// Currently only reads are supported, no writting.
//
// Also values from this type are not yet able to track changes on disk, this is,
// Gitdir values will get outdated as soon as repositories change on disk.
type ObjectStorage struct {
	dir   *dotgit.DotGit
	index index.Index
}

func (s *ObjectStorage) NewObject() core.Object {
	return &core.MemoryObject{}
}

// Writer method not supported on Memory storage
func (o *ObjectStorage) Writer() (io.WriteCloser, error) {
	return o.dir.NewObjectPack()
}

// Set adds a new object to the storage. As this functionality is not
// yet supported, this method always returns a "not implemented yet"
// error an zero hash.
func (s *ObjectStorage) Set(core.Object) (core.Hash, error) {
	return core.ZeroHash, fmt.Errorf("not implemented yet")
}

// Get returns the object with the given hash, by searching for it in
// the packfile and the git object directories.
func (s *ObjectStorage) Get(t core.ObjectType, h core.Hash) (core.Object, error) {
	obj, err := s.getFromUnpacked(t, h)
	if err == dotgit.ErrObjfileNotFound {
		if s.index == nil {
			return nil, core.ErrObjectNotFound
		}
		return s.getFromPackfile(t, h)
	}

	return obj, err
}

func (s *ObjectStorage) getFromUnpacked(t core.ObjectType, h core.Hash) (obj core.Object, err error) {
	f, err := s.dir.Object(h)
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	obj = s.NewObject()
	objReader, err := objfile.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer func() {
		errClose := objReader.Close()
		if err == nil {
			err = errClose
		}
	}()

	err = objReader.FillObject(obj)
	if err != nil {
		return nil, err
	}
	if core.AnyObject != t && obj.Type() != t {
		return nil, core.ErrObjectNotFound
	}
	return obj, nil
}

// Get returns the object with the given hash, by searching for it in
// the packfile.
func (s *ObjectStorage) getFromPackfile(t core.ObjectType, h core.Hash) (core.Object, error) {
	offset, err := s.index.Get(h)
	if err != nil {
		return nil, err
	}

	packs, err := s.dir.ObjectsPacks()
	if err != nil {
		return nil, err
	}

	if len(packs) == 0 {
		return nil, nil
	}

	f, _, err := s.dir.ObjectPack(packs[0].Name())
	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	p := packfile.NewScanner(f)
	d := packfile.NewDecoder(p, nil)

	obj, err := d.ReadObjectAt(offset)
	if err != nil {
		return nil, err
	}

	if core.AnyObject != t && obj.Type() != t {
		return nil, core.ErrObjectNotFound
	}

	return obj, nil
}

// Iter returns an iterator for all the objects in the packfile with the
// given type.
func (s *ObjectStorage) Iter(t core.ObjectType) (core.ObjectIter, error) {
	var objects []core.Object

	hashes, err := s.dir.Objects()
	if err != nil {
		return nil, err
	}

	for _, hash := range hashes {
		object, err := s.getFromUnpacked(core.AnyObject, hash)
		if err != nil {
			return nil, err
		}
		if object.Type() == t {
			objects = append(objects, object)
		}
	}

	for hash := range s.index {
		object, err := s.getFromPackfile(core.AnyObject, hash)
		if err != nil {
			return nil, err
		}
		if t == core.AnyObject || object.Type() == t {
			objects = append(objects, object)
		}
	}

	return core.NewObjectSliceIter(objects), nil
}

func buildIndex(fs fs.Filesystem, dir *dotgit.DotGit) (index.Index, error) {
	packs, err := dir.ObjectsPacks()
	if err != nil {
		return nil, err
	}

	if len(packs) == 0 {
		return nil, nil
	}

	_, _, err = dir.ObjectPack(packs[0].Name())
	if err != nil {
		if err == dotgit.ErrIdxNotFound {
			return buildIndexFromPackfile(dir)
		}
		return nil, err
	}

	return buildIndexFromIdxfile(fs, "")
}

func buildIndexFromPackfile(dir *dotgit.DotGit) (index.Index, error) {
	packs, err := dir.ObjectsPacks()
	if err != nil {
		return nil, err
	}

	if len(packs) == 0 {
		return nil, nil
	}

	f, _, err := dir.ObjectPack(packs[0].Name())
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	index, _, err := index.NewFromPackfile(f)
	return index, err
}

func buildIndexFromIdxfile(fs fs.Filesystem, path string) (index.Index, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	i := index.New()
	return i, i.Decode(f)
}

func (o *ObjectStorage) Begin() core.TxObjectStorage {
	return &TxObjectStorage{}
}

type TxObjectStorage struct{}

func (tx *TxObjectStorage) Set(obj core.Object) (core.Hash, error) {
	return core.ZeroHash, fmt.Errorf("not implemented yet")
}

func (tx *TxObjectStorage) Get(core.ObjectType, core.Hash) (core.Object, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (tx *TxObjectStorage) Commit() error {
	return fmt.Errorf("not implemented yet")
}

func (tx *TxObjectStorage) Rollback() error {
	return fmt.Errorf("not implemented yet")
}