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

// Basic example of how to checkout a specific commit.
func main() {
	CheckArgs("<url>", "<directory>", "<commit>")
	url, directory, commit := os.Args[1], os.Args[2], os.Args[3]

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
	})

	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD, it's shows that the
	// repository is poiting to the giving commit in detached mode
	Info("git show-ref --head HEAD")
	ref, err = r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())
}
