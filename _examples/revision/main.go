/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package main

import (
	"fmt"
	"os"

	"github.com/sniperkit/snk.fork.go-git.v4"
	. "github.com/sniperkit/snk.fork.go-git.v4/_examples"
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing"
)

// Example how to resolve a revision into its commit counterpart
func main() {
	CheckArgs("<path>", "<revision>")

	path := os.Args[1]
	revision := os.Args[2]

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Resolve revision into a sha1 commit, only some revisions are resolved
	// look at the doc to get more details
	Info("git rev-parse %s", revision)

	h, err := r.ResolveRevision(plumbing.Revision(revision))

	CheckIfError(err)

	fmt.Println(h.String())
}
