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
	"github.com/sniperkit/snk.fork.go-git.v4/plumbing/object"
)

// Basic example of how to list tags.
func main() {
	CheckArgs("<path>")
	path := os.Args[1]

	// We instance a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// List all tag references, both lightweight tags and annotated tags
	Info("git show-ref --tag")

	tagrefs, err := r.Tags()
	CheckIfError(err)
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t)
		return nil
	})
	CheckIfError(err)

	// Print each annotated tag object (lightweight tags are not included)
	Info("for t in $(git show-ref --tag); do if [ \"$(git cat-file -t $t)\" = \"tag\" ]; then git cat-file -p $t ; fi; done")

	tags, err := r.TagObjects()
	CheckIfError(err)
	err = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t)
		return nil
	})
	CheckIfError(err)
}
