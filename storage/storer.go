/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package storage

import (
	"github.com/sniperkit/snk.fork.go-git.v4/config"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/storer"
)

// Storer is a generic storage of objects, references and any information
// related to a particular repository. The package github.com/sniperkit/snk.fork.go-git.v4/storage
// contains two implementation a filesystem base implementation (such as `.git`)
// and a memory implementations being ephemeral
type Storer interface {
	storer.EncodedObjectStorer
	storer.ReferenceStorer
	storer.ShallowStorer
	storer.IndexStorer
	config.ConfigStorer
	ModuleStorer
}

// ModuleStorer allows interact with the modules' Storers
type ModuleStorer interface {
	// Module returns a Storer representing a submodule, if not exists returns a
	// new empty Storer is returned
	Module(name string) (Storer, error)
}
