/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package main

import (
	"os"

	"github.com/sniperkit/snk.fork.go-git.v4"
	. "github.com/sniperkit/snk.fork.go-git.v4/_examples"
)

// Example of how to open a repository in a specific path, and push to
// its default remote (origin).
func main() {
	CheckArgs("<repository-path>")
	path := os.Args[1]

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	Info("git push")
	// push using default options
	err = r.Push(&git.PushOptions{})
	CheckIfError(err)
}
