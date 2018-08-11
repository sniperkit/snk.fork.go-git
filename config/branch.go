/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package config

import (
	"errors"

	"github.com/sniperkit/snk.fork.go-git.v4/plumbing"
	format "github.com/sniperkit/snk.fork.go-git.v4/plumbing/format/config"
)

var (
	errBranchEmptyName    = errors.New("branch config: empty name")
	errBranchInvalidMerge = errors.New("branch config: invalid merge")
)

// Branch contains information on the
// local branches and which remote to track
type Branch struct {
	// Name of branch
	Name string
	// Remote name of remote to track
	Remote string
	// Merge is the local refspec for the branch
	Merge plumbing.ReferenceName

	raw *format.Subsection
}

// Validate validates fields of branch
func (b *Branch) Validate() error {
	if b.Name == "" {
		return errBranchEmptyName
	}

	if b.Merge != "" && !b.Merge.IsBranch() {
		return errBranchInvalidMerge
	}

	return nil
}

func (b *Branch) marshal() *format.Subsection {
	if b.raw == nil {
		b.raw = &format.Subsection{}
	}

	b.raw.Name = b.Name

	if b.Remote == "" {
		b.raw.RemoveOption(remoteSection)
	} else {
		b.raw.SetOption(remoteSection, b.Remote)
	}

	if b.Merge == "" {
		b.raw.RemoveOption(mergeKey)
	} else {
		b.raw.SetOption(mergeKey, string(b.Merge))
	}

	return b.raw
}

func (b *Branch) unmarshal(s *format.Subsection) error {
	b.raw = s

	b.Name = b.raw.Name
	b.Remote = b.raw.Options.Get(remoteSection)
	b.Merge = plumbing.ReferenceName(b.raw.Options.Get(mergeKey))

	return b.Validate()
}
