package packfile

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	// ErrReferenceDeltaNotFound is returned when the reference delta is not
	// found.
	ErrReferenceDeltaNotFound = errors.New("reference delta not found")

	// ErrNotSeekableSource is returned when the source for the parser is not
	// seekable and a storage was not provided, so it can't be parsed.
	ErrNotSeekableSource = errors.New("parser source is not seekable and storage was not provided")

	// ErrDeltaNotCached is returned when the delta could not be found in cache.
	ErrDeltaNotCached = errors.New("delta could not be found in cache")
)

// Observer interface is implemented by index encoders.
type Observer interface {
	// OnHeader is called when a new packfile is opened.
	OnHeader(count uint32) error
	// OnInflatedObjectHeader is called for each object header read.
	OnInflatedObjectHeader(t plumbing.ObjectType, objSize int64, pos int64) error
	// OnInflatedObjectContent is called for each decoded object.
	OnInflatedObjectContent(h plumbing.Hash, pos int64, crc uint32, content []byte) error
	// OnFooter is called when decoding is done.
	OnFooter(h plumbing.Hash) error
}

// Parser decodes a packfile and calls any observer associated to it. Is used
// to generate indexes.
type Parser struct {
	storage          storer.EncodedObjectStorer
	scanner          *Scanner
	count            uint32
	oi               []*objectInfo
	oiByHash         map[plumbing.Hash]*objectInfo
	oiByOffset       map[int64]*objectInfo
	hashOffset       map[plumbing.Hash]int64
	pendingRefDeltas map[plumbing.Hash][]*objectInfo
	checksum         plumbing.Hash

	cache *cache.ObjectLRU
	// delta content by offset, only used if source is not seekable
	deltas map[int64][]byte

	ob []Observer
}

// NewParser creates a new Parser. The Scanner source must be seekable.
// If it's not, NewParserWithStorage should be used instead.
func NewParser(scanner *Scanner, ob ...Observer) (*Parser, error) {
	return NewParserWithStorage(scanner, nil, ob...)
}

// NewParserWithStorage creates a new Parser. The scanner source must either
// be seekable or a storage must be provided.
func NewParserWithStorage(
	scanner *Scanner,
	storage storer.EncodedObjectStorer,
	ob ...Observer,
) (*Parser, error) {
	if !scanner.IsSeekable && storage == nil {
		return nil, ErrNotSeekableSource
	}

	var deltas map[int64][]byte
	if !scanner.IsSeekable {
		deltas = make(map[int64][]byte)
	}

	return &Parser{
		storage:          storage,
		scanner:          scanner,
		ob:               ob,
		count:            0,
		cache:            cache.NewObjectLRUDefault(),
		pendingRefDeltas: make(map[plumbing.Hash][]*objectInfo),
		deltas:           deltas,
	}, nil
}

func (p *Parser) forEachObserver(f func(o Observer) error) error {
	for _, o := range p.ob {
		if err := f(o); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) onHeader(count uint32) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnHeader(count)
	})
}

func (p *Parser) onInflatedObjectHeader(
	t plumbing.ObjectType,
	objSize int64,
	pos int64,
) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnInflatedObjectHeader(t, objSize, pos)
	})
}

func (p *Parser) onInflatedObjectContent(
	h plumbing.Hash,
	pos int64,
	crc uint32,
	content []byte,
) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnInflatedObjectContent(h, pos, crc, content)
	})
}

func (p *Parser) onFooter(h plumbing.Hash) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnFooter(h)
	})
}

// Parse start decoding phase of the packfile.
func (p *Parser) Parse() (plumbing.Hash, error) {
	if err := p.init(); err != nil {
		return plumbing.ZeroHash, err
	}

	if err := p.indexObjects(); err != nil {
		return plumbing.ZeroHash, err
	}

	var err error
	p.checksum, err = p.scanner.Checksum()
	if err != nil && err != io.EOF {
		return plumbing.ZeroHash, err
	}

	if err := p.resolveDeltas(); err != nil {
		return plumbing.ZeroHash, err
	}

	if len(p.pendingRefDeltas) > 0 {
		return plumbing.ZeroHash, ErrReferenceDeltaNotFound
	}

	if err := p.onFooter(p.checksum); err != nil {
		return plumbing.ZeroHash, err
	}

	return p.checksum, nil
}

func (p *Parser) init() error {
	_, c, err := p.scanner.Header()
	if err != nil {
		return err
	}

	if err := p.onHeader(c); err != nil {
		return err
	}

	p.count = c
	p.oiByHash = make(map[plumbing.Hash]*objectInfo, p.count)
	p.oiByOffset = make(map[int64]*objectInfo, p.count)
	p.oi = make([]*objectInfo, p.count)

	return nil
}

func (p *Parser) indexObjects() error {
	buf := new(bytes.Buffer)

	for i := uint32(0); i < p.count; i++ {
		buf.Reset()

		oh, err := p.scanner.NextObjectHeader()
		if err != nil {
			return err
		}

		delta := false
		var ota *objectInfo
		switch t := oh.Type; t {
		case plumbing.OFSDeltaObject:
			delta = true

			parent, ok := p.oiByOffset[oh.OffsetReference]
			if !ok {
				return plumbing.ErrObjectNotFound
			}

			ota = newDeltaObject(oh.Offset, oh.Length, t, parent)
			parent.Children = append(parent.Children, ota)
		case plumbing.REFDeltaObject:
			delta = true

			parent, ok := p.oiByHash[oh.Reference]
			if ok {
				ota = newDeltaObject(oh.Offset, oh.Length, t, parent)
				parent.Children = append(parent.Children, ota)
			} else {
				ota = newBaseObject(oh.Offset, oh.Length, t)
				p.pendingRefDeltas[oh.Reference] = append(
					p.pendingRefDeltas[oh.Reference],
					ota,
				)
			}
		default:
			ota = newBaseObject(oh.Offset, oh.Length, t)
		}

		size, crc, err := p.scanner.NextObject(buf)
		if err != nil {
			return err
		}

		ota.Crc32 = crc
		ota.PackSize = size
		ota.Length = oh.Length

		data := buf.Bytes()
		if !delta {
			if _, err := ota.Write(data); err != nil {
				return err
			}
			ota.SHA1 = ota.Sum()
			p.oiByHash[ota.SHA1] = ota
		}

		if p.storage != nil && !delta {
			obj := new(plumbing.MemoryObject)
			obj.SetSize(oh.Length)
			obj.SetType(oh.Type)
			if _, err := obj.Write(data); err != nil {
				return err
			}

			if _, err := p.storage.SetEncodedObject(obj); err != nil {
				return err
			}
		}

		if delta && !p.scanner.IsSeekable {
			p.deltas[oh.Offset] = make([]byte, len(data))
			copy(p.deltas[oh.Offset], data)
		}

		p.oiByOffset[oh.Offset] = ota
		p.oi[i] = ota
	}

	return nil
}

func (p *Parser) resolveDeltas() error {
	for _, obj := range p.oi {
		content, err := p.get(obj)
		if err != nil {
			return err
		}

		if err := p.onInflatedObjectHeader(obj.Type, obj.Length, obj.Offset); err != nil {
			return err
		}

		if err := p.onInflatedObjectContent(obj.SHA1, obj.Offset, obj.Crc32, content); err != nil {
			return err
		}

		if !obj.IsDelta() && len(obj.Children) > 0 {
			for _, child := range obj.Children {
				if _, err := p.resolveObject(child, content); err != nil {
					return err
				}
			}

			// Remove the delta from the cache.
			if obj.DiskType.IsDelta() && !p.scanner.IsSeekable {
				delete(p.deltas, obj.Offset)
			}
		}

		obj.Content = nil
	}

	return nil
}

func (p *Parser) get(o *objectInfo) ([]byte, error) {
	if len(o.Content) > 0 {
		return o.Content, nil
	}

	e, ok := p.cache.Get(o.SHA1)
	// If it's not on the cache and is not a delta we can try to find it in the
	// storage, if there's one.
	if !ok && p.storage != nil && !o.Type.IsDelta() {
		var err error
		e, err = p.storage.EncodedObject(plumbing.AnyObject, o.SHA1)
		if err != nil {
			return nil, err
		}
	}

	if e != nil {
		r, err := e.Reader()
		if err != nil {
			return nil, err
		}

		buf := make([]byte, e.Size())
		if _, err = r.Read(buf); err != nil {
			return nil, err
		}

		return buf, nil
	}

	var data []byte
	if o.DiskType.IsDelta() {
		base, err := p.get(o.Parent)
		if err != nil {
			return nil, err
		}

		data, err = p.resolveObject(o, base)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		data, err = p.readData(o)
		if err != nil {
			return nil, err
		}
	}

	if len(o.Children) > 0 {
		m := &plumbing.MemoryObject{}
		m.Write(data)
		m.SetType(o.Type)
		m.SetSize(o.Size())
		p.cache.Put(m)
	}

	return data, nil
}

func (p *Parser) resolveObject(
	o *objectInfo,
	base []byte,
) ([]byte, error) {
	if !o.DiskType.IsDelta() {
		return nil, nil
	}

	data, err := p.readData(o)
	if err != nil {
		return nil, err
	}

	data, err = applyPatchBase(o, data, base)
	if err != nil {
		return nil, err
	}

	if pending, ok := p.pendingRefDeltas[o.SHA1]; ok {
		for _, po := range pending {
			po.Parent = o
			o.Children = append(o.Children, po)
		}
		delete(p.pendingRefDeltas, o.SHA1)
	}

	if p.storage != nil {
		obj := new(plumbing.MemoryObject)
		obj.SetSize(o.Size())
		obj.SetType(o.Type)
		if _, err := obj.Write(data); err != nil {
			return nil, err
		}

		if _, err := p.storage.SetEncodedObject(obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (p *Parser) readData(o *objectInfo) ([]byte, error) {
	if !p.scanner.IsSeekable && o.DiskType.IsDelta() {
		data, ok := p.deltas[o.Offset]
		if !ok {
			return nil, ErrDeltaNotCached
		}

		return data, nil
	}

	// TODO: skip header. Header size can be calculated with the offset of the
	// next offset in the first pass.
	if _, err := p.scanner.SeekFromStart(o.Offset); err != nil {
		return nil, err
	}

	if _, err := p.scanner.NextObjectHeader(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, _, err := p.scanner.NextObject(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func applyPatchBase(ota *objectInfo, data, base []byte) ([]byte, error) {
	patched, err := PatchDelta(base, data)
	if err != nil {
		return nil, err
	}

	ota.Type = ota.Parent.Type
	ota.Hasher = plumbing.NewHasher(ota.Type, int64(len(patched)))
	if _, err := ota.Write(patched); err != nil {
		return nil, err
	}
	ota.SHA1 = ota.Sum()
	ota.Length = int64(len(patched))

	return patched, nil
}

type objectInfo struct {
	plumbing.Hasher

	Offset       int64
	Length       int64
	HeaderLength int64
	PackSize     int64
	Type         plumbing.ObjectType
	DiskType     plumbing.ObjectType

	Crc32 uint32

	Parent   *objectInfo
	Children []*objectInfo
	SHA1     plumbing.Hash

	Content []byte
}

func newBaseObject(offset, length int64, t plumbing.ObjectType) *objectInfo {
	return newDeltaObject(offset, length, t, nil)
}

func newDeltaObject(
	offset, length int64,
	t plumbing.ObjectType,
	parent *objectInfo,
) *objectInfo {
	children := make([]*objectInfo, 0)

	obj := &objectInfo{
		Hasher:   plumbing.NewHasher(t, length),
		Offset:   offset,
		Length:   length,
		PackSize: 0,
		Type:     t,
		DiskType: t,
		Crc32:    0,
		Parent:   parent,
		Children: children,
	}

	return obj
}

func (o *objectInfo) Write(b []byte) (int, error) {
	o.Content = make([]byte, len(b))
	copy(o.Content, b)
	return o.Hasher.Write(b)
}

func (o *objectInfo) IsDelta() bool {
	return o.Type.IsDelta()
}

func (o *objectInfo) Size() int64 {
	return o.Length
}
