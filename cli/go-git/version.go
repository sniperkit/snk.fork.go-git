/*
Sniperkit-Bot
- Date: 2018-08-11 15:40:00.935176804 +0200 CEST m=+0.032827986
- Status: analyzed
*/

package main

import "fmt"

var version string
var build string

type CmdVersion struct{}

func (c *CmdVersion) Execute(args []string) error {
	fmt.Printf("%s (%s) - build %s\n", bin, version, build)

	return nil
}
